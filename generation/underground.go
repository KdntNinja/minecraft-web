package generation

import (
	"math/rand"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/settings"
)

// GetUndergroundBlock determines the block type for underground positions
func GetUndergroundBlock(worldX, worldY, surfaceHeight int, rng *rand.Rand) coretypes.BlockType {
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
				return coretypes.CopperOre
			case 2:
				return coretypes.IronOre
			case 3:
				return coretypes.GoldOre
			default:
				return coretypes.Stone
			}
		}
	}

	// Shallow underground (already handled in surface.go for <= 4)
	if depthFromSurface <= 8 { // Increased from 4 to 8 for thicker above-ground and shallow layers
		noiseVal := terrainNoise.Noise2D(float64(worldX)/10.0, float64(worldY)/10.0)
		if noiseVal > 0.3 {
			return coretypes.Clay
		} else {
			return coretypes.Dirt
		}
	} else if worldY < settings.ChunkHeight*settings.WorldChunksY-15 {
		// Interwoven stone layers using noise for each type
		graniteNoise := terrainNoise.Noise2D(float64(worldX)/22.0+100, float64(worldY)/22.0+100)
		andesiteNoise := terrainNoise.Noise2D(float64(worldX)/22.0+200, float64(worldY)/22.0+200)
		dioriteNoise := terrainNoise.Noise2D(float64(worldX)/22.0+300, float64(worldY)/22.0+300)
		slateNoise := terrainNoise.Noise2D(float64(worldX)/22.0+400, float64(worldY)/22.0+400)
		clayPockets := terrainNoise.Noise2D(float64(worldX)/8.0+500, float64(worldY)/8.0+500)

		if clayPockets > 0.8 {
			return coretypes.Clay
		}
		if graniteNoise > 0.45 {
			return coretypes.Granite
		}
		if andesiteNoise > 0.45 {
			return coretypes.Andesite
		}
		if dioriteNoise > 0.45 {
			return coretypes.Diorite
		}
		if slateNoise > 0.45 {
			return coretypes.Slate
		}
		// Ash pockets in deeper stone
		stoneVariation := terrainNoise.Noise2D(float64(worldX)/15.0, float64(worldY)/15.0)
		if stoneVariation > 0.4 && depthFromSurface > 15 && rng.Float64() < 0.10 {
			return coretypes.Ash
		}
		return coretypes.Stone
	} else if worldY < settings.ChunkHeight*settings.WorldChunksY-5 {
		// Deep underground transition zone
		deepNoise := terrainNoise.Noise2D(float64(worldX)/20.0, float64(worldY)/20.0)

		if deepNoise > 0.5 {
			return coretypes.Ash
		} else if deepNoise < -0.3 && rng.Float64() < 0.2 {
			// Rare water pockets in deep areas
			return coretypes.Water
		} else {
			return coretypes.Stone
		}
	} else {
		// Deepest levels - hellstone with some variation
		hellNoise := terrainNoise.Noise2D(float64(worldX)/12.0, float64(worldY)/12.0)

		if hellNoise > 0.3 {
			return coretypes.Hellstone
		} else if hellNoise < -0.4 {
			return coretypes.Ash
		} else {
			return coretypes.Stone
		}
	}
}
