package noise

// Cave generation functions for underground tunnel systems

func (sn *SimplexNoise) FastCaveNoise(x, y float64) float64 {
	return sn.Noise2D(x*0.05, y*0.08) // Single octave for performance
}

func (sn *SimplexNoise) HybridCaveNoise(x, y float64) float64 {
	minecraft := sn.Noise2D(x*0.04, y*0.02) // Horizontal tunnel bias
	terraria := sn.TerrariaCaveNoise(x, y)  // More chaotic caves

	return (minecraft + terraria) * 0.5
}

func (sn *SimplexNoise) TerrariaCaveNoise(x, y float64) float64 {
	primaryCaves := sn.Noise2D(x*0.03, y*0.015)        // Main tunnel system
	secondaryCaves := sn.Noise2D(x*0.06, y*0.04) * 0.7 // Smaller side tunnels
	largeCaverns := sn.Noise2D(x*0.01, y*0.008) * 1.2  // Rare large chambers

	return primaryCaves + secondaryCaves + largeCaverns*0.6
}
