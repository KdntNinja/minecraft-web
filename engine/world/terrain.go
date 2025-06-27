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

	// Use hybrid terrain generation combining Minecraft and Terraria styles
	baseHeight := 12 // Base terrain height

	// Generate hybrid terrain
	terrainNoise := surfaceNoise.HybridTerrainNoise(float64(x))

	// Different terrain variations based on biome
	biome := getBiome(x)
	var heightVariation int

	switch biome {
	case DesertBiome:
		// Desert: Minecraft-style flat with some dunes
		heightVariation = int(terrainNoise * 3) // Flatter terrain
	case SnowBiome:
		// Snow: More Terraria-style mountains with ridges
		ridgeNoise := surfaceNoise.RidgedNoise1D(float64(x), 2, 0.015, 1.0)
		heightVariation = int(terrainNoise*5 + ridgeNoise*3) // Mountain-like terrain
	case ClayCanyonBiome:
		// Clay Canyon: Mix of both styles
		heightVariation = int(terrainNoise * 4) // Canyon-like terrain
	default: // Forest
		// Forest: Minecraft-style rolling hills
		heightVariation = int(terrainNoise * 4) // Gentler rolling hills
	}

	height := baseHeight + heightVariation

	// Ensure height is always within reasonable bounds
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

// getBiome determines the biome at a given x coordinate using Terraria-like distribution
func getBiome(x int) BiomeType {
	initNoiseGenerators()

	// Use Terraria-style biome noise for more natural biome distribution
	biomeNoise := biomeNoise.TerrariaBiomeNoise(float64(x))

	// Terraria-like biome thresholds with more forest (default biome)
	if biomeNoise < -0.4 {
		return DesertBiome
	} else if biomeNoise > 0.5 {
		return SnowBiome
	} else if biomeNoise > 0.2 && biomeNoise < 0.35 {
		return ClayCanyonBiome
	}
	return ForestBiome // Most common biome like in Terraria
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

// shouldPlaceTree determines if a tree should be placed at this x coordinate
func shouldPlaceTree(x int) bool {
	initNoiseGenerators()

	biome := getBiome(x)
	if biome != ForestBiome {
		return false // Only place trees in forest biome for now
	}

	// Simple tree placement using basic noise
	treeNoise := surfaceNoise.Noise1D(float64(x) * 0.1)

	// Trees need adequate spacing and some randomness
	return treeNoise > 0.3 && (x%6 == 0 || x%7 == 0)
}

// getOreType determines what ore should be placed based on depth and Terraria-like noise
func getOreType(x, y, depthFromSurface int) block.BlockType {
	initNoiseGenerators()

	// Use Terraria-style ore generation with different patterns for each ore

	// Copper ore - most common, shallow
	if depthFromSurface > 3 {
		copperNoise := oreNoise.TerrariaOreNoise(float64(x), float64(y), 0)
		if copperNoise > 0 && depthFromSurface < 20 {
			return block.CopperOre
		}
	}

	// Iron ore - medium depth, medium rarity
	if depthFromSurface > 8 {
		ironNoise := oreNoise.TerrariaOreNoise(float64(x), float64(y), 1)
		if ironNoise > 0 && depthFromSurface < 30 {
			return block.IronOre
		}
	}

	// Silver ore - deeper, less common
	if depthFromSurface > 12 {
		silverNoise := oreNoise.TerrariaOreNoise(float64(x), float64(y), 2)
		if silverNoise > 0 && depthFromSurface < 35 {
			return block.SilverOre
		}
	}

	// Gold ore - deep, rare
	if depthFromSurface > 15 {
		goldNoise := oreNoise.TerrariaOreNoise(float64(x), float64(y), 3)
		if goldNoise > 0 {
			return block.GoldOre
		}
	}

	// Platinum ore - very deep, very rare
	if depthFromSurface > 25 {
		platinumNoise := oreNoise.TerrariaOreNoise(float64(x), float64(y), 4)
		if platinumNoise > 0 {
			return block.PlatinumOre
		}
	}

	return block.Stone // Default to stone
}

// isInCave determines if a position should be a cave using hybrid approach
func isInCave(x, y int) bool {
	initNoiseGenerators()

	// Don't generate caves too close to surface
	if y < 6 {
		return false
	}

	// Use hybrid cave noise combining Minecraft and Terraria styles
	caveNoiseValue := caveNoise.HybridCaveNoise(float64(x), float64(y))

	// Different cave thresholds for different depths
	var threshold float64
	if y > 40 { // Deep underground - larger caverns
		threshold = -0.3
	} else if y > 20 { // Medium depth - medium caves
		threshold = -0.4
	} else { // Shallow - small tunnels
		threshold = -0.5
	}

	return caveNoiseValue < threshold
}

// isUnderworld checks if we're in the underworld layer
func isUnderworld(y, worldHeight int) bool {
	return y > worldHeight-6 // Bottom 6 layers are underworld
}
