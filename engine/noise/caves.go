package noise

// Cave generation functions for different game styles

// FastCaveNoise generates simple cave patterns
func (sn *SimplexNoise) FastCaveNoise(x, y float64) float64 {
	// Single octave for maximum performance
	return sn.Noise2D(x*0.05, y*0.08)
}

// HybridCaveNoise creates caves that blend Minecraft and Terraria styles
func (sn *SimplexNoise) HybridCaveNoise(x, y float64) float64 {
	// Minecraft-style caves - more horizontal tunnels
	minecraft := sn.Noise2D(x*0.04, y*0.02)

	// Terraria-style caves - more varied
	terraria := sn.TerrariaCaveNoise(x, y)

	// Combine them
	return (minecraft + terraria) * 0.5
}

// TerrariaCaveNoise generates cave patterns similar to Terraria
func (sn *SimplexNoise) TerrariaCaveNoise(x, y float64) float64 {
	// Primary cave tunnels - horizontal bias
	primaryCaves := sn.Noise2D(x*0.03, y*0.015)

	// Secondary cave systems - more chaotic
	secondaryCaves := sn.Noise2D(x*0.06, y*0.04) * 0.7

	// Large caverns - rare but spacious
	largeCaverns := sn.Noise2D(x*0.01, y*0.008) * 1.2

	// Combine cave systems
	return primaryCaves + secondaryCaves + largeCaverns*0.6
}
