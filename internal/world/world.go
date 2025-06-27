package world

import (
	"github.com/KdntNinja/webcraft/internal/block"
	"github.com/KdntNinja/webcraft/internal/entity"
	"github.com/KdntNinja/webcraft/internal/player"
)

type World struct {
	Blocks    [][]block.Chunk // [vertical][horizontal] for multiple columns
	Entities  entity.Entities
	MinChunkX int // Minimum chunk X coordinate in the world grid
	MinChunkY int // Minimum chunk Y coordinate in the world grid
}

// NewWorld constructs a new World instance with generated chunks
func NewWorld(numChunksY int, centerChunkX int) *World {
	width := 7 // 3 chunks left, 1 center, 3 right
	blocks := make([][]block.Chunk, numChunksY)
	for cy := 0; cy < numChunksY; cy++ {
		blocks[cy] = make([]block.Chunk, width)
		for cx := 0; cx < width; cx++ {
			chunkX := centerChunkX + cx - 3 // -3 to +3 relative to center
			blocks[cy][cx] = GenerateChunk(chunkX, cy)
		}
	}
	w := &World{
		Blocks:    blocks,
		Entities:  entity.Entities{},
		MinChunkX: centerChunkX - 3, // Start with -3 offset
		MinChunkY: 0,                // Start at Y=0
	}
	// Add player entity at center
	centerChunkCol := len(blocks[0]) / 2                                 // Get center chunk column
	centerBlockX := centerChunkCol*block.ChunkWidth + block.ChunkWidth/2 // Center of center chunk
	px := float64(centerBlockX * block.TileSize)

	// Convert to grid coordinates to find spawn position
	grid := w.ToIntGrid()

	// Find the surface by searching from top down to find first solid block (surface)
	spawnY := 0
	found := false

	// Search from top of world down to find the first solid block (surface)
	for y := 0; y < len(grid); y++ {
		if centerBlockX < len(grid[y]) {
			if entity.IsSolid(grid, centerBlockX, y) {
				// Found the first solid block from top, this is the surface
				// Spawn player directly above it with enough clearance
				spawnY = y - 2 // Two blocks above the surface block for player height
				found = true
				break
			}
		}
	}

	// If no surface found, spawn at a reasonable default height near top
	if !found {
		spawnY = 10 // Spawn near the top of the world
	}

	// Ensure spawn position is valid (not negative and has clearance)
	if spawnY < 0 {
		spawnY = 0
	}

	// Verify there's clear air space for the player (player is 2 blocks tall)
	// Move spawn position up if needed to ensure clear space
	for spawnY >= 0 && centerBlockX < len(grid[0]) {
		// Check if current position and position above are clear
		if spawnY < len(grid) && spawnY+1 < len(grid) &&
			!entity.IsSolid(grid, centerBlockX, spawnY) &&
			!entity.IsSolid(grid, centerBlockX, spawnY+1) {
			break // Found clear air space for player
		}
		spawnY-- // Move up one block
		if spawnY < 0 {
			spawnY = 0 // Don't go below world
			break
		}
	}

	py := float64(spawnY * block.TileSize)
	w.Entities = append(w.Entities, player.NewPlayer(px, py))
	return w
}

// ToIntGrid flattens the world's blocks into a [][]int grid for entity collision
func (w *World) ToIntGrid() [][]int {
	if len(w.Blocks) == 0 || len(w.Blocks[0]) == 0 {
		return [][]int{}
	}

	height := len(w.Blocks) * block.ChunkHeight
	width := len(w.Blocks[0]) * block.ChunkWidth
	grid := make([][]int, height)

	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
		cy := y / block.ChunkHeight
		inChunkY := y % block.ChunkHeight

		// Bounds check for cy
		if cy >= len(w.Blocks) {
			continue
		}

		for x := 0; x < width; x++ {
			cx := x / block.ChunkWidth
			inChunkX := x % block.ChunkWidth

			// Bounds check for cx
			if cx >= len(w.Blocks[cy]) {
				grid[y][x] = int(block.Air) // Default to air if chunk doesn't exist
				continue
			}

			grid[y][x] = int(w.Blocks[cy][cx][inChunkY][inChunkX])
		}
	}
	return grid
}
