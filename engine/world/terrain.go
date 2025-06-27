package world

import (
	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/noise"
)

var (
	surfaceHeights = make(map[int]int)
	chunkCache     = make(map[string]block.Chunk) // Cache chunks to avoid regeneration

	// Enhanced noise generator with crypto-random seeds
	enhancedNoise *noise.EnhancedNoiseGenerator

	// Individual noise instances for backward compatibility
	surfaceNoise    *noise.SimplexNoise // For surface terrain height
	dirtNoise       *noise.SimplexNoise // For dirt layer transitions
	stoneNoise      *noise.SimplexNoise // For stone layer variations
	caveNoise       *noise.SimplexNoise // For cave generation
	oreNoise        *noise.SimplexNoise // For ore distribution
	biomeNoise      *noise.SimplexNoise // For biome distribution
	underworldNoise *noise.SimplexNoise // For underworld generation
)

func initNoiseGenerators() {
	if enhancedNoise == nil {
		// Use enhanced noise generator with crypto-random seed
		enhancedNoise = noise.NewEnhancedNoiseGenerator()

		// Clear caches for fresh world generation
		surfaceHeights = make(map[int]int)
		chunkCache = make(map[string]block.Chunk)

		// Assign individual noise instances for compatibility
		surfaceNoise = enhancedNoise.Primary
		dirtNoise = enhancedNoise.Secondary
		stoneNoise = enhancedNoise.Detail
		caveNoise = enhancedNoise.Cave
		oreNoise = enhancedNoise.Ore
		biomeNoise = enhancedNoise.Biome
		underworldNoise = enhancedNoise.Structure
	}
}

func getSurfaceHeight(x int) int {
	h, ok := surfaceHeights[x]
	if ok {
		return h
	}
	initNoiseGenerators()

	// Use enhanced terrain generation for better variety
	baseHeight := 12 // Base terrain height

	// Generate sophisticated terrain using enhanced algorithms
	terrainNoise := surfaceNoise.AdvancedHybridTerrainNoise(float64(x))

	// Enhanced biome-based terrain variations
	biome := getBiome(x)
	var heightVariation int
	var biomeModifier float64

	switch biome {
	case DesertBiome:
		// Desert: flatter with occasional dunes
		biomeModifier = surfaceNoise.PlainsTerrainNoise(float64(x)) * 0.5
		heightVariation = int((terrainNoise + biomeModifier) * 3)
	case SnowBiome:
		// Snow: dramatic mountainous terrain
		biomeModifier = surfaceNoise.MountainousTerrainNoise(float64(x)) * 0.7
		heightVariation = int((terrainNoise + biomeModifier) * 8)
	case ClayCanyonBiome:
		// Canyon: sharp ridges and valleys
		biomeModifier = surfaceNoise.RidgedNoise1D(float64(x), 2, 0.01, 1.0) * 0.6
		heightVariation = int((terrainNoise + biomeModifier) * 6)
	default: // Forest
		// Forest: gentle rolling hills with enhanced detail
		biomeModifier = surfaceNoise.EnhancedTerrainNoise(float64(x)) * 0.3
		heightVariation = int((terrainNoise + biomeModifier) * 5)
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

// getBiome determines the biome at a given x coordinate using fast noise
func getBiome(x int) BiomeType {
	initNoiseGenerators()

	// Use simple noise for better performance
	biomeNoise := biomeNoise.Noise1D(float64(x) * 0.005) // Simplified scale

	// Simplified biome thresholds
	if biomeNoise < -0.3 {
		return DesertBiome
	} else if biomeNoise > 0.4 {
		return SnowBiome
	} else if biomeNoise > 0.1 {
		return ClayCanyonBiome
	}
	return ForestBiome // Most common biome
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

// getOreType determines what ore should be placed using enhanced noise
func getOreType(x, y, depthFromSurface int) block.BlockType {
	initNoiseGenerators()

	// Use enhanced ore generation with depth-based rarity
	oreType := oreNoise.EnhancedOreNoise(float64(x), float64(y), depthFromSurface)

	// Map ore type numbers to block types
	switch oreType {
	case 0:
		return block.CopperOre
	case 1:
		return block.IronOre
	case 2:
		return block.SilverOre
	case 3:
		return block.GoldOre
	case 4:
		return block.PlatinumOre
	default:
		return block.Stone // Default to stone
	}
}

// isInCave determines if a position should be a cave using enhanced approach
func isInCave(x, y int) bool {
	initNoiseGenerators()

	// Don't generate caves too close to surface
	if y < 8 {
		return false
	}

	// Use enhanced cave generation for better variety
	caveNoiseValue := caveNoise.EnhancedCaveNoise(float64(x), float64(y))

	// Dynamic cave threshold based on depth for more realistic cave systems
	depth := y - 8
	depthThreshold := -0.4 + float64(depth)*0.01 // More caves deeper down

	return caveNoiseValue < depthThreshold
}

// isUnderworld checks if we're in the underworld layer
func isUnderworld(y, worldHeight int) bool {
	return y > worldHeight-6 // Bottom 6 layers are underworld
}

// ResetWorldGeneration forces regeneration with new random seeds
func ResetWorldGeneration() {
	// Reset enhanced noise generator to nil to force regeneration with new crypto-random seed
	enhancedNoise = nil
	
	// Reset all individual noise generators to nil
	surfaceNoise = nil
	dirtNoise = nil
	stoneNoise = nil
	caveNoise = nil
	oreNoise = nil
	biomeNoise = nil
	underworldNoise = nil

	// Clear all caches
	surfaceHeights = make(map[int]int)
	chunkCache = make(map[string]block.Chunk)
}

// GetWorldSeed returns the current world seed for display or saving
func GetWorldSeed() int64 {
	initNoiseGenerators()
	if enhancedNoise != nil {
		return enhancedNoise.Seed
	}
	return 0
}
