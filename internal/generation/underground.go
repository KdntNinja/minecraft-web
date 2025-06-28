package generation

import (
	"math/rand"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// GetUndergroundBlock determines the block type for underground positions
func GetUndergroundBlock(worldX, worldY, surfaceHeight int, rng *rand.Rand) block.BlockType {
	depthFromSurface := worldY - surfaceHeight
	terrainNoise := GetTerrainNoise()

	// Check for ore veins with enhanced vein generation
	oreType := GetOreType(worldX, worldY)
	if oreType > 0 {
		// Base ore chance from settings, but increased for visibility
		oreChance := settings.OreVeinChance * 4.0 // Significantly increased

		// Check if this position extends an existing ore vein (makes ores cluster)
		if IsOreVeinExtension(worldX, worldY, oreType) {
			oreChance *= 3.0 // Much more likely to place ore near other ore
		}

		if rng.Float64() < oreChance {
			switch oreType {
			case 1:
				return block.CopperOre
			case 2:
				return block.IronOre
			case 3:
				return block.GoldOre
			default:
				return block.Stone
			}
		}
	}

	// Shallow underground (already handled in surface.go for <= 4)
	if depthFromSurface <= 4 {
		noiseVal := terrainNoise.Noise2D(float64(worldX)/10.0, float64(worldY)/10.0)
		if noiseVal > 0.3 {
			return block.Clay
		} else {
			return block.Dirt
		}
	} else if worldY < settings.ChunkHeight*settings.WorldChunksY-15 {
		// Stone layer with reduced clay and more variety
		stoneVariation := terrainNoise.Noise2D(float64(worldX)/15.0, float64(worldY)/15.0)
		clayPockets := terrainNoise.Noise2D(float64(worldX)/8.0+500, float64(worldY)/8.0+500)

		// Create clay pockets in stone (reduced frequency)
		if clayPockets > 0.8 { // Made much rarer (was 0.6)
			return block.Clay
		} else if stoneVariation > 0.4 && depthFromSurface > 15 {
			// Deeper stone areas can have ash pockets
			if rng.Float64() < 0.10 { // Reduced ash frequency too
				return block.Ash
			} else {
				return block.Stone
			}
		} else {
			return block.Stone
		}
	} else if worldY < settings.ChunkHeight*settings.WorldChunksY-5 {
		// Deep underground transition zone
		deepNoise := terrainNoise.Noise2D(float64(worldX)/20.0, float64(worldY)/20.0)

		if deepNoise > 0.5 {
			return block.Ash
		} else if deepNoise < -0.3 && rng.Float64() < 0.2 {
			// Rare water pockets in deep areas
			return block.Water
		} else {
			return block.Stone
		}
	} else {
		// Deepest levels - hellstone with some variation
		hellNoise := terrainNoise.Noise2D(float64(worldX)/12.0, float64(worldY)/12.0)

		if hellNoise > 0.3 {
			return block.Hellstone
		} else if hellNoise < -0.4 {
			return block.Ash
		} else {
			return block.Stone
		}
	}
}
