package noise

// Biome and underground layer generation

func (sn *SimplexNoise) TerrariaBiomeNoise(x float64) float64 {
	largeBiomes := sn.FractalNoise1D(x, 2, 0.002, 1.0, 0.5)      // Continental biome zones
	mediumBiomes := sn.FractalNoise1D(x*1.3, 3, 0.008, 0.6, 0.6) // Regional variations
	transitions := sn.Noise1D(x*0.01) * 0.3                      // Smooth transitions

	return largeBiomes + mediumBiomes + transitions
}

func (sn *SimplexNoise) TerrariaUndergroundNoise(x, y float64) float64 {
	dirtStoneTransition := sn.FractalNoise2D(x, y, 3, 0.03, 1.0, 0.6)    // Dirt to stone layer
	stoneVariation := sn.FractalNoise2D(x*1.2, y*0.8, 2, 0.02, 0.8, 0.5) // Stone layer patterns

	return dirtStoneTransition + stoneVariation*0.5
}

func (sn *SimplexNoise) TerrariaUnderworldNoise(x, y float64) float64 {
	ashPockets := sn.FractalNoise2D(x, y, 4, 0.08, 1.0, 0.7)          // Ash formations
	lavaPockets := sn.FractalNoise2D(x*0.7, y*1.3, 3, 0.05, 1.2, 0.6) // Lava chambers
	hellstoneVeins := sn.RidgedNoise2D(x*2, y*0.5, 2, 0.1, 1.0)       // Hellstone ore veins

	return ashPockets + lavaPockets*0.8 + hellstoneVeins*0.6
}
