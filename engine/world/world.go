package world

import (
	"math/rand"

	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/aquilax/go-perlin"
)

const (
	seed       = 100
	smoothness = 100.0
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

var (
	surfaceHeights = make(map[int]int)
	perlinInstance *perlin.Perlin
)

func getPerlin() *perlin.Perlin {
	if perlinInstance == nil {
		perlinInstance = perlin.NewPerlin(2, 2, 3, rand.Int63())
	}
	return perlinInstance
}

func getSurfaceHeight(x int) int {
	h, ok := surfaceHeights[x]
	if ok {
		return h
	}
	// Use Perlin noise for smoother terrain
	scale := 0.08 // Adjust for frequency
	noise := getPerlin().Noise1D(float64(x) * scale)
	height := int((noise+1)*0.5*float64(block.ChunkHeight-8)) + 4
	surfaceHeights[x] = height
	return height
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
