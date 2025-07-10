package chunks

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/progress"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation"
)

// ChunkCoord represents a chunk coordinate
type ChunkCoord struct {
	X int
	Y int
}

// ChunkGenerationJob represents a chunk generation job
type ChunkGenerationJob struct {
	coord    ChunkCoord
	priority int // Higher values = higher priority
}

// ChunkManager handles dynamic chunk loading and unloading with multithreaded generation
type ChunkManager struct {
	chunks          map[ChunkCoord]*block.Chunk
	loadedChunks    map[ChunkCoord]bool
	generating      map[ChunkCoord]bool     // Track chunks being generated
	chunkQueue      chan chunkResult        // Channel for async chunk results
	jobQueue        chan ChunkGenerationJob // Job queue for generation workers
	mutex           sync.RWMutex
	viewDistance    int // Chunks to keep loaded around player
	lastPlayerChunk ChunkCoord

	// Worker pool for chunk generation
	workerPool   sync.WaitGroup
	workerCtx    context.Context
	workerCancel context.CancelFunc
	numWorkers   int

	// Performance metrics
	generationMetrics struct {
		totalGenerated int64
		totalTime      time.Duration
		mutex          sync.RWMutex
	}
}

// NewChunkManager creates a new chunk manager
type chunkResult struct {
	coord          ChunkCoord
	chunk          *block.Chunk
	generationTime time.Duration
}

func NewChunkManager(viewDistance int) *ChunkManager {
	cm := &ChunkManager{
		chunks:          make(map[ChunkCoord]*block.Chunk),
		loadedChunks:    make(map[ChunkCoord]bool),
		generating:      make(map[ChunkCoord]bool),
		chunkQueue:      make(chan chunkResult, 32),
		jobQueue:        make(chan ChunkGenerationJob, 64), // Buffered job queue
		viewDistance:    viewDistance,
		lastPlayerChunk: ChunkCoord{X: math.MaxInt32, Y: math.MaxInt32}, // Force initial load
		numWorkers:      runtime.NumCPU(),                               // Use all available CPUs
	}
	go cm.chunkInsertWorker()
	go cm.startWorkerPool()
	return cm
}

// chunkInsertWorker runs in the background and inserts generated chunks into the map
func (cm *ChunkManager) chunkInsertWorker() {
	for res := range cm.chunkQueue {
		cm.mutex.Lock()
		cm.chunks[res.coord] = res.chunk
		cm.loadedChunks[res.coord] = true
		delete(cm.generating, res.coord)
		cm.mutex.Unlock()

		// Log performance if slow
		if res.generationTime > 100*time.Millisecond {
			fmt.Printf("CHUNK_MANAGER: Slow chunk generation at (%d, %d): %v\n",
				res.coord.X, res.coord.Y, res.generationTime)
		}
	}
}

// startWorkerPool starts the worker pool for chunk generation
func (cm *ChunkManager) startWorkerPool() {
	cm.workerCtx, cm.workerCancel = context.WithCancel(context.Background())
	for i := 0; i < cm.numWorkers; i++ {
		cm.workerPool.Add(1)
		go cm.worker(i)
	}
}

// worker is a single worker goroutine for chunk generation
func (cm *ChunkManager) worker(workerID int) {
	defer cm.workerPool.Done()
	fmt.Printf("CHUNK_MANAGER: Worker %d starting\n", workerID)

	for {
		select {
		case <-cm.workerCtx.Done():
			fmt.Printf("CHUNK_MANAGER: Worker %d shutting down\n", workerID)
			return
		case job := <-cm.jobQueue:
			start := time.Now()

			// Generate the chunk
			chunk := generation.GenerateChunk(job.coord.X, job.coord.Y)
			generationTime := time.Since(start)

			// Update metrics
			cm.generationMetrics.mutex.Lock()
			cm.generationMetrics.totalGenerated++
			cm.generationMetrics.totalTime += generationTime
			cm.generationMetrics.mutex.Unlock()

			// Send result to insertion worker
			result := chunkResult{
				coord:          job.coord,
				chunk:          &chunk,
				generationTime: generationTime,
			}

			select {
			case cm.chunkQueue <- result:
				fmt.Printf("CHUNK_MANAGER: Worker %d generated chunk (%d, %d) in %v\n",
					workerID, job.coord.X, job.coord.Y, generationTime)
			case <-cm.workerCtx.Done():
				return
			}
		}
	}
}

// GetChunk returns a chunk at the given coordinates, generating it if necessary
func (cm *ChunkManager) GetChunk(chunkX, chunkY int) *block.Chunk {
	coord := ChunkCoord{X: chunkX, Y: chunkY}

	cm.mutex.RLock()
	chunk, exists := cm.chunks[coord]
	generating := cm.generating[coord]
	cm.mutex.RUnlock()

	if exists {
		return chunk
	}
	if generating {
		return nil // Still generating, return nil for now
	}

	// Mark as generating and start async generation
	cm.mutex.Lock()
	if _, already := cm.generating[coord]; !already {
		cm.generating[coord] = true
		// Add job to the queue with priority based on distance from player
		playerDist := int(math.Abs(float64(coord.X-cm.lastPlayerChunk.X))) + int(math.Abs(float64(coord.Y-cm.lastPlayerChunk.Y)))
		priority := cm.viewDistance*2 - playerDist // Closer chunks have higher priority
		if priority < 0 {
			priority = 0
		}
		cm.jobQueue <- ChunkGenerationJob{coord: coord, priority: priority}
	}
	cm.mutex.Unlock()
	return nil // Not ready yet
}

