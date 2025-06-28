package noise

// Utility functions are in helpers.go

import (
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// Biome and underground layer generation

func (sn *PerlinNoise) TerrariaBiomeNoise(x float64) float64 {
	largeBiomes := sn.FractalNoise1D(x, 2, 0.002, 1.0, settings.PerlinPersistence)                           // Continental biome zones
	mediumBiomes := sn.FractalNoise1D(x*1.3, settings.PerlinOctaves, 0.008, 0.6, settings.PerlinPersistence) // Regional variations
	transitions := sn.Noise1D(x*0.01) * 0.3                                                                  // Smooth transitions

	return largeBiomes + mediumBiomes + transitions
}

func (sn *PerlinNoise) TerrariaUndergroundNoise(x, y float64) float64 {
	dirtStoneTransition := sn.FractalNoise2D(x, y, settings.PerlinOctaves, 0.03, 1.0, 0.6)      // Dirt to stone layer
	stoneVariation := sn.FractalNoise2D(x*1.2, y*0.8, 2, 0.02, 0.8, settings.PerlinPersistence) // Stone layer patterns

	return dirtStoneTransition + stoneVariation*0.5
}

func (sn *PerlinNoise) TerrariaUnderworldNoise(x, y float64) float64 {
	ashPockets := sn.FractalNoise2D(x, y, 4, 0.08, 1.0, 0.7)          // Ash formations
	lavaPockets := sn.FractalNoise2D(x*0.7, y*1.3, 3, 0.05, 1.2, 0.6) // Lava chambers
	hellstoneVeins := sn.RidgedNoise2D(x*2, y*0.5, 2, 0.1, 1.0)       // Hellstone ore veins

	return ashPockets + lavaPockets*0.8 + hellstoneVeins*0.6
}

// BiomeType represents different biome types
type BiomeType int

const (
	PlainseBiome BiomeType = iota
	ForestBiome
	MountainBiome
	JungleBiome
	OceanBiome
)

// BiomeData contains biome characteristics
type BiomeData struct {
	Type        BiomeType
	Temperature float64 // -1.0 (cold) to 1.0 (hot)
	Humidity    float64 // -1.0 (dry) to 1.0 (wet)
	Elevation   float64 // Terrain height modifier
}

// GetBiomeAt determines the biome at a given x coordinate (optimized for WASM)
func (sn *PerlinNoise) GetBiomeAt(x float64) BiomeData {
	// Use fewer octaves for performance
	temperature := sn.FractalNoise1D(x*0.001, 2, 0.002, 1.0, 0.5)
	humidity := sn.FractalNoise1D(x*0.0015+1000, 2, 0.0025, 1.0, 0.5)
	elevation := sn.FractalNoise1D(x*0.0008+2000, 1, 0.001, 1.0, 0.6)

	// Clamp values
	if temperature < -1.0 {
		temperature = -1.0
	} else if temperature > 1.0 {
		temperature = 1.0
	}
	if humidity < -1.0 {
		humidity = -1.0
	} else if humidity > 1.0 {
		humidity = 1.0
	}
	if elevation < -1.0 {
		elevation = -1.0
	} else if elevation > 1.0 {
		elevation = 1.0
	}

	// Determine biome based on temperature and humidity - simplified to 5 biomes
	var biomeType BiomeType
	if elevation > 0.6 {
		biomeType = MountainBiome
	} else if temperature > 0.5 && humidity > 0.3 {
		biomeType = JungleBiome
	} else if humidity > 0.2 && temperature > -0.2 {
		biomeType = ForestBiome
	} else if elevation < -0.5 {
		biomeType = OceanBiome
	} else {
		biomeType = PlainseBiome
	}

	return BiomeData{
		Type:        biomeType,
		Temperature: temperature,
		Humidity:    humidity,
		Elevation:   elevation,
	}
}

// GetBiomeTerrainHeight returns Terraria-style terrain height modified by biome characteristics
func (sn *PerlinNoise) GetBiomeTerrainHeight(x float64, biome BiomeData) float64 {
	baseNoise := sn.FractalNoise1D(x*0.01, 2, 0.015, 1.0, 0.6)

	switch biome.Type {
	case MountainBiome:
		// Extremely tall and dramatic mountains with sharp peaks
		mountains := sn.FractalNoise1D(x*0.008, 3, 0.012, 2.0, 0.5)
		peaks := sn.RidgedNoise1D(x*0.02, 2, 0.025, 1.5)
		plateaus := sn.FractalNoise1D(x*0.005, 2, 0.008, 1.2, 0.7)
		return baseNoise*0.4 + mountains*1.8 + peaks*1.2 + plateaus*0.8

	case PlainseBiome:
		// Gentle rolling hills with occasional dramatic rises
		hills := sn.FractalNoise1D(x*0.015, 2, 0.018, 1.0, 0.6)
		rolling := sn.FractalNoise1D(x*0.04, 1, 0.045, 0.6, 0.5)
		return baseNoise*0.5 + hills*0.8 + rolling*0.4

	case ForestBiome:
		// Varied mountainous forest terrain with deep valleys
		forestHills := sn.FractalNoise1D(x*0.012, 3, 0.015, 1.2, 0.6)
		valleys := sn.FractalNoise1D(x*0.008, 2, 0.01, 1.0, 0.7)
		ridges := sn.FractalNoise1D(x*0.03, 2, 0.035, 0.8, 0.5)
		return baseNoise*0.4 + forestHills*1.3 + valleys*1.0 + ridges*0.6

	case JungleBiome:
		// Dense, varied terrain with deep ravines and towering canopy areas
		jungleHills := sn.FractalNoise1D(x*0.014, 3, 0.017, 1.5, 0.6)
		ravines := sn.FractalNoise1D(x*0.009, 2, 0.011, 1.2, 0.7)
		canopyRidges := sn.RidgedNoise1D(x*0.02, 2, 0.025, 0.8)
		return baseNoise*0.4 + jungleHills*1.6 + ravines*1.1 + canopyRidges*0.7

	case OceanBiome:
		// Underwater terrain with deep trenches and seamounts
		seafloor := sn.FractalNoise1D(x*0.02, 2, 0.025, 0.8, 0.6)
		trenches := sn.FractalNoise1D(x*0.008, 2, 0.01, 1.0, 0.7)
		return baseNoise*0.2 + seafloor*0.4 + trenches*0.6 - 1.5 // Deep underwater

	default:
		return baseNoise * 1.0
	}
}
