package terrain

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation/noise"
)

var (
	terrainHeights = make(map[int]int)
	terrainNoise   *noise.PerlinNoise
)

func initTerrainNoise() {
	if terrainNoise == nil {
		terrainNoise = noise.NewPerlinNoise(42) // Use a fixed seed for simplicity
		terrainHeights = make(map[int]int)
	}
}

func getSurfaceHeight(x int) int {
	h, ok := terrainHeights[x]
	if ok {
		return h
	}
	initTerrainNoise()

	baseHeight := 12 // Base terrain height
	heightNoise := terrainNoise.SimpleTerrainNoise(float64(x))
	height := baseHeight + int(heightNoise*8) // Simple scaling

	if height < 3 {
		height = 3
	}
	if height > settings.ChunkHeight-2 {
		height = settings.ChunkHeight - 2
	}

	terrainHeights[x] = height
	return height
}

// BiomeType represents different surface biomes
// Only a single biome for simplicity

type BiomeType int

const (
	DefaultBiome BiomeType = iota
)

func getBiome(x int) BiomeType {
	return DefaultBiome
}

// getOreType determines what ore should be placed using simple noise
func getOreType(x, y, depthFromSurface int) block.BlockType {
	initTerrainNoise()

	// Use simple noise for ore generation
	oreNoise := terrainNoise.Noise2D(float64(x), float64(y))

	// Map ore type numbers to block types
	switch {
	case oreNoise < -0.3:
		return block.CopperOre
	case oreNoise < 0:
		return block.IronOre
	case oreNoise < 0.3:
		return block.SilverOre
	case oreNoise < 0.6:
		return block.GoldOre
	case oreNoise < 0.9:
		return block.PlatinumOre
	default:
		return block.Stone // Default to stone
	}
}

// isInCave determines if a position should be a cave using simple noise
func isInCave(x, y int) bool {
	initTerrainNoise()

	// Don't generate caves too close to surface
	if y < 8 {
		return false
	}

	// Use simple noise for cave generation
	caveNoise := terrainNoise.Noise2D(float64(x), float64(y))

	// Simple cave threshold
	return caveNoise < -0.2
}

// isUnderworld checks if we're in the underworld layer
func isUnderworld(y, worldHeight int) bool {
	return y > worldHeight-6 // Bottom 6 layers are underworld
}

// ResetWorldGeneration forces regeneration with new random seeds
func ResetWorldGeneration() {
	// Reset terrain noise to nil to force regeneration with new seed
	terrainNoise = nil

	// Clear all caches
	terrainHeights = make(map[int]int)
	chunkCache = make(map[string]block.Chunk)
}

// GetWorldSeed returns the current world seed for display or saving
func GetWorldSeed() int64 {
	initTerrainNoise()
	return 42 // Fixed seed, as randomness is not used in this simplified version
}

// getBlockType returns the block type for a given world position, generating Minecraft-like layers
func getBlockType(x, y int) block.BlockType {
	initTerrainNoise()

	surfaceY := getSurfaceHeight(x)
	depth := y - surfaceY

	if y < 0 || y >= settings.ChunkHeight {
		return block.Air
	}

	if y < surfaceY {
		return block.Air
	}

	// Surface block
	if y == surfaceY {
		return block.Grass
	}

	// Dirt layer (3 blocks below surface)
	if y > surfaceY && y <= surfaceY+3 {
		return block.Dirt
	}

	// Stone and ores below dirt
	if y > surfaceY+3 && y < settings.ChunkHeight-6 {
		// Occasionally generate ores
		ore := getOreType(x, y, depth)
		if ore != block.Stone {
			return ore
		}
		return block.Stone
	}

	// Underworld (bottom 6 layers)
	if y >= settings.ChunkHeight-6 {
		return block.Hellstone
	}

	return block.Stone
}
