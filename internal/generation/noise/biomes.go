package noise

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

// GetBiomeAt determines the biome at a given x coordinate
func (sn *PerlinNoise) GetBiomeAt(x float64) BiomeData {
	// Use multiple noise layers for biome determination
	temperature := sn.FractalNoise1D(x*0.001, 3, 0.002, 1.0, 0.5)
	humidity := sn.FractalNoise1D(x*0.0015+1000, 3, 0.0025, 1.0, 0.5)
	elevation := sn.FractalNoise1D(x*0.0008+2000, 2, 0.001, 1.0, 0.6)

	// Normalize values to -1 to 1 range
	temperature = clamp(temperature, -1.0, 1.0)
	humidity = clamp(humidity, -1.0, 1.0)
	elevation = clamp(elevation, -1.0, 1.0)

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

// GetBiomeTerrainHeight returns terrain height modified by biome characteristics
func (sn *PerlinNoise) GetBiomeTerrainHeight(x float64, biome BiomeData) float64 {
	baseHeight := sn.SimpleTerrainNoise(x)

	switch biome.Type {
	case MountainBiome:
		mountainHeight := sn.SimpleTerrainNoise(x)
		return baseHeight*0.3 + mountainHeight*0.7 + biome.Elevation*0.4

	case PlainseBiome:
		plainHeight := sn.SimpleTerrainNoise(x)
		return baseHeight*0.2 + plainHeight*0.8 + biome.Elevation*0.2

	case DesertBiome:
		// Desert has dunes and occasional mesas
		dunes := sn.SimpleTerrainNoise(x)
		mesas := sn.SimpleTerrainNoise(x) * 0.3
		return baseHeight*0.4 + dunes + mesas + biome.Elevation*0.3

	case ForestBiome:
		// Forest has gentle rolling hills
		return baseHeight + biome.Elevation*0.2

	case SwampBiome:
		// Swamps are mostly flat with occasional small hills
		return baseHeight + biome.Elevation*0.2

	case TundraBiome:
		// Tundra has permafrost ridges and gentle slopes
		return baseHeight + biome.Elevation*0.2

	case JungleBiome:
		// Jungle has varied terrain with steep hills
		return baseHeight + biome.Elevation*0.2

	case OceanBiome:
		// Ocean floor variation
		return baseHeight*0.5 + biome.Elevation*0.2

	default:
		return baseHeight + biome.Elevation*0.2
	}
}

// Helper function to clamp values
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
