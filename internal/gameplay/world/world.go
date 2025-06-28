package world

import (
	"fmt"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
	"github.com/KdntNinja/webcraft/internal/generation/terrain"
)

type ChunkCoord struct {
	X int
	Y int
}

type World struct {
	Chunks   map[ChunkCoord]block.Chunk // Fixed world: map of chunk coordinates to chunks
	Entities entity.Entities

	// Performance optimization caches
	cachedGrid        [][]int // Cached collision grid
	cachedGridOffsetX int     // Cached grid offset X
	cachedGridOffsetY int     // Cached grid offset Y
	gridDirty         bool    // Flag to indicate grid needs regeneration
}

// NewWorld constructs a new World instance with a fixed set of pre-generated chunks
func NewWorld(seed int64) *World {
	terrain.ResetWorldGeneration(seed)
	w := &World{
		Chunks:    make(map[ChunkCoord]block.Chunk),
		Entities:  entity.Entities{},
		gridDirty: true, // Grid needs initial generation
	}
	// Generate a large fixed world area - no dynamic loading
	worldWidth := settings.WorldChunksX  // Total chunks horizontally
	worldHeight := settings.WorldChunksY // Total chunks vertically

	fmt.Printf("DEBUG: WorldChunksX=%d, WorldChunksY=%d\n", worldWidth, worldHeight)

	// Calculate chunk range to center the world around (0,0)
	halfWidth := worldWidth / 2
	var startX, endX int
	if worldWidth%2 == 0 {
		// Even number of chunks: generate equal chunks on both sides
		// For 24 chunks: -12 to 11 (24 total)
		startX = -halfWidth
		endX = halfWidth - 1
	} else {
		// Odd number of chunks: center chunk at 0
		// For 25 chunks: -12 to 12 (25 total)
		startX = -halfWidth
		endX = halfWidth
	}

	fmt.Printf("DEBUG: Generating chunks from X=%d to X=%d (total: %d chunks)\n", startX, endX, endX-startX+1)

	for cy := 0; cy < worldHeight; cy++ {
		for cx := startX; cx <= endX; cx++ {
			coord := ChunkCoord{X: cx, Y: cy}
			w.Chunks[coord] = GenerateChunk(coord.X, coord.Y)
		}
	}

	fmt.Printf("DEBUG: Generated %d chunks total\n", len(w.Chunks))

	// Print first few and last few chunk coordinates for verification
	chunkCount := 0
	fmt.Printf("DEBUG: First 5 chunks generated: ")
	for coord := range w.Chunks {
		if chunkCount < 5 {
			fmt.Printf("(%d,%d) ", coord.X, coord.Y)
		}
		chunkCount++
	}
	fmt.Printf("\n")

	// Add player entity at pixel (0, 0) in world coordinates
	fmt.Printf("DEBUG: Player spawning at pixel (0, 0)\n")

	// Spawn player at pixel (0, 0), which corresponds to block (0, 0)
	spawnBlockX := 0
	spawnBlockY := 0

	// Find the surface height at block X=0 to determine proper Y spawn
	surfaceY := FindSurfaceHeight(spawnBlockX, w)
	fmt.Printf("DEBUG: Surface height at X=%d is Y=%d\n", spawnBlockX, surfaceY)

	// Spawn player 3 blocks above surface for safety
	spawnBlockY = surfaceY - 3

	// Ensure spawn position is reasonable
	if spawnBlockY < 5 {
		spawnBlockY = 5
	}
	if spawnBlockY > 200 {
		spawnBlockY = 200
	}

	// Convert to pixel coordinates - spawn at exact pixel (0, 0) for X, calculated Y
	px := 0.0
	py := float64(spawnBlockY * settings.TileSize)

	fmt.Printf("DEBUG: Final player spawn at pixel position (%f, %f), block position (%d, %d)\n", px, py, spawnBlockX, spawnBlockY)

	w.Entities = append(w.Entities, player.NewPlayer(px, py))
	return w
}

// ToIntGrid flattens the world's blocks into a [][]int grid for entity collision, and returns the offset (minX, minY)
// Uses caching to avoid regenerating the grid every frame
func (w *World) ToIntGrid() ([][]int, int, int) {
	if len(w.Chunks) == 0 {
		return [][]int{}, 0, 0
	}

	// Return cached grid if it's still valid
	if !w.gridDirty && w.cachedGrid != nil {
		return w.cachedGrid, w.cachedGridOffsetX, w.cachedGridOffsetY
	}

	// Regenerate grid only when necessary
	minX, maxX, minY, maxY := 0, 0, 0, 0
	first := true
	for coord := range w.Chunks {
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

	// Reuse cached grid if dimensions match, otherwise allocate new
	if w.cachedGrid == nil || len(w.cachedGrid) != height || (height > 0 && len(w.cachedGrid[0]) != width) {
		w.cachedGrid = make([][]int, height)
		for y := 0; y < height; y++ {
			w.cachedGrid[y] = make([]int, width)
		}
	}

	grid := w.cachedGrid
	for y := 0; y < height; y++ {
		cy := minY + y/settings.ChunkHeight
		inChunkY := y % settings.ChunkHeight
		for x := 0; x < width; x++ {
			cx := minX + x/settings.ChunkWidth
			inChunkX := x % settings.ChunkWidth
			coord := ChunkCoord{X: cx, Y: cy}
			chunk, ok := w.Chunks[coord]
			if !ok || len(chunk) == 0 {
				grid[y][x] = int(block.Air)
				continue
			}
			grid[y][x] = int(chunk[inChunkY][inChunkX])
		}
	}

	// Cache the results
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

	coord := ChunkCoord{X: chunkX, Y: chunkY}
	chunk, exists := w.Chunks[coord]
	if !exists {
		return block.Air // Return air for non-existent chunks
	}

	// Bounds check
	if inChunkY < 0 || inChunkY >= settings.ChunkHeight || inChunkX < 0 || inChunkX >= settings.ChunkWidth {
		return block.Air
	}

	return chunk[inChunkY][inChunkX]
}

// SetBlockAt sets the block type at the given world coordinates
func (w *World) SetBlockAt(blockX, blockY int, blockType block.BlockType) bool {
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

	coord := ChunkCoord{X: chunkX, Y: chunkY}
	chunk, exists := w.Chunks[coord]
	if !exists {
		return false // Cannot modify non-existent chunks
	}

	// Bounds check
	if inChunkY < 0 || inChunkY >= settings.ChunkHeight || inChunkX < 0 || inChunkX >= settings.ChunkWidth {
		return false
	}

	// Set the block
	chunk[inChunkY][inChunkX] = blockType
	w.Chunks[coord] = chunk // Update the chunk in the map

	// Instead of marking entire grid dirty, update just this block in the cached grid
	w.updateCachedGridBlock(blockX, blockY, blockType)

	return true
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

	// Check if any entity would collide with the new block
	if w.wouldBlockCollideWithEntity(blockX, blockY) {
		return false // Cannot place block where an entity is present
	}

	return w.SetBlockAt(blockX, blockY, blockType)
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
