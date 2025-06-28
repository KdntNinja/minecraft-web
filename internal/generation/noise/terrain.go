package noise

// SimpleTerrainNoise provides basic terrain height using 1D noise
func (sn *PerlinNoise) SimpleTerrainNoise(x float64) float64 {
	return sn.Noise1D(x * 0.01)
}
