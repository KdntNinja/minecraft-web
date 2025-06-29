package generation

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// GetSurfaceBlockType determines the surface block type based on biome
func GetSurfaceBlockType(worldX int) block.BlockType {
	// Always grass for now - could add biome variation later
	return block.Grass
}

// GetShallowUndergroundBlock determines shallow underground block types
func GetShallowUndergroundBlock(worldX, worldY int) block.BlockType {
	terrainNoise := GetTerrainNoise()

	// Biome-based surface variation
	biomeNoise := terrainNoise.Noise2D(float64(worldX)/settings.TreeBiomeNoiseScale, 0)
	noiseVal := terrainNoise.Noise2D(float64(worldX)/settings.TreeNoiseScale, float64(worldY)/settings.TreeNoiseScale)

	if biomeNoise > settings.TreeClayBiomeThresh && noiseVal > settings.TreeClayNoiseThresh {
		return block.Clay // Clay deposits in certain biomes
	} else {
		return block.Dirt // Normal dirt layer under grass
	}
}