// UpdatePlayerPosition updates the chunk loading based on player position
func (cm *ChunkManager) UpdatePlayerPosition(playerX, playerY float64) {
	// Calculate player's current chunk
	chunkX := int(math.Floor(playerX / float64(settings.ChunkWidth*settings.TileSize)))
	chunkY := int(math.Floor(playerY / float64(settings.ChunkHeight*settings.TileSize)))

	currentChunk := ChunkCoord{X: chunkX, Y: chunkY}

	// Always check for new chunks, not just on chunk change
	cm.loadChunksAroundPlayer(chunkX, chunkY)

	// Unload distant chunks
	cm.unloadDistantChunks(chunkX, chunkY)

	// Only update lastPlayerChunk if player moved to a different chunk
	if currentChunk != cm.lastPlayerChunk {
		fmt.Printf("CHUNK_MANAGER: Player moved to chunk (%d, %d)\n", chunkX, chunkY)
		cm.lastPlayerChunk = currentChunk
	}
}

// Load up to N chunks per frame to reduce stutter
const MaxChunksPerFrame = 2

// loadChunksAroundPlayer loads chunks within view distance of the player (no bias, true square)
func (cm *ChunkManager) loadChunksAroundPlayer(playerChunkX, playerChunkY int) {
	loadCount := 0
	chunksToLoad := []ChunkCoord{}

	for dx := -cm.viewDistance; dx <= cm.viewDistance; dx++ {
		for dy := -cm.viewDistance; dy <= cm.viewDistance; dy++ {
			chunkX := playerChunkX + dx
			chunkY := playerChunkY + dy

			coord := ChunkCoord{X: chunkX, Y: chunkY}

			cm.mutex.RLock()
			_, exists := cm.chunks[coord]
			generating := cm.generating[coord]
			cm.mutex.RUnlock()

			if !exists && !generating {
				chunksToLoad = append(chunksToLoad, coord)
			}
		}
	}

	// Only load up to MaxChunksPerFrame per call to reduce stutter
	for i := 0; i < len(chunksToLoad) && i < MaxChunksPerFrame; i++ {
		coord := chunksToLoad[i]
		cm.GetChunk(coord.X, coord.Y) // Will start async generation if not present
		loadCount++
	}

	if loadCount > 0 {
		fmt.Printf("CHUNK_MANAGER: Loading %d new chunks around player (limited per frame)\n", loadCount)
	}
}

// unloadDistantChunks unloads chunks that are too far from the player
func (cm *ChunkManager) unloadDistantChunks(playerChunkX, playerChunkY int) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	unloadDistance := cm.viewDistance + 2 // Keep 2 extra chunks before unloading
	var toUnload []ChunkCoord

	for coord := range cm.chunks {
		dx := coord.X - playerChunkX
		dy := coord.Y - playerChunkY
		distance := int(math.Sqrt(float64(dx*dx + dy*dy)))

		if distance > unloadDistance {
			toUnload = append(toUnload, coord)
		}
	}

	for _, coord := range toUnload {
		delete(cm.chunks, coord)
		delete(cm.loadedChunks, coord)
	}

	if len(toUnload) > 0 {
		fmt.Printf("CHUNK_MANAGER: Unloaded %d distant chunks\n", len(toUnload))
	}
}

// InitialLoadWithProgress loads chunks around the spawn point during world creation and updates progress
func (cm *ChunkManager) InitialLoadWithProgress(spawnX, spawnY float64) {
	spawnChunkX := int(math.Floor(spawnX / float64(settings.ChunkWidth*settings.TileSize)))
	spawnChunkY := int(math.Floor(spawnY / float64(settings.ChunkHeight*settings.TileSize)))

	fmt.Printf("CHUNK_MANAGER: Initial load around spawn chunk (%d, %d)\n", spawnChunkX, spawnChunkY)

	// Calculate total chunks to load
	totalChunks := (cm.viewDistance*2 + 1) * (cm.viewDistance*2 + 1)
	generatedChunks := 0

	// Set progress for initial chunk loading
	progress.SetCurrentStepSubSteps(totalChunks, fmt.Sprintf("Loading %d initial chunks...", totalChunks))

	// Load chunks around spawn point in parallel
	var wg sync.WaitGroup
	for dx := -cm.viewDistance; dx <= cm.viewDistance; dx++ {
		for dy := -cm.viewDistance; dy <= cm.viewDistance; dy++ {
			chunkX := spawnChunkX + dx
			chunkY := spawnChunkY + dy
			wg.Add(1)
			go func(chunkX, chunkY int) {
				defer wg.Done()
				cm.GetChunk(chunkX, chunkY)
			}(chunkX, chunkY)
			generatedChunks++
			// Update progress for each chunk
			progress.UpdateCurrentStepProgress(generatedChunks,
				fmt.Sprintf("Loaded chunk %d/%d at (%d, %d)", generatedChunks, totalChunks, chunkX, chunkY))
		}
	}
	wg.Wait()
	cm.lastPlayerChunk = ChunkCoord{X: spawnChunkX, Y: spawnChunkY}
}

