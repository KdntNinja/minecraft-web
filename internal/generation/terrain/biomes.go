package terrain

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/generation/noise"
)

// GetBiomeAt is a standalone wrapper for the noise package method
// This maintains backward compatibility with the old biome system
func GetBiomeAt(terrainNoise *noise.PerlinNoise, x float64) int {
	biome := terrainNoise.GetBiomeAt(x)
	return int(biome.Type)
}

// GetBiomeTerrainHeight is a standalone wrapper for the noise package method
func GetBiomeTerrainHeight(terrainNoise *noise.PerlinNoise, x float64, biome int) float64 {
	biomeData := noise.BiomeData{Type: noise.BiomeType(biome)}
	return terrainNoise.GetBiomeTerrainHeight(x, biomeData)
}

// getSurfaceBlockByID maps biome ID to surface block
func getSurfaceBlockByID(biome int) block.BlockType {
	switch biome {
	case 0: // PlainseBiome
		return block.Grass
	case 1: // ForestBiome
		return block.Grass
	case 2: // MountainBiome (renumbered from 3)
		return block.Stone
	case 3: // JungleBiome (renumbered from 6)
		return block.Grass
	case 4: // OceanBiome (renumbered from 7)
		return block.Water
	default:
		return block.Grass
	}
}

// getTerrariaStyleSurfaceBlock returns biome-appropriate surface blocks with variation
func getTerrariaStyleSurfaceBlock(biome noise.BiomeData, terrainNoise *noise.PerlinNoise, globalX int) block.BlockType {
	// Add some surface variation within biomes
	surfaceVariation := terrainNoise.Noise1D(float64(globalX)*0.03 + 500)

	switch biome.Type {
	case 0: // PlainseBiome
		if surfaceVariation > 0.7 {
			return block.Dirt // Exposed dirt patches
		}
		return block.Grass
	case 1: // ForestBiome
		return block.Grass
	case 2: // MountainBiome
		if biome.Temperature < -0.2 {
			return block.Dirt // Cold mountain areas
		} else if biome.Elevation > 0.5 || surfaceVariation > 0.4 {
			return block.Stone // Rocky mountain peaks
		}
		return block.Grass
	case 3: // JungleBiome
		if surfaceVariation > 0.6 {
			return block.Dirt // Jungle dirt patches
		}
		return block.Grass
	case 4: // OceanBiome
		return block.Dirt // Ocean floor
	default:
		return block.Grass
	}
}

// getTerrariaStyleSoilBlock returns biome-appropriate soil blocks
func getTerrariaStyleSoilBlock(biome noise.BiomeData, terrainNoise *noise.PerlinNoise, globalX, globalY int) block.BlockType {
	soilNoise := terrainNoise.FractalNoise2D(float64(globalX)*0.05, float64(globalY)*0.05, 2, 0.08, 1.0, 0.5)

	// Simplified soil logic for remaining biomes
	if soilNoise > 0.7 {
		return block.Clay
	}
	return block.Dirt
}

// getSurfaceBlock returns the basic surface block for a biome (legacy function)
func getSurfaceBlock(biome noise.BiomeData) block.BlockType {
	switch biome.Type {
	case 0: // PlainseBiome
		return block.Grass
	case 1: // ForestBiome
		return block.Grass
	case 2: // MountainBiome
		if biome.Temperature < -0.2 {
			return block.Dirt // Cold mountain areas
		} else if biome.Elevation > 0.5 {
			return block.Stone
		}
		return block.Grass
	case 3: // JungleBiome
		return block.Grass
	case 4: // OceanBiome
		return block.Water
	default:
		return block.Grass
	}
}
