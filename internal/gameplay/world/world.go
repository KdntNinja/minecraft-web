package world

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/progress"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
	"github.com/KdntNinja/webcraft/internal/gameplay/world/chunks"
	"github.com/KdntNinja/webcraft/internal/generation"
)

// AsyncUpdateTask represents a task for async world updates
type AsyncUpdateTask struct {
	taskType string
	data     interface{}
	callback func(interface{})
}

type World struct {
	ChunkManager *chunks.ChunkManager // Dynamic chunk loading system
	Entities     entity.Entities      // All entities in the world

	// Performance optimization caches
	cachedGrid        [][]int // Cached collision grid
	cachedGridOffsetX int     // Cached grid offset X
	cachedGridOffsetY int     // Cached grid offset Y
	gridDirty         bool    // Flag to indicate grid needs regeneration

	// Async update system
	updateTasks      chan AsyncUpdateTask
	updateWorkers    sync.WaitGroup
	updateCtx        context.Context
	updateCancel     context.CancelFunc
	numUpdateWorkers int

	// Grid generation pool
	gridGenerationMutex sync.RWMutex
	gridGenerationPool  sync.Pool
}

// NewWorld constructs a new World instance with dynamic chunk loading
func NewWorld(seed int64) *World {
	// Step: World Setup
	progress.UpdateCurrentStepProgress(1, "Setting up world structure...")
	generation.ResetWorldGeneration(seed)
	progress.UpdateCurrentStepProgress(2, "Reset world generation")

	w := &World{
		ChunkManager:     chunks.NewChunkManager(settings.ChunkViewDistance),
		Entities:         entity.Entities{},
		gridDirty:        true,
		numUpdateWorkers: runtime.NumCPU(),                // Use all available CPUs for updates
		updateTasks:      make(chan AsyncUpdateTask, 100), // Buffered channel for update tasks
	}

	// Initialize async update system
	w.updateCtx, w.updateCancel = context.WithCancel(context.Background())
	w.startUpdateWorkers()

	// Initialize grid generation pool
	w.gridGenerationPool = sync.Pool{
		New: func() interface{} {
			return make([][]int, 0, 256) // Pre-allocate slice capacity
		},
	}

	progress.UpdateCurrentStepProgress(3, "Created world structure")
	progress.CompleteCurrentStep()

	// Step: Finding Player Spawn
	progress.UpdateCurrentStepProgress(1, "Finding spawn location...")
	spawnPoint := chunks.FindSafeSpawnPoint()
	progress.UpdateCurrentStepProgress(2, fmt.Sprintf("Found spawn at (%.1f, %.1f)", spawnPoint.X, spawnPoint.Y))
	progress.CompleteCurrentStep()

	// Step: Spawning Player (before chunk loading)
	progress.UpdateCurrentStepProgress(1, "Creating player entity...")
	playerEntity := player.NewPlayer(spawnPoint.X, spawnPoint.Y, w)
	w.Entities = append(w.Entities, playerEntity)
	progress.UpdateCurrentStepProgress(2, "Created player entity")
	progress.CompleteCurrentStep()

	// Step: Generating Terrain around player
	progress.UpdateCurrentStepProgress(1, "Loading initial chunks around player...")
	w.ChunkManager.InitialLoadWithProgress(playerEntity.X, playerEntity.Y)
	progress.CompleteCurrentStep()

	// Step: Finalizing
	progress.UpdateCurrentStepProgress(1, "World generation finished!")
	progress.CompleteCurrentStep()

	// Initialize the grid generation pool
	w.gridGenerationPool = sync.Pool{
		New: func() interface{} {
			// Allocate a new grid slice
			return make([][]int, 0)
		},
	}

	// Start async update workers
	w.updateCtx, w.updateCancel = context.WithCancel(context.Background())
	for i := 0; i < w.numUpdateWorkers; i++ {
		w.updateWorkers.Add(1)
		go w.updateWorker(i)
	}

	return w
}

// startUpdateWorkers starts the worker pool for async world updates
func (w *World) startUpdateWorkers() {
	for i := 0; i < w.numUpdateWorkers; i++ {
		w.updateWorkers.Add(1)
		go w.updateWorker(i)
	}
}

