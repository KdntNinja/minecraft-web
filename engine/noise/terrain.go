package noise

// Terrain generation functions for different world styles

func (sn *SimplexNoise) FastTerrainNoise(x float64) float64 {
	base := sn.Noise1D(x * 0.01)       // Large landforms
	detail := sn.Noise1D(x*0.03) * 0.5 // Surface detail
	return base + detail
}

func (sn *SimplexNoise) MinecraftTerrainNoise(x float64) float64 {
	continental := sn.Noise1D(x*0.005) * 0.8 // Large continental shape
	hills := sn.Noise1D(x*0.02) * 0.3        // Rolling hills
	details := sn.Noise1D(x*0.08) * 0.1      // Fine surface variation

	return continental + hills + details
}

func (sn *SimplexNoise) TerrariaTerrainNoise(x float64) float64 {
	largeTerrain := sn.Noise1D(x * 0.008)     // Base landmasses
	mediumTerrain := sn.Noise1D(x*0.02) * 0.5 // Medium hills
	smallTerrain := sn.Noise1D(x*0.05) * 0.25 // Small bumps

	return largeTerrain + mediumTerrain + smallTerrain
}

func (sn *SimplexNoise) HybridTerrainNoise(x float64) float64 {
	minecraft := sn.MinecraftTerrainNoise(x)
	terraria := sn.TerrariaTerrainNoise(x)

	blendFactor := (sn.Noise1D(x*0.003) + 1.0) * 0.5 // 0 to 1 blend ratio

	return minecraft*(1.0-blendFactor) + terraria*blendFactor
}
