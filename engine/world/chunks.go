package world

import (
	"fmt"

	"github.com/KdntNinja/webcraft/engine/block"
)

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
