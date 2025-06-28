package terrain

import (
	"fmt"

	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation/noise"
)

// Helper function to generate chunk key
func getChunkKey(chunkX, chunkY int) string {
	// Use a string format to avoid collisions and support negative coordinates
	return fmt.Sprintf("%d,%d", chunkX, chunkY)
}

// getSurfaceHeightRange calculates the min and max surface heights for a chunk column
func getSurfaceHeightRange(chunkX int, terrainNoise *noise.PerlinNoise) (int, int) {
	minHeight := 1000
	maxHeight := -1000

	// Sample a few points across the chunk width to get height range
	for x := 0; x < settings.ChunkWidth; x++ {
		globalX := chunkX*settings.ChunkWidth + x
		biome := terrainNoise.GetBiomeAt(float64(globalX))
		terrainHeight := terrainNoise.GetBiomeTerrainHeight(float64(globalX), biome)

		baseHeight := 50.0
		heightVariation := terrainHeight * 20.0
		surfaceHeight := int(baseHeight + heightVariation)

		if surfaceHeight < minHeight {
			minHeight = surfaceHeight
		}
		if surfaceHeight > maxHeight {
			maxHeight = surfaceHeight
		}
	}

	return minHeight, maxHeight
}
