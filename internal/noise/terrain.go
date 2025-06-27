package noise

import (
	"math"
)

// Enhanced terrain generation functions with better variety and randomness

func (sn *SimplexNoise) FastTerrainNoise(x float64) float64 {
	base := sn.Noise1D(x * 0.01)       // Large landforms
	detail := sn.Noise1D(x*0.03) * 0.5 // Surface detail
	return base + detail
}

// EnhancedTerrainNoise provides better terrain with multiple octaves and seed-based variation
func (sn *SimplexNoise) EnhancedTerrainNoise(x float64) float64 {
	// Use seed-based variation for unique worlds
	seedVariation := float64(sn.seed%1000) / 1000.0 * 0.2
	
	// Multiple octaves for complex terrain
	continental := sn.FractalNoise1D(x, 2, 0.003, 1.0, 0.5)    // Large continental features
	hills := sn.FractalNoise1D(x*1.3, 3, 0.015, 0.6, 0.6)     // Rolling hills
	details := sn.FractalNoise1D(x*2.1, 2, 0.05, 0.3, 0.5)    // Fine surface details
	microDetails := sn.Noise1D(x*0.08) * 0.1                  // Micro-scale variation
	
	// Combine with seed-based offset for variety
	return continental + hills + details + microDetails + seedVariation
}

func (sn *SimplexNoise) MinecraftTerrainNoise(x float64) float64 {
	// Enhance with seed-based variation
	seedOffset := float64(sn.seed%500) / 500.0 * 0.15
	
	continental := sn.Noise1D(x*0.005) * 0.8 // Large continental shape
	hills := sn.Noise1D(x*0.02) * 0.3        // Rolling hills
	details := sn.Noise1D(x*0.08) * 0.1      // Fine surface variation

	return continental + hills + details + seedOffset
}

func (sn *SimplexNoise) TerrariaTerrainNoise(x float64) float64 {
	// Enhanced with seed-based chaos
	seedChaos := math.Sin(float64(sn.seed)*0.001) * 0.1
	
	largeTerrain := sn.Noise1D(x * 0.008)     // Base landmasses
	mediumTerrain := sn.Noise1D(x*0.02) * 0.5 // Medium hills
	smallTerrain := sn.Noise1D(x*0.05) * 0.25 // Small bumps
	
	// Add chaotic elements for Terraria-style terrain
	chaos := sn.RidgedNoise1D(x, 2, 0.03, 0.4) * 0.3

	return largeTerrain + mediumTerrain + smallTerrain + chaos + seedChaos
}

// AdvancedHybridTerrainNoise combines multiple terrain styles with better blending
func (sn *SimplexNoise) AdvancedHybridTerrainNoise(x float64) float64 {
	minecraft := sn.MinecraftTerrainNoise(x)
	terraria := sn.TerrariaTerrainNoise(x)
	enhanced := sn.EnhancedTerrainNoise(x)

	// Dynamic blending based on position and seed
	blendNoise1 := (sn.Noise1D(x*0.003) + 1.0) * 0.5    // 0 to 1 blend ratio
	blendNoise2 := (sn.Noise1D(x*0.007 + 100) + 1.0) * 0.5
	
	// Seed-based style preference
	seedBias := float64((sn.seed / 1000) % 3) / 3.0
	
	// Three-way blend with seed influence
	weight1 := (blendNoise1 + seedBias) / 2.0
	weight2 := (blendNoise2 + (1.0-seedBias)) / 2.0
	weight3 := 1.0 - weight1 - weight2
	
	// Normalize weights
	totalWeight := weight1 + weight2 + weight3
	weight1 /= totalWeight
	weight2 /= totalWeight
	weight3 /= totalWeight

	return minecraft*weight1 + terraria*weight2 + enhanced*weight3
}

func (sn *SimplexNoise) HybridTerrainNoise(x float64) float64 {
	return sn.AdvancedHybridTerrainNoise(x)
}

// MountainousTerrainNoise creates dramatic mountain terrain
func (sn *SimplexNoise) MountainousTerrainNoise(x float64) float64 {
	// Mountain ridges using ridged noise
	ridges := sn.RidgedNoise1D(x, 3, 0.01, 1.2)
	
	// Base mountain shape
	mountains := sn.FractalNoise1D(x, 4, 0.008, 1.0, 0.6)
	
	// Sharp peaks and valleys
	peaks := math.Abs(sn.Noise1D(x*0.02)) * 0.8
	
	// Seed-based mountain character
	seedCharacter := math.Sin(float64(sn.seed)*0.001) * 0.2
	
	return ridges*0.7 + mountains*0.5 + peaks + seedCharacter
}

// PlainsTerrainNoise creates flat terrain with gentle rolling hills
func (sn *SimplexNoise) PlainsTerrainNoise(x float64) float64 {
	// Very gentle rolling hills
	gentle := sn.FractalNoise1D(x, 2, 0.005, 0.4, 0.5)
	
	// Micro-variations
	micro := sn.Noise1D(x*0.08) * 0.05
	
	// Seed-based plain character
	seedFlat := (float64(sn.seed%100) / 100.0 - 0.5) * 0.1
	
	return gentle + micro + seedFlat
}
