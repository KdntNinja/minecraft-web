package generation

import (
	"github.com/KdntNinja/webcraft/settings"
	"github.com/aquilax/go-perlin"
)

var (
	generationSeed int64
	terrainNoise   *perlin.Perlin
	cavesNoise     *perlin.Perlin
	oresNoise      *perlin.Perlin
)

// InitializeNoise initializes all noise generators with the given seed
func InitializeNoise(seed int64) {
	generationSeed = seed
	terrainNoise = perlin.NewPerlin(settings.PerlinAlpha, settings.PerlinBeta, settings.PerlinOctaves, seed)
	cavesNoise = perlin.NewPerlin(settings.PerlinAlpha, settings.PerlinBeta, settings.PerlinOctaves, seed+1000)
	oresNoise = perlin.NewPerlin(settings.PerlinAlpha, settings.PerlinBeta, settings.PerlinOctaves, seed+2000)
}

// ResetGeneration forces regeneration with a new provided seed
func ResetGeneration(seed int64) {
	ResetHeightCache()
	generationSeed = 0
	terrainNoise = nil
	cavesNoise = nil
	oresNoise = nil
	InitializeNoise(seed)
}

// ResetWorldGeneration forces regeneration with a new provided seed (legacy name for compatibility)
func ResetWorldGeneration(seed int64) {
	ResetGeneration(seed)
}

// GetSeed returns the current generation seed for display or saving
func GetSeed() int64 {
	if terrainNoise == nil {
		InitializeNoise(42) // Default seed
	}
	return generationSeed
}

// GetWorldSeed returns the current world seed for display or saving (legacy name for compatibility)
func GetWorldSeed() int64 {
	return GetSeed()
}

// GetTerrainNoise returns the terrain noise generator
func GetTerrainNoise() *perlin.Perlin {
	if terrainNoise == nil {
		InitializeNoise(42) // Default seed
	}
	return terrainNoise
}

// GetCaveNoise returns the cave noise generator
func GetCaveNoise() *perlin.Perlin {
	if cavesNoise == nil {
		InitializeNoise(42) // Default seed
	}
	return cavesNoise
}

// GetOreNoise returns the ore noise generator
func GetOreNoise() *perlin.Perlin {
	if oresNoise == nil {
		InitializeNoise(42) // Default seed
	}
	return oresNoise
}
