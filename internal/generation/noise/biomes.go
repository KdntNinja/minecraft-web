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
	DesertBiome
	MountainBiome
	SwampBiome
	TundraBiome
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

	// Determine biome based on temperature and humidity
	var biomeType BiomeType
	if elevation > 0.6 {
		biomeType = MountainBiome
	} else if temperature < -0.4 {
		biomeType = TundraBiome
	} else if humidity < -0.4 {
		biomeType = DesertBiome
	} else if temperature > 0.5 && humidity > 0.3 {
		biomeType = JungleBiome
	} else if humidity > 0.2 && temperature > -0.2 {
		biomeType = ForestBiome
	} else if humidity > 0.5 && elevation < -0.2 {
		biomeType = SwampBiome
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

// GetBiomeTerrainHeight returns terrain height modified by biome characteristics (optimized for WASM)
func (sn *PerlinNoise) GetBiomeTerrainHeight(x float64, biome BiomeData) float64 {
	baseHeight := sn.TerrainNoise(x)

	switch biome.Type {
	case MountainBiome:
		mountainHeight := sn.TerrainNoise(x)
		return baseHeight*0.3 + mountainHeight*0.7 + biome.Elevation*0.4
	case PlainseBiome:
		// Jungle has varied terrain with steep hills
		return baseHeight + biome.Elevation*0.2
	case OceanBiome:
		// Ocean floor variation
		return baseHeight*0.5 + biome.Elevation*0.2
	default:
		return baseHeight + biome.Elevation*0.2
	}
}
