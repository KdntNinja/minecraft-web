package world

import (
	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/noise"
)

var (
	surfaceHeights = make(map[int]int)
	chunkCache     = make(map[string]block.Chunk) // Cache chunks to avoid regeneration

	// Different noise instances for each terrain layer
	surfaceNoise    *noise.SimplexNoise // For surface terrain height
	dirtNoise       *noise.SimplexNoise // For dirt layer transitions
	stoneNoise      *noise.SimplexNoise // For stone layer variations
	caveNoise       *noise.SimplexNoise // For cave generation
	oreNoise        *noise.SimplexNoise // For ore distribution
	biomeNoise      *noise.SimplexNoise // For biome distribution
	underworldNoise *noise.SimplexNoise // For underworld generation
)

func initNoiseGenerators() {
	if surfaceNoise == nil {
		// Surface terrain - smoother, larger features
		surfaceNoise = noise.NewSimplexNoise(12345)

		// Dirt layer - medium frequency transitions
		dirtNoise = noise.NewSimplexNoise(67890)

		// Stone layer - higher frequency, more chaotic
		stoneNoise = noise.NewSimplexNoise(54321)

		// Cave generation - creates caverns and tunnels
		caveNoise = noise.NewSimplexNoise(98765)

		// Ore distribution - for placing ore veins
		oreNoise = noise.NewSimplexNoise(11111)

		// Biome distribution - determines surface biomes
		biomeNoise = noise.NewSimplexNoise(22222)

		// Underworld generation - for hell layer
		underworldNoise = noise.NewSimplexNoise(33333)
	}
}

func getSurfaceHeight(x int) int {
	h, ok := surfaceHeights[x]
	if ok {
		return h
	}
	initNoiseGenerators()

	// Use surface noise for main terrain height with larger scale features
	scale := 0.02    // Lower frequency for smoother terrain
	baseHeight := 12 // Base terrain height

	// Multi-octave noise for more interesting terrain using fractal noise
	combinedNoise := surfaceNoise.FractalNoise1D(float64(x), 3, scale, 1.0, 0.5)

	// Different terrain variations based on biome
	biome := getBiome(x)
	var heightVariation int

	switch biome {
	case DesertBiome:
		heightVariation = int(combinedNoise * 4) // Flatter terrain
	case SnowBiome:
		// Use ridged noise for mountain peaks
		ridgeNoise := surfaceNoise.RidgedNoise1D(float64(x), 2, scale*0.5, 1.0)
		heightVariation = int((combinedNoise + ridgeNoise*0.5) * 8) // Mountain-like terrain
	case ClayCanyonBiome:
		heightVariation = int(combinedNoise * 6) // Canyon-like terrain
	default: // Forest
		heightVariation = int(combinedNoise * 6) // Rolling hills
	}

	height := baseHeight + heightVariation

	// Ensure height is always within reasonable bounds regardless of x coordinate
	if height < 3 {
		height = 3
	}
	if height > block.ChunkHeight-2 {
		height = block.ChunkHeight - 2
	}

	surfaceHeights[x] = height
	return height
}

// BiomeType represents different surface biomes
type BiomeType int

const (
	ForestBiome BiomeType = iota
	DesertBiome
	SnowBiome
	ClayCanyonBiome
)

// getBiome determines the biome at a given x coordinate
func getBiome(x int) BiomeType {
	initNoiseGenerators()

	scale := 0.01 // Very low frequency for large biome areas
	noise := biomeNoise.Noise1D(float64(x) * scale)

	if noise < -0.3 {
		return DesertBiome
	} else if noise > 0.4 {
		return SnowBiome
	} else if noise > 0.1 && noise < 0.3 {
		return ClayCanyonBiome
	}
	return ForestBiome
}

// getSurfaceBlock returns the appropriate surface block for the biome
func getSurfaceBlock(x int) block.BlockType {
	biome := getBiome(x)

	switch biome {
	case DesertBiome:
		return block.Sand
	case SnowBiome:
		return block.Snow
	case ClayCanyonBiome:
		return block.Clay
	default:
		return block.Grass
	}
}

// getOreType determines what ore should be placed based on depth and noise
func getOreType(x, y, depthFromSurface int) block.BlockType {
	initNoiseGenerators()

	oreScale := 0.1
	oreNoise2D := oreNoise.Noise2D(float64(x)*oreScale, float64(y)*oreScale)

	// Different ores at different depths with rarity
	if depthFromSurface > 15 { // Deep underground
		if oreNoise2D > 0.7 {
			return block.GoldOre
		} else if oreNoise2D > 0.6 {
			return block.SilverOre
		} else if depthFromSurface > 25 && oreNoise2D > 0.5 {
			return block.PlatinumOre // Very deep platinum
		}
	}

	if depthFromSurface > 10 { // Medium depth
		if oreNoise2D > 0.75 {
			return block.IronOre
		}
	}

	if depthFromSurface > 5 { // Shallow underground
		if oreNoise2D > 0.8 {
			return block.CopperOre
		}
	}

	return block.Stone // Default to stone
}

// isInCave determines if a position should be a cave
func isInCave(x, y int) bool {
	initNoiseGenerators()

	caveScale := 0.08
	caveNoise2D := caveNoise.Noise2D(float64(x)*caveScale, float64(y)*caveScale)

	// Add secondary cave layer for more complex cave systems
	caveNoise2 := caveNoise.Noise2D(float64(x)*caveScale*2, float64(y)*caveScale*2) * 0.5

	combinedCaveNoise := caveNoise2D + caveNoise2

	// Create caves only below surface and not too shallow
	// Use different thresholds for different depths to create varied cave systems
	if y < 8 {
		return false // No caves near surface
	}

	if y > 30 { // Deep caves - larger caverns
		return combinedCaveNoise < -0.4
	} else { // Medium depth caves - smaller tunnels
		return combinedCaveNoise < -0.5
	}
}

// isUnderworld checks if we're in the underworld layer
func isUnderworld(y, worldHeight int) bool {
	return y > worldHeight-6 // Bottom 6 layers are underworld
}