// GetAllChunks returns all currently loaded chunks (for rendering)
func (cm *ChunkManager) GetAllChunks() map[ChunkCoord]*block.Chunk {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// Return a shallow copy of the map (pointers, not arrays)
	result := make(map[ChunkCoord]*block.Chunk)
	for coord, chunk := range cm.chunks {
		result[coord] = chunk
	}
	return result
}

// GetLoadedChunkCount returns the number of currently loaded chunks
func (cm *ChunkManager) GetLoadedChunkCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.chunks)
}

// IsChunkLoaded checks if a chunk is currently loaded
func (cm *ChunkManager) IsChunkLoaded(chunkX, chunkY int) bool {
	coord := ChunkCoord{X: chunkX, Y: chunkY}
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.loadedChunks[coord]
}

// SetBlock sets a block at the given world coordinates through the chunk manager
func (cm *ChunkManager) SetBlock(blockX, blockY int, blockType block.BlockType) bool {
	chunkX := blockX / settings.ChunkWidth
	chunkY := blockY / settings.ChunkHeight
	inChunkX := blockX % settings.ChunkWidth
	inChunkY := blockY % settings.ChunkHeight

	// Handle negative coordinates properly
	if blockX < 0 {
		chunkX = (blockX - settings.ChunkWidth + 1) / settings.ChunkWidth
		inChunkX = ((blockX % settings.ChunkWidth) + settings.ChunkWidth) % settings.ChunkWidth
	}
	if blockY < 0 {
		chunkY = (blockY - settings.ChunkHeight + 1) / settings.ChunkHeight
		inChunkY = ((blockY % settings.ChunkHeight) + settings.ChunkHeight) % settings.ChunkHeight
	}

	// Get or generate the chunk
	chunk := cm.GetChunk(chunkX, chunkY)
	if chunk == nil || len(*chunk) == 0 {
		return false
	}

	// Bounds check
	if inChunkY < 0 || inChunkY >= settings.ChunkHeight || inChunkX < 0 || inChunkX >= settings.ChunkWidth {
		return false
	}

	// Set the block in the chunk
	cm.mutex.Lock()
	(*chunk)[inChunkY][inChunkX] = blockType
	cm.mutex.Unlock()

	return true
}

// GetBlock gets a block at the given world coordinates through the chunk manager
func (cm *ChunkManager) GetBlock(blockX, blockY int) block.BlockType {
	chunkX := blockX / settings.ChunkWidth
	chunkY := blockY / settings.ChunkHeight
	inChunkX := blockX % settings.ChunkWidth
	inChunkY := blockY % settings.ChunkHeight

	// Handle negative coordinates properly
	if blockX < 0 {
		chunkX = (blockX - settings.ChunkWidth + 1) / settings.ChunkWidth
		inChunkX = ((blockX % settings.ChunkWidth) + settings.ChunkWidth) % settings.ChunkWidth
	}
	if blockY < 0 {
		chunkY = (blockY - settings.ChunkHeight + 1) / settings.ChunkHeight
		inChunkY = ((blockY % settings.ChunkHeight) + settings.ChunkHeight) % settings.ChunkHeight
	}

	chunk := cm.GetChunk(chunkX, chunkY)
	if chunk == nil || len(*chunk) == 0 {
		return block.Air
	}

	// Bounds check
	if inChunkY < 0 || inChunkY >= settings.ChunkHeight || inChunkX < 0 || inChunkX >= settings.ChunkWidth {
		return block.Air
	}

	cm.mutex.RLock()
	blockType := (*chunk)[inChunkY][inChunkX]
	cm.mutex.RUnlock()

	return blockType
}

// ViewDistance returns the chunk view distance
func (cm *ChunkManager) ViewDistance() int {
	return cm.viewDistance
}

// GetGenerationMetrics returns performance metrics
func (cm *ChunkManager) GetGenerationMetrics() (int64, time.Duration) {
	cm.generationMetrics.mutex.RLock()
	defer cm.generationMetrics.mutex.RUnlock()
	return cm.generationMetrics.totalGenerated, cm.generationMetrics.totalTime
}

// Shutdown cleanly stops all workers
func (cm *ChunkManager) Shutdown() {
	fmt.Println("CHUNK_MANAGER: Shutting down...")
	cm.workerCancel()
	cm.workerPool.Wait()
	close(cm.chunkQueue)
	close(cm.jobQueue)
	fmt.Println("CHUNK_MANAGER: Shutdown complete")
}

// Stop stops the chunk manager and waits for workers to finish
func (cm *ChunkManager) Stop() {
	cm.workerCancel()
	cm.workerPool.Wait()
	fmt.Println("CHUNK_MANAGER: All workers stopped")
}
