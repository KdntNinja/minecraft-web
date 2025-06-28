package generation

// IsCave determines if a position should be a cave using 3D-like noise
func IsCave(worldX, worldY int) bool {
	// Only generate caves underground (below surface + some buffer)
	surfaceHeight := GetHeightAt(worldX)
	if worldY <= surfaceHeight+8 {
		return false
	}

	caveNoise := GetCaveNoise()

	x := float64(worldX)
	y := float64(worldY)
	depth := worldY - surfaceHeight

	// Large cave systems using low-frequency noise
	largeCaveNoise := caveNoise.Noise2D(x/40.0, y/40.0)

	// Smaller tunnels using medium-frequency noise
	smallCaveNoise := caveNoise.Noise2D(x/20.0+500, y/20.0+500)

	// Tiny air pockets using high-frequency noise
	pocketNoise := caveNoise.Noise2D(x/8.0+1000, y/8.0+1000)

	// Combine different cave sizes with depth-based probability
	var caveThreshold float64

	if depth > 100 {
		// Deep caves - larger and more common
		caveThreshold = 0.45
		combinedCave := largeCaveNoise*0.6 + smallCaveNoise*0.3 + pocketNoise*0.1
		return combinedCave > caveThreshold
	} else if depth > 50 {
		// Medium depth caves
		caveThreshold = 0.55
		combinedCave := largeCaveNoise*0.4 + smallCaveNoise*0.5 + pocketNoise*0.1
		return combinedCave > caveThreshold
	} else {
		// Shallow caves - smaller and rarer
		caveThreshold = 0.65
		combinedCave := smallCaveNoise*0.7 + pocketNoise*0.3
		return combinedCave > caveThreshold
	}
}
