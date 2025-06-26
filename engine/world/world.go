package world

import (
	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/entity"
	"github.com/KdntNinja/webcraft/engine/player"
)

type BlockType int

const (
	Air BlockType = iota
	Grass
	Dirt
	Stone
)

type Chunk [block.ChunkHeight][block.ChunkWidth]BlockType

type World struct {
	Blocks    [][]Chunk // [vertical][horizontal] for multiple columns
	Entities  entity.Entities
	MinChunkX int // Minimum chunk X coordinate in the world grid
	MinChunkY int // Minimum chunk Y coordinate in the world grid
}

// NewWorld constructs a new World instance with generated chunks
func NewWorld(numChunksY int, centerChunkX int) *World {
	width := 5 // 2 chunks left, 1 center, 2 right
	blocks := make([][]Chunk, numChunksY)
	for cy := 0; cy < numChunksY; cy++ {
		blocks[cy] = make([]Chunk, width)
		for cx := 0; cx < width; cx++ {
			chunkX := centerChunkX + cx - 2
			blocks[cy][cx] = GenerateChunk(chunkX, cy)
		}
	}
	w := &World{
		Blocks:    blocks,
		Entities:  entity.Entities{},
		MinChunkX: centerChunkX - 2, // Start with -2 offset
		MinChunkY: 0,                // Start at Y=0
	}
	// Add player entity at center
	px := (len(blocks[0])*block.ChunkWidth/2)*block.TileSize + block.TileSize/2
	playerGlobalX := px / block.TileSize

	// Find surface height more efficiently by looking for the grass block
	// instead of searching from the top
	spawnY := 0
	found := false

	// First, calculate the expected surface height using the same logic as chunk generation
	worldX := playerGlobalX
	surfaceHeight := getSurfaceHeight(worldX)

	// Start searching from around the expected surface height
	searchStart := max(0, surfaceHeight-5)
	searchEnd := min(len(blocks)*block.ChunkHeight, surfaceHeight+10)

	grid := w.ToIntGrid()

	// Look for the surface (grass block) or first solid block
	for y := searchStart; y < searchEnd; y++ {
		if y < len(grid) && playerGlobalX < len(grid[y]) {
			if grid[y][playerGlobalX] == int(Grass) {
				spawnY = y - 1 // One block above grass
				found = true
				break
			} else if entity.IsSolid(grid, playerGlobalX, y) {
				spawnY = y - 1 // One block above first solid block
				found = true
				break
			}
		}
	}

	// If we didn't find a surface, default to a reasonable height
	if !found {
		spawnY = surfaceHeight - 1
	}

	// Ensure spawn position is valid
	if spawnY < 0 {
		spawnY = 0
	}

	// Center player in the block, fully inside
	py := float64(spawnY * block.TileSize)
	w.Entities = append(w.Entities, player.NewPlayer(float64(px), py))
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
				grid[y][x] = int(Air) // Default to air if chunk doesn't exist
				continue
			}

			grid[y][x] = int(w.Blocks[cy][cx][inChunkY][inChunkX])
		}
	}
	return grid
}

// Helper functions for min and max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
