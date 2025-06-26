package world

import (
	"fmt"

	"github.com/aquilax/go-perlin"

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
	Blocks   [][]Chunk // [vertical][horizontal] for multiple columns
	Entities entity.Entities
}

var (
	surfaceHeights = make(map[int]int)

	// Different Perlin noise instances for each terrain layer
	surfaceNoise *perlin.Perlin // For surface terrain height
	dirtNoise    *perlin.Perlin // For dirt layer transitions
	stoneNoise   *perlin.Perlin // For stone layer variations

	chunkCache = make(map[string]Chunk) // Cache chunks to avoid regeneration
)

func initNoiseGenerators() {
	if surfaceNoise == nil {
		// Surface terrain - smoother, larger features
		surfaceNoise = perlin.NewPerlin(2, 2, 3, 12345)

		// Dirt layer - medium frequency transitions
		dirtNoise = perlin.NewPerlin(3, 3, 4, 67890)

		// Stone layer - higher frequency, more chaotic
		stoneNoise = perlin.NewPerlin(4, 4, 5, 54321)
	}
}

func getSurfaceHeight(x int) int {
	h, ok := surfaceHeights[x]
	if ok {
		return h
	}
	initNoiseGenerators()

	// Use surface noise for main terrain height with larger scale features
	scale := 0.05 // Lower frequency for smoother terrain
	noise := surfaceNoise.Noise1D(float64(x) * scale)
	height := int((noise+1)*0.5*float64(block.ChunkHeight-8)) + 4
	surfaceHeights[x] = height
	return height
}

func GenerateChunk(chunkX, chunkY int) Chunk {
	// Check cache first to avoid regenerating identical chunks
	cacheKey := fmt.Sprintf("%d,%d", chunkX, chunkY)
	if cached, exists := chunkCache[cacheKey]; exists {
		return cached
	}

	var chunk Chunk
	// Optimize: calculate surface heights for entire chunk width at once
	surfaces := make([]int, block.ChunkWidth)
	for x := 0; x < block.ChunkWidth; x++ {
		globalX := chunkX*block.ChunkWidth + x
		surfaces[x] = getSurfaceHeight(globalX)
	}

	// Generate chunk blocks using pre-calculated surface heights and multiple noise layers
	for y := 0; y < block.ChunkHeight; y++ {
		globalY := chunkY*block.ChunkHeight + y
		for x := 0; x < block.ChunkWidth; x++ {
			globalX := chunkX*block.ChunkWidth + x
			surface := surfaces[x]

			if globalY < surface {
				chunk[y][x] = Air
			} else if globalY == surface {
				chunk[y][x] = Grass
			} else {
				// Use different noise algorithms for natural layer transitions
				depthFromSurface := globalY - surface

				// Dirt layer transition using dirt noise
				dirtScale := 0.12
				dirtNoise2D := dirtNoise.Noise2D(float64(globalX)*dirtScale, float64(globalY)*dirtScale)
				dirtThickness := 3 + int(dirtNoise2D*2) // 1-5 blocks thick

				if depthFromSurface <= dirtThickness {
					chunk[y][x] = Dirt
				} else {
					// Stone layer with some variation using stone noise
					stoneScale := 0.15
					stoneNoise2D := stoneNoise.Noise2D(float64(globalX)*stoneScale, float64(globalY)*stoneScale)

					// Occasionally create air pockets (caves) in stone
					if stoneNoise2D < -0.6 && depthFromSurface > 5 {
						chunk[y][x] = Air
					} else {
						chunk[y][x] = Stone
					}
				}
			}
		}
	}

	// Cache the generated chunk
	chunkCache[cacheKey] = chunk
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
	grid := w.ToIntGrid()
	spawnY := 0
	for y := 0; y < len(grid); y++ {
		if entity.IsSolid(grid, playerGlobalX, y) {
			spawnY = y - 1 // one block above the first solid block
			break
		}
	}
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
