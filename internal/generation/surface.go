package generation

import "github.com/KdntNinja/webcraft/internal/core/engine/block"

// GetSurfaceBlockType determines the surface block type based on biome
func GetSurfaceBlockType(worldX int) block.BlockType {
	// Always grass for now - could add biome variation later
	return block.Grass
}

// GetShallowUndergroundBlock determines shallow underground block types
func GetShallowUndergroundBlock(worldX, worldY int) block.BlockType {
	terrainNoise := GetTerrainNoise()

	// Biome-based surface variation
	biomeNoise := terrainNoise.Noise2D(float64(worldX)/80.0, 0)
	noiseVal := terrainNoise.Noise2D(float64(worldX)/10.0, float64(worldY)/10.0)

	if biomeNoise > 0.4 && noiseVal > 0.2 {
		return block.Clay // Clay deposits in certain biomes
	} else {
		return block.Dirt // Normal dirt layer under grass
	}
}
