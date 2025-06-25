package world

import (
	"math/rand"

	"github.com/KdntNinja/webcraft/engine/block"
)

type BlockType int

const (
	Air BlockType = iota
	Grass
	Dirt
	Stone
)

type Chunk [block.ChunkHeight][block.ChunkWidth]BlockType

type World [][]Chunk // [vertical][horizontal] for multiple columns

var surfaceHeights = make(map[int]int)

func getSurfaceHeight(x int) int {
	h, ok := surfaceHeights[x]
	if ok {
		return h
	}
	prev := block.ChunkHeight / 2
	if x > 0 {
		prev = getSurfaceHeight(x - 1)
	}
	// Random walk for demo
	change := rand.Intn(3) - 1 // -1, 0, or +1
	h = prev + change
	if h < 4 {
		h = 4
	}
	if h > block.ChunkHeight-4 {
		h = block.ChunkHeight - 4
	}
	surfaceHeights[x] = h
	return h
}

func GenerateChunk(chunkX, chunkY int) Chunk {
	var chunk Chunk
	for y := 0; y < block.ChunkHeight; y++ {
		for x := 0; x < block.ChunkWidth; x++ {
			globalX := chunkX*block.ChunkWidth + x
			globalY := chunkY*block.ChunkHeight + y
			surface := getSurfaceHeight(globalX)
			if globalY < surface {
				chunk[y][x] = Air
			} else if globalY == surface {
				chunk[y][x] = Grass
			} else if globalY > surface && globalY < surface+3 {
				chunk[y][x] = Dirt
			} else {
				chunk[y][x] = Stone
			}
		}
	}
	return chunk
}

// GenerateWorld generates a 2D slice of chunks: [vertical][horizontal]
func GenerateWorld(numChunksY int, centerChunkX int) [][]Chunk {
	width := 5 // 2 chunks left, 1 center, 2 right
	world := make([][]Chunk, numChunksY)
	for cy := 0; cy < numChunksY; cy++ {
		world[cy] = make([]Chunk, width)
		for cx := 0; cx < width; cx++ {
			chunkX := centerChunkX + cx - 2
			world[cy][cx] = GenerateChunk(chunkX, cy)
		}
	}
	return world
}