// updateWorker processes async update tasks
func (w *World) updateWorker(workerID int) {
	defer w.updateWorkers.Done()

	for {
		select {
		case <-w.updateCtx.Done():
			return
		case task := <-w.updateTasks:
			// Process different types of update tasks
			switch task.taskType {
			case "grid_generation":
				w.processGridGeneration(task.data, task.callback)
			case "entity_update":
				w.processEntityUpdate(task.data, task.callback)
			case "physics_update":
				w.processPhysicsUpdate(task.data, task.callback)
			}
		}
	}
}

// processGridGeneration handles grid generation tasks
func (w *World) processGridGeneration(data interface{}, callback func(interface{})) {
	if _, ok := data.(map[string]interface{}); ok {
		// Generate grid asynchronously
		allChunks := w.ChunkManager.GetAllChunks()
		result := w.generateGridDataAsync(allChunks)
		if callback != nil {
			callback(result)
		}
	}
}

// processEntityUpdate handles entity update tasks
func (w *World) processEntityUpdate(data interface{}, callback func(interface{})) {
	if entities, ok := data.([]entity.Entity); ok {
		// Update entities in parallel
		var wg sync.WaitGroup
		for _, e := range entities {
			wg.Add(1)
			go func(ent entity.Entity) {
				defer wg.Done()
				ent.Update()
			}(e)
		}
		wg.Wait()
		if callback != nil {
			callback(nil)
		}
	}
}

// processPhysicsUpdate handles physics update tasks
func (w *World) processPhysicsUpdate(data interface{}, callback func(interface{})) {
	if _, ok := data.(map[string]interface{}); ok {
		// Process physics updates asynchronously
		if callback != nil {
			callback(nil)
		}
	}
}

// asyncUpdateWorker processes update tasks from the channel
func (w *World) asyncUpdateWorker() {
	defer w.updateWorkers.Done()
	for {
		select {
		case task := <-w.updateTasks:
			// Process the update task
			switch task.taskType {
			case "entityUpdate":
				if entities, ok := task.data.(entity.Entities); ok {
					for i := range entities {
						entities[i].Update()
					}
				}
			}

			// Call the task's callback function if provided
			if task.callback != nil {
				task.callback(task.data)
			}

		case <-w.updateCtx.Done():
			return // Context cancelled, exit the worker
		}
	}
}

// generateGridDataAsync generates collision grid data asynchronously
func (w *World) generateGridDataAsync(allChunks map[chunks.ChunkCoord]*block.Chunk) interface{} {
	if len(allChunks) == 0 {
		return map[string]interface{}{
			"grid":    [][]int{},
			"offsetX": 0,
			"offsetY": 0,
		}
	}

	// Calculate bounds
	minX, maxX, minY, maxY := 0, 0, 0, 0
	first := true
	for coord := range allChunks {
		if first {
			minX, maxX, minY, maxY = coord.X, coord.X, coord.Y, coord.Y
			first = false
		} else {
			if coord.X < minX {
				minX = coord.X
			}
			if coord.X > maxX {
				maxX = coord.X
			}
			if coord.Y < minY {
				minY = coord.Y
			}
			if coord.Y > maxY {
				maxY = coord.Y
			}
		}
	}

	width := (maxX - minX + 1) * settings.ChunkWidth
	height := (maxY - minY + 1) * settings.ChunkHeight

	// Get grid from pool or create new
	var grid [][]int
	if poolGrid := w.gridGenerationPool.Get(); poolGrid != nil {
		if pGrid, ok := poolGrid.([][]int); ok {
			grid = pGrid[:0] // Reset slice but keep capacity
		}
	}
	if grid == nil {
		grid = make([][]int, height)
	}

	// Ensure grid has correct dimensions
	for len(grid) < height {
		grid = append(grid, make([]int, width))
	}
	for i := 0; i < height; i++ {
		if len(grid[i]) < width {
			grid[i] = make([]int, width)
		}
	}

	// Fill grid in parallel
	var wg sync.WaitGroup
	for coord, chunk := range allChunks {
		wg.Add(1)
		go func(coord chunks.ChunkCoord, chunk *block.Chunk) {
			defer wg.Done()
			for y := 0; y < settings.ChunkHeight; y++ {
				for x := 0; x < settings.ChunkWidth; x++ {
					globalX := (coord.X-minX)*settings.ChunkWidth + x
					globalY := (coord.Y-minY)*settings.ChunkHeight + y
					if y < len((*chunk)) && x < len((*chunk)[y]) {
						grid[globalY][globalX] = int((*chunk)[y][x])
					} else {
						grid[globalY][globalX] = int(block.Air)
					}
				}
			}
		}(coord, chunk)
	}
	wg.Wait()

	return map[string]interface{}{
		"grid":    grid,
		"offsetX": minX * settings.ChunkWidth,
		"offsetY": minY * settings.ChunkHeight,
	}
}

