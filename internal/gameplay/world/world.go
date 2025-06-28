package world

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
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

	// Find the surface height at the center position
	surfaceY := FindSurfaceHeight(centerBlockX, blocks)

	// Spawn player 2 blocks above the surface for safety
	spawnY := surfaceY - 2

	// Ensure spawn position is reasonable
	if spawnY < 0 {
		spawnY = 0
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
