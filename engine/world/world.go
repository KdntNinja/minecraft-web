package world

import (
	"math/rand"

	"github.com/aquilax/go-perlin"

	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/entity"
	"github.com/KdntNinja/webcraft/engine/player"
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

type World struct {
	Blocks   [][]Chunk // [vertical][horizontal] for multiple columns
	Entities entity.Entities
}

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
		Blocks:   blocks,
		Entities: entity.Entities{},
	}
	// Add player entity at center
	px := (len(blocks[0])*block.ChunkWidth/2)*block.TileSize + block.TileSize/2
	playerGlobalX := px / block.TileSize
	// Use ToIntGrid to get the world as a 2D grid
	grid := w.ToIntGrid()
	spawnY := 0
	for y := 0; y < len(grid); y++ {
		if entity.IsSolid(grid, playerGlobalX, y) {
			spawnY = y - 1 // one block above the first solid block
			break
		}
	}
	if spawnY < 0 {
		spawnY = 0 // fallback to top if no solid block found
	}
	py := float64(spawnY*block.TileSize - player.Height/2)
	w.Entities = append(w.Entities, player.NewPlayer(float64(px), py))
	return w
}

// ToIntGrid flattens the world's blocks into a [][]int grid for entity collision
func (w *World) ToIntGrid() [][]int {
	height := len(w.Blocks) * block.ChunkHeight
	width := len(w.Blocks[0]) * block.ChunkWidth
	grid := make([][]int, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
		cy := y / block.ChunkHeight
		inChunkY := y % block.ChunkHeight
		for x := 0; x < width; x++ {
			cx := x / block.ChunkWidth
			inChunkX := x % block.ChunkWidth
			grid[y][x] = int(w.Blocks[cy][cx][inChunkY][inChunkX])
		}
	}
	return grid
}
