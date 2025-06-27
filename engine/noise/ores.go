package noise

// Ore generation functions for different game styles

// FastOreNoise generates ore patterns with minimal computation
func (sn *SimplexNoise) FastOreNoise(x, y float64) float64 {
	// Very simple ore generation
	return sn.Noise2D(x*0.1, y*0.1)
}

// TerrariaOreNoise generates ore vein patterns like Terraria
func (sn *SimplexNoise) TerrariaOreNoise(x, y float64, oreType int) float64 {
	// Different scales for different ore types
	var scale, threshold float64

	switch oreType {
	case 0: // Copper - common, small veins
		scale = 0.15
		threshold = 0.7
	case 1: // Iron - medium rarity, medium veins
		scale = 0.12
		threshold = 0.75
	case 2: // Silver - less common, smaller veins
		scale = 0.1
		threshold = 0.8
	case 3: // Gold - rare, small but rich veins
		scale = 0.08
		threshold = 0.85
	case 4: // Platinum - very rare, tiny veins
		scale = 0.06
		threshold = 0.9
	default:
		scale = 0.1
		threshold = 0.8
	}

	// Generate ore vein pattern
	oreNoise := sn.FractalNoise2D(x, y, 3, scale, 1.0, 0.5)

	// Add some randomness to vein shapes
	veinShape := sn.Noise2D(x*scale*2, y*scale*2) * 0.3

	return oreNoise + veinShape - threshold
}
