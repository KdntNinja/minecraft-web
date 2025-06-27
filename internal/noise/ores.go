package noise

// Ore vein generation for underground mining with enhanced randomness

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

	// Seed-based ore pattern variation
	seedOffset := float64(sn.seed%1000) * 0.001
	
	oreNoise := sn.FractalNoise2D(x+seedOffset, y-seedOffset, 3, scale, 1.0, 0.5) // Base vein pattern
	veinShape := sn.Noise2D(x*scale*2, y*scale*2) * 0.3     // Vein shape variation
	
	// Add seed-based vein character
	veinCharacter := float64((sn.seed/int64(oreType+1))%100) / 100.0 * 0.1

	return oreNoise + veinShape + veinCharacter - threshold
}

// EnhancedOreNoise provides sophisticated ore distribution with depth-based rarity
func (sn *SimplexNoise) EnhancedOreNoise(x, y float64, depthFromSurface int) int {
	// Seed-based ore variation
	seedVar := float64(sn.seed%555) * 0.001
	
	// Base ore noise
	oreNoise := sn.FractalNoise2D(x+seedVar, y-seedVar, 2, 0.08, 1.0, 0.6)
	
	// Depth-based ore type determination
	switch {
	case depthFromSurface > 30:
		// Deep ores: gold, platinum
		if oreNoise > 0.85 {
			return 4 // Platinum
		} else if oreNoise > 0.75 {
			return 3 // Gold
		} else if oreNoise > 0.65 {
			return 2 // Silver
		}
	case depthFromSurface > 20:
		// Medium depth: silver, iron
		if oreNoise > 0.8 {
			return 2 // Silver
		} else if oreNoise > 0.7 {
			return 1 // Iron
		}
	case depthFromSurface > 10:
		// Shallow: copper, iron
		if oreNoise > 0.75 {
			return 1 // Iron
		} else if oreNoise > 0.65 {
			return 0 // Copper
		}
	default:
		// Very shallow: only copper
		if oreNoise > 0.7 {
			return 0 // Copper
		}
	}

	return -1 // No ore
}
