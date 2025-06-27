package world

import (
	"time"

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
		// Generate random seed based on current time for unique worlds each run
		baseSeed := time.Now().UnixNano()

		// Clear caches for fresh world generation
		surfaceHeights = make(map[int]int)
		chunkCache = make(map[string]block.Chunk)

		// Surface terrain - smoother, larger features
		surfaceNoise = noise.NewSimplexNoise(baseSeed)

		// Dirt layer - medium frequency transitions
		dirtNoise = noise.NewSimplexNoise(baseSeed + 1000)

		// Stone layer - higher frequency, more chaotic
		stoneNoise = noise.NewSimplexNoise(baseSeed + 2000)

		// Cave generation - organic cave systems
		caveNoise = noise.NewSimplexNoise(baseSeed + 3000)

		// Ore distribution - rare mineral veins
		oreNoise = noise.NewSimplexNoise(baseSeed + 4000)

		// Biome distribution - large scale climate zones
		biomeNoise = noise.NewSimplexNoise(baseSeed + 5000)

		// Underworld generation - chaotic hellish terrain
		underworldNoise = noise.NewSimplexNoise(baseSeed + 6000)
	}
}

func getSurfaceHeight(x int) int {
	h, ok := surfaceHeights[x]
	if ok {
		return h
	}
	initNoiseGenerators()

	// Use fast terrain generation for better performance
	baseHeight := 12 // Base terrain height

	// Generate optimized terrain for low-end hardware
	terrainNoise := surfaceNoise.FastTerrainNoise(float64(x))

	// Simplified biome-based terrain variations
	biome := getBiome(x)
	var heightVariation int

	switch biome {
	case DesertBiome:
		heightVariation = int(terrainNoise * 3) // Flatter terrain
	case SnowBiome:
		heightVariation = int(terrainNoise * 6) // Mountain-like terrain
	case ClayCanyonBiome:
		heightVariation = int(terrainNoise * 4) // Canyon-like terrain
	default: // Forest
		heightVariation = int(terrainNoise * 4) // Rolling hills
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

// getOreType determines what ore should be placed using fast noise
func getOreType(x, y, depthFromSurface int) block.BlockType {
	initNoiseGenerators()

	// Use fast ore generation for better performance
	oreNoiseValue := oreNoise.FastOreNoise(float64(x), float64(y))

	// Simplified ore distribution based on depth only
	if depthFromSurface > 25 && oreNoiseValue > 0.8 {
		return block.GoldOre
	} else if depthFromSurface > 20 && oreNoiseValue > 0.75 {
		return block.SilverOre
	} else if depthFromSurface > 15 && oreNoiseValue > 0.7 {
		return block.IronOre
	} else if depthFromSurface > 8 && oreNoiseValue > 0.65 {
		return block.CopperOre
	}

	return block.Stone // Default to stone
}

// isInCave determines if a position should be a cave using fast approach
func isInCave(x, y int) bool {
	initNoiseGenerators()

	// Don't generate caves too close to surface
	if y < 6 {
		return false
	}

	// Use fast cave noise for better performance
	caveNoiseValue := caveNoise.FastCaveNoise(float64(x), float64(y))

	// Simple cave threshold - no depth variation for performance
	return caveNoiseValue < -0.5
}

// isUnderworld checks if we're in the underworld layer
func isUnderworld(y, worldHeight int) bool {
	return y > worldHeight-6 // Bottom 6 layers are underworld
}

// ResetWorldGeneration forces regeneration with new random seeds
func ResetWorldGeneration() {
	// Reset all noise generators to nil to force regeneration with new seeds
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
