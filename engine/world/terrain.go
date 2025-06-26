package world

import (
	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/aquilax/go-perlin"
)

var (
	surfaceHeights = make(map[int]int)
	chunkCache     = make(map[string]Chunk) // Cache chunks to avoid regeneration

	// Different Perlin noise instances for each terrain layer
	surfaceNoise *perlin.Perlin // For surface terrain height
	dirtNoise    *perlin.Perlin // For dirt layer transitions
	stoneNoise   *perlin.Perlin // For stone layer variations
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
	// Ensure height is always within reasonable bounds regardless of x coordinate
	height := int((noise+1)*0.5*float64(block.ChunkHeight-8)) + 4

	// Clamp height to reasonable values
	if height < 2 {
		height = 2
	}
	if height >= block.ChunkHeight-2 {
		height = block.ChunkHeight - 3
	}

	surfaceHeights[x] = height
	return height
}
