package noise

// Enhanced cave generation functions for underground tunnel systems with seed-based variety

func (sn *SimplexNoise) FastCaveNoise(x, y float64) float64 {
	return sn.Noise2D(x*0.05, y*0.08) // Single octave for performance
}

func (sn *SimplexNoise) HybridCaveNoise(x, y float64) float64 {
	minecraft := sn.Noise2D(x*0.04, y*0.02) // Horizontal tunnel bias
	terraria := sn.TerrariaCaveNoise(x, y)  // More chaotic caves

	// Seed-based blending for variety
	seedBlend := float64((sn.seed/100)%100) / 100.0
	return minecraft*(1.0-seedBlend) + terraria*seedBlend
}

func (sn *SimplexNoise) TerrariaCaveNoise(x, y float64) float64 {
	// Seed-based cave system variation
	seedOffset := float64(sn.seed%1000) * 0.001
	
	primaryCaves := sn.Noise2D(x*0.03+seedOffset, y*0.015)        // Main tunnel system
	secondaryCaves := sn.Noise2D(x*0.06-seedOffset, y*0.04) * 0.7 // Smaller side tunnels
	largeCaverns := sn.Noise2D(x*0.01, y*0.008+seedOffset) * 1.2  // Rare large chambers

	// Add seed-based chaos for unique cave patterns
	seedChaos := sn.RidgedNoise2D(x+seedOffset*1000, y-seedOffset*1000, 2, 0.02, 0.3)

	return primaryCaves + secondaryCaves + largeCaverns*0.6 + seedChaos
}

// EnhancedCaveNoise provides sophisticated cave generation with multiple layers
func (sn *SimplexNoise) EnhancedCaveNoise(x, y float64) float64 {
	// Depth-based cave density
	depth := y * 0.01
	depthModifier := depth * 0.3 // More caves deeper down
	
	// Primary cave network with seed variation
	seedVar := float64(sn.seed%777) * 0.001
	mainCaves := sn.FractalNoise2D(x+seedVar, y-seedVar, 3, 0.025, 1.0, 0.6)
	
	// Secondary tunnel network
	tunnels := sn.Noise2D(x*0.045, y*0.035) * 0.8
	
	// Large chamber system
	chambers := sn.RidgedNoise2D(x, y, 2, 0.015, 1.5) * 0.4
	
	// Micro-cave details
	microCaves := sn.Noise2D(x*0.1, y*0.12) * 0.2
	
	return mainCaves + tunnels + chambers + microCaves + depthModifier
}
