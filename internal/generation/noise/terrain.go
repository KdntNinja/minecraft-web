package noise

// Utility functions are in helpers.go

// TerrainNoise provides basic but richer terrain height using lightweight fractal noise
func (sn *PerlinNoise) TerrainNoise(x float64) float64 {
	// Use 2 octaves for gentle hills and valleys, tuned for web performance
	base := sn.Noise1D(x * 0.01)
	detail := sn.Noise1D(x*0.05) * 0.4
	return base + detail
}