// ToIntGrid flattens the world's blocks into a [][]int grid for entity collision, and returns the offset (minX, minY)
// Uses caching to avoid regenerating the grid every frame
func (w *World) ToIntGrid() ([][]int, int, int) {
	allChunks := w.ChunkManager.GetAllChunks()
	if len(allChunks) == 0 {
		return [][]int{}, 0, 0
	}

	// Regenerate grid only when necessary
	minX, maxX, minY, maxY := 0, 0, 0, 0
	first := true
	for coord := range allChunks {
		if first {
			minX, maxX, minY, maxY = coord.X, coord.X, coord.Y, coord.Y
			first = false
		} else {
			if coord.X < minX {
				minX = coord.X
			}
			if coord.X > maxX {
				maxX = coord.X
			}
			if coord.Y < minY {
				minY = coord.Y
			}
			if coord.Y > maxY {
				maxY = coord.Y
			}
		}
	}
	width := (maxX - minX + 1) * settings.ChunkWidth
	height := (maxY - minY + 1) * settings.ChunkHeight

	// Always allocate a new grid to avoid stale data when moving into negative coords
	grid := make([][]int, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
	}

	var wg sync.WaitGroup
	for coord, chunk := range allChunks {
		wg.Add(1)
		go func(coord chunks.ChunkCoord, chunk *block.Chunk) {
			defer wg.Done()
			for y := 0; y < settings.ChunkHeight; y++ {
				for x := 0; x < settings.ChunkWidth; x++ {
					globalX := (coord.X-minX)*settings.ChunkWidth + x
					globalY := (coord.Y-minY)*settings.ChunkHeight + y
					if y < len((*chunk)) && x < len((*chunk)[y]) {
						grid[globalY][globalX] = int((*chunk)[y][x])
					} else {
						grid[globalY][globalX] = int(block.Air)
					}
				}
			}
		}(coord, chunk)
	}
	wg.Wait()

	w.cachedGrid = grid
	w.cachedGridOffsetX = minX * settings.ChunkWidth
	w.cachedGridOffsetY = minY * settings.ChunkHeight
	w.gridDirty = false

	return grid, w.cachedGridOffsetX, w.cachedGridOffsetY
}

// IsGridDirty returns whether the collision grid needs to be regenerated
func (w *World) IsGridDirty() bool {
	return w.gridDirty
}

// GetCachedGrid returns the cached collision grid without regenerating it
func (w *World) GetCachedGrid() ([][]int, int, int) {
	if w.cachedGrid == nil {
		return [][]int{}, 0, 0
	}
	return w.cachedGrid, w.cachedGridOffsetX, w.cachedGridOffsetY
}

// Block interaction methods

// GetBlockAt returns the block type at the given world coordinates
func (w *World) GetBlockAt(blockX, blockY int) block.BlockType {
	return w.ChunkManager.GetBlock(blockX, blockY)
}

// SetBlockAt sets the block type at the given world coordinates
func (w *World) SetBlockAt(blockX, blockY int, blockType block.BlockType) bool {
	success := w.ChunkManager.SetBlock(blockX, blockY, blockType)

	if success {
		w.updateCachedGridBlock(blockX, blockY, blockType)
		// Do NOT always mark gridDirty here; updateCachedGridBlock will do so only if needed
	}

	return success
}

// BreakBlock removes a block at the given coordinates
func (w *World) BreakBlock(blockX, blockY int) bool {
	currentBlock := w.GetBlockAt(blockX, blockY)
	if currentBlock == block.Air {
		return false // Cannot break air
	}

	return w.SetBlockAt(blockX, blockY, block.Air)
}

