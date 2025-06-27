package noise

// Biome and underground generation functions

// TerrariaBiomeNoise generates biome boundaries similar to Terraria
func (sn *SimplexNoise) TerrariaBiomeNoise(x float64) float64 {
	// Large-scale biome distribution
	largeBiomes := sn.FractalNoise1D(x, 2, 0.002, 1.0, 0.5)

	// Medium-scale biome variations
	mediumBiomes := sn.FractalNoise1D(x*1.3, 3, 0.008, 0.6, 0.6)

	// Small transition zones
	transitions := sn.Noise1D(x*0.01) * 0.3

	return largeBiomes + mediumBiomes + transitions
}

// TerrariaUndergroundNoise generates underground layer transitions
func (sn *SimplexNoise) TerrariaUndergroundNoise(x, y float64) float64 {
	// Dirt to stone transition - should be somewhat uneven
	dirtStoneTransition := sn.FractalNoise2D(x, y, 3, 0.03, 1.0, 0.6)

	// Stone layer variations
	stoneVariation := sn.FractalNoise2D(x*1.2, y*0.8, 2, 0.02, 0.8, 0.5)

	return dirtStoneTransition + stoneVariation*0.5
}

// TerrariaUnderworldNoise generates hell/underworld terrain like Terraria
func (sn *SimplexNoise) TerrariaUnderworldNoise(x, y float64) float64 {
	// Jagged, chaotic terrain for underworld
	ashPockets := sn.FractalNoise2D(x, y, 4, 0.08, 1.0, 0.7)
	lavaPockets := sn.FractalNoise2D(x*0.7, y*1.3, 3, 0.05, 1.2, 0.6)
	hellstoneVeins := sn.RidgedNoise2D(x*2, y*0.5, 2, 0.1, 1.0)

	return ashPockets + lavaPockets*0.8 + hellstoneVeins*0.6
}
