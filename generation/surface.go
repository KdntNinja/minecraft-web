package generation

import (
	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/settings"
)

// GetSurfaceBlockType determines the surface block type based on biome
func GetSurfaceBlockType(worldX int) coretypes.BlockType {
	// Always grass for now - could add biome variation later
	return coretypes.Grass
}

// GetShallowUndergroundBlock determines shallow underground block types
func GetShallowUndergroundBlock(worldX, worldY int) coretypes.BlockType {
	terrainNoise := GetTerrainNoise()

	// Biome-based surface variation
	biomeNoise := terrainNoise.Noise2D(float64(worldX)/settings.TreeBiomeNoiseScale, 0)
	noiseVal := terrainNoise.Noise2D(float64(worldX)/settings.TreeNoiseScale, float64(worldY)/settings.TreeNoiseScale)

	if biomeNoise > settings.TreeClayBiomeThresh && noiseVal > settings.TreeClayNoiseThresh {
		return coretypes.Clay // Clay deposits in certain biomes
	} else {
		return coretypes.Dirt // Normal dirt layer under grass
	}
}