// wouldBlockCollideWithEntity checks if placing a block at the given coordinates would collide with any entity
func (w *World) wouldBlockCollideWithEntity(blockX, blockY int) bool {
	// Convert block coordinates to world coordinates
	blockWorldX := float64(blockX * settings.TileSize)
	blockWorldY := float64(blockY * settings.TileSize)
	blockWidth := float64(settings.TileSize)
	blockHeight := float64(settings.TileSize)

	// Check collision with all entities
	for _, entity := range w.Entities {
		entityX, entityY := entity.GetPosition()

		// Get entity dimensions based on type
		var entityWidth, entityHeight float64
		if player, ok := entity.(*player.Player); ok {
			entityWidth = float64(player.AABB.Width)
			entityHeight = float64(player.AABB.Height)
		} else {
			// Default entity size for other entity types
			entityWidth = float64(settings.TileSize)
			entityHeight = float64(settings.TileSize)
		}

		// Check AABB collision between entity and potential block position
		if entityX < blockWorldX+blockWidth &&
			entityX+entityWidth > blockWorldX &&
			entityY < blockWorldY+blockHeight &&
			entityY+entityHeight > blockWorldY {
			return true // Collision detected
		}
	}
	return false // No collision
}

// PlaceBlock places a block at the given coordinates
func (w *World) PlaceBlock(blockX, blockY int, blockType block.BlockType) bool {
	currentBlock := w.GetBlockAt(blockX, blockY)
	if currentBlock != block.Air {
		return false // Cannot place block where one already exists
	}

	// Don't allow placing air blocks
	if blockType == block.Air {
		return false
	}

	// Prevent placing a block inside any entity (including player)
	if w.wouldBlockCollideWithEntity(blockX, blockY) {
		return false
	}

	return w.SetBlockAt(blockX, blockY, blockType)
}

// Update method to handle dynamic chunk loading
func (w *World) Update() {
	// Update chunk loading based on player position
	if len(w.Entities) > 0 {
		if player, ok := w.Entities[0].(*player.Player); ok {
			w.ChunkManager.UpdatePlayerPosition(player.X, player.Y)
		}
	}

	// Update entities directly in parallel (more efficient than task queue for this)
	var wg sync.WaitGroup
	for _, e := range w.Entities {
		wg.Add(1)
		go func(ent entity.Entity) {
			defer wg.Done()
			ent.Update()
		}(e)
	}
	wg.Wait()
}

// GetChunkCount returns the number of currently loaded chunks
func (w *World) GetChunkCount() int {
	return w.ChunkManager.GetLoadedChunkCount()
}

// FindSurfaceHeight finds the surface height at the given X coordinate
func FindSurfaceHeight(blockX int, w *World) int {
	// Find which chunk this block belongs to
	chunkX, _ := chunks.BlockToChunk(blockX, 0)

	// We need to look through multiple chunks vertically to find surface
	for chunkY := 0; chunkY < 10; chunkY++ { // Search downward through chunks
		chunk := w.ChunkManager.GetChunk(chunkX, chunkY)

		// Look through this chunk for the surface
		localX := blockX - (chunkX * settings.ChunkWidth)
		if localX < 0 || localX >= settings.ChunkWidth {
			continue
		}

		for localY := 0; localY < settings.ChunkHeight; localY++ {
			if chunk[localY][localX] != block.Air {
				// Found the first non-air block - this is the surface
				return (chunkY * settings.ChunkHeight) + localY
			}
		}
	}

	// Default surface height if not found
	return settings.SurfaceBaseHeight
}

// updateCachedGridBlock efficiently updates a single block in the cached collision grid
func (w *World) updateCachedGridBlock(blockX, blockY int, blockType block.BlockType) {
	// Only update if we have a cached grid
	if w.cachedGrid == nil {
		return
	}

	// Convert world coordinates to grid coordinates
	gridX := blockX - w.cachedGridOffsetX
	gridY := blockY - w.cachedGridOffsetY

	// Bounds check for the cached grid
	if gridY < 0 || gridY >= len(w.cachedGrid) || gridX < 0 || gridX >= len(w.cachedGrid[0]) {
		// Block is outside cached grid bounds, mark as dirty for next full regeneration
		w.gridDirty = true
		return
	}

	// Update just this one block in the cached grid
	w.cachedGrid[gridY][gridX] = int(blockType)

	// Grid is still valid, no need to mark as dirty
}

// GetChunksForRendering returns all currently loaded chunks for rendering
func (w *World) GetChunksForRendering() map[chunks.ChunkCoord]*block.Chunk {
	return w.ChunkManager.GetAllChunks()
}

// Stop stops the world and cleans up resources
func (w *World) Stop() {
	// Cancel the update context to stop all workers
	w.updateCancel()

	// Wait for all update workers to finish
	w.updateWorkers.Wait()
	close(w.updateTasks)
	w.ChunkManager.Shutdown()
	fmt.Println("WORLD: Shutdown complete")
}
