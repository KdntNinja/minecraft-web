package terrain

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/generation/noise"
)

// Get underground block based on depth, biome, and noise
func getUndergroundBlock(depthFromSurface int, biome noise.BiomeData, terrainNoise *noise.PerlinNoise, globalX, globalY int) block.BlockType {
	x, y := float64(globalX), float64(globalY)

	if depthFromSurface <= 3 {
		// Surface soil layer - simplified to just dirt for all biomes
		return block.Dirt
	} else if depthFromSurface <= 8 {
		// Transition layer with some variation
		transitionNoise := terrainNoise.Noise2D(x*0.05, y*0.05)
		if transitionNoise > 0.4 {
			return block.Dirt
		}
		return block.Stone
	} else if depthFromSurface <= 20 {
		// Mixed stone and dirt layer
		mixNoise := terrainNoise.Noise2D(x*0.03, y*0.03)
		if mixNoise > 0.6 {
			return block.Dirt
		} else if mixNoise > 0.2 {
			return block.Stone
		} else { // Some variety in stone types
			switch biome.Type {
			case 2: // DesertBiome - use clay instead of sand
				return block.Clay
			case 3: // MountainBiome
				return block.Stone // Replace Granite with Stone
			default:
				return block.Stone
			}
		}
	} else {
		// Deep stone with occasional variation
		stoneVariation := terrainNoise.FractalNoise2D(x*0.03, y*0.03, 2, 0.02, 1.0, 0.5)

		if stoneVariation > 0.7 {
			// Stone variants based on biome
			switch biome.Type {
			case 3: // MountainBiome
				return block.Stone // Replace Granite with Stone
			case 7: // OceanBiome (if underground)
				return block.Stone // Replace Marble with Stone
			case 2: // DesertBiome - use clay instead of sand
				return block.Clay
			default:
				return block.Stone
			}
		} else if stoneVariation > 0.3 {
			return block.Stone
		} else {
			// Add some clay for variety instead of limestone
			return block.Clay
		}
	}
}
