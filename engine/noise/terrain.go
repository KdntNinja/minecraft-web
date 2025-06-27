package noise

// Terrain generation functions for different game styles

// FastTerrainNoise generates simple but fast terrain
func (sn *SimplexNoise) FastTerrainNoise(x float64) float64 {
	// Only use 2 octaves for speed
	base := sn.Noise1D(x * 0.01)
	detail := sn.Noise1D(x*0.03) * 0.5
	return base + detail
}

// MinecraftTerrainNoise generates Minecraft-like terrain with smoother hills
func (sn *SimplexNoise) MinecraftTerrainNoise(x float64) float64 {
	// Minecraft-style terrain with gentler slopes
	continental := sn.Noise1D(x*0.005) * 0.8 // Continental shape
	hills := sn.Noise1D(x*0.02) * 0.3        // Rolling hills
	details := sn.Noise1D(x*0.08) * 0.1      // Fine details

	return continental + hills + details
}

// TerrariaTerrainNoise generates Terraria-like surface terrain with hills and valleys
func (sn *SimplexNoise) TerrariaTerrainNoise(x float64) float64 {
	// Base terrain using multiple octaves like Terraria
	largeTerrain := sn.Noise1D(x * 0.008)     // Large landmasses
	mediumTerrain := sn.Noise1D(x*0.02) * 0.5 // Hills
	smallTerrain := sn.Noise1D(x*0.05) * 0.25 // Small details

	return largeTerrain + mediumTerrain + smallTerrain
}

// HybridTerrainNoise combines Minecraft smoothness with Terraria variety
func (sn *SimplexNoise) HybridTerrainNoise(x float64) float64 {
	// Get both terrain types
	minecraft := sn.MinecraftTerrainNoise(x)
	terraria := sn.TerrariaTerrainNoise(x)

	// Blend them based on position for variety
	blendFactor := (sn.Noise1D(x*0.003) + 1.0) * 0.5 // 0 to 1

	// More Minecraft-like in some areas, more Terraria-like in others
	return minecraft*(1.0-blendFactor) + terraria*blendFactor
}
