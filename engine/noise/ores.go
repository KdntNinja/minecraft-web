package noise

// Ore vein generation for underground mining

func (sn *SimplexNoise) FastOreNoise(x, y float64) float64 {
	return sn.Noise2D(x*0.1, y*0.1) // Simple ore distribution
}

func (sn *SimplexNoise) TerrariaOreNoise(x, y float64, oreType int) float64 {
	var scale, threshold float64

	switch oreType {
	case 0: // Copper - common, small veins
		scale = 0.15
		threshold = 0.7
	case 1: // Iron - medium rarity
		scale = 0.12
		threshold = 0.75
	case 2: // Silver - less common
		scale = 0.1
		threshold = 0.8
	case 3: // Gold - rare, rich veins
		scale = 0.08
		threshold = 0.85
	case 4: // Platinum - very rare
		scale = 0.06
		threshold = 0.9
	default:
		scale = 0.1
		threshold = 0.8
	}

	oreNoise := sn.FractalNoise2D(x, y, 3, scale, 1.0, 0.5) // Base vein pattern
	veinShape := sn.Noise2D(x*scale*2, y*scale*2) * 0.3     // Vein shape variation

	return oreNoise + veinShape - threshold
}
