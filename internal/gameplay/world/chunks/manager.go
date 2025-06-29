package chunks

import (
	"fmt"
	"math"
	"sync"

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

// ChunkManager handles dynamic chunk loading and unloading
type ChunkManager struct {
	chunks          map[ChunkCoord]*block.Chunk
	loadedChunks    map[ChunkCoord]bool
	mutex           sync.RWMutex
	viewDistance    int // Chunks to keep loaded around player
	lastPlayerChunk ChunkCoord
}

// NewChunkManager creates a new chunk manager
func NewChunkManager(viewDistance int) *ChunkManager {
	return &ChunkManager{
		chunks:          make(map[ChunkCoord]*block.Chunk),
		loadedChunks:    make(map[ChunkCoord]bool),
		viewDistance:    viewDistance,
		lastPlayerChunk: ChunkCoord{X: math.MaxInt32, Y: math.MaxInt32}, // Force initial load
	}
}

// GetChunk returns a chunk at the given coordinates, generating it if necessary
func (cm *ChunkManager) GetChunk(chunkX, chunkY int) *block.Chunk {
	coord := ChunkCoord{X: chunkX, Y: chunkY}

	cm.mutex.RLock()
	chunk, exists := cm.chunks[coord]
	cm.mutex.RUnlock()

	if exists {
		return chunk
	}

	// Generate the chunk
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Double-check after acquiring write lock
	if chunk, exists := cm.chunks[coord]; exists {
		return chunk
	}

	// Generate new chunk
	newChunk := generation.GenerateChunk(chunkX, chunkY)
	cm.chunks[coord] = &newChunk
	cm.loadedChunks[coord] = true

	fmt.Printf("CHUNK_MANAGER: Generated chunk (%d, %d)\n", chunkX, chunkY)
	return &newChunk
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
			cm.mutex.RUnlock()

			if !exists {
				chunksToLoad = append(chunksToLoad, coord)
			}
		}
	}

	// Only load up to MaxChunksPerFrame per call to reduce stutter
	for i := 0; i < len(chunksToLoad) && i < MaxChunksPerFrame; i++ {
		coord := chunksToLoad[i]
		go cm.GetChunk(coord.X, coord.Y)
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

	// Load chunks around spawn point
	for dx := -cm.viewDistance; dx <= cm.viewDistance; dx++ {
		for dy := -cm.viewDistance; dy <= cm.viewDistance; dy++ {
			chunkX := spawnChunkX + dx
			chunkY := spawnChunkY + dy

			cm.GetChunk(chunkX, chunkY)
			generatedChunks++

			// Update progress for each chunk
			progress.UpdateCurrentStepProgress(generatedChunks,
				fmt.Sprintf("Loaded chunk %d/%d at (%d, %d)", generatedChunks, totalChunks, chunkX, chunkY))
		}
	}

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
