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
	veinShape := sn.Noise2D(x*scale*2, y*scale*2) * 0.3                           // Vein shape variation

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

// Advanced ore generation with realistic distribution patterns
func (sn *SimplexNoise) RealisticOreGeneration(x, y float64, depthFromSurface int, biome BiomeData) int {
	// Base ore density increases with depth
	depthMultiplier := float64(depthFromSurface) / 50.0 // More ores deeper down

	// Biome-specific ore preferences
	biomeModifier := 1.0
	switch biome.Type {
	case MountainBiome:
		biomeModifier = 1.3 // Mountains have more ores
	case DesertBiome:
		biomeModifier = 0.8 // Deserts have fewer ores
	case SwampBiome:
		biomeModifier = 0.9 // Swamps have different ore distribution
	}

	// Layered ore generation based on depth
	if depthFromSurface < 5 {
		// Surface layer - mostly common materials
		coalNoise := sn.FractalNoise2D(x*0.1, y*0.08, 3, 0.15, 1.0, 0.6) * biomeModifier
		if coalNoise > 0.6 {
			return 1 // Coal
		}
	} else if depthFromSurface < 20 {
		// Shallow underground - iron and copper
		ironNoise := sn.FractalNoise2D(x*0.08, y*0.1, 3, 0.12, 1.0, 0.5) * biomeModifier * depthMultiplier
		copperNoise := sn.FractalNoise2D(x*0.12+100, y*0.09, 3, 0.14, 1.0, 0.6) * biomeModifier

		if ironNoise > 0.7 {
			return 2 // Iron
		}
		if copperNoise > 0.65 {
			return 3 // Copper
		}
	} else if depthFromSurface < 50 {
		// Medium depth - precious metals
		goldNoise := sn.FractalNoise2D(x*0.06, y*0.07, 4, 0.08, 1.2, 0.4) * biomeModifier * depthMultiplier
		silverNoise := sn.FractalNoise2D(x*0.07+200, y*0.08, 3, 0.1, 1.0, 0.5) * biomeModifier * depthMultiplier

		if goldNoise > 0.85 {
			return 4 // Gold
		}
		if silverNoise > 0.8 {
			return 5 // Silver
		}
	} else {
		// Deep underground - rare ores and gems
		diamondNoise := sn.FractalNoise2D(x*0.04, y*0.05, 4, 0.06, 1.5, 0.3) * biomeModifier * depthMultiplier
		platinumNoise := sn.FractalNoise2D(x*0.05+300, y*0.06, 3, 0.07, 1.2, 0.4) * biomeModifier * depthMultiplier

		if diamondNoise > 0.9 {
			return 6 // Diamond
		}
		if platinumNoise > 0.88 {
			return 7 // Platinum
		}
	}

	return 0 // No ore
}

// CaveGeneration creates realistic cave systems
func (sn *SimplexNoise) GenerateCaves(x, y float64, depthFromSurface int) bool {
	// Caves are more common at medium depths
	depthFactor := 1.0
	if depthFromSurface > 10 && depthFromSurface < 80 {
		depthFactor = 1.5 // More caves at medium depth
	} else if depthFromSurface > 80 {
		depthFactor = 2.0 // Large cave systems deep underground
	}

	// Multi-octave cave noise for complex cave systems
	largeCaves := sn.FractalNoise2D(x*0.02, y*0.02, 3, 0.03, 1.0, 0.6) * depthFactor
	mediumCaves := sn.FractalNoise2D(x*0.05, y*0.04, 2, 0.06, 0.8, 0.5) * depthFactor
	smallCaves := sn.Noise2D(x*0.1, y*0.08) * 0.3 * depthFactor

	// Cave worm patterns for interesting cave shapes
	wormX := sn.FractalNoise2D(x*0.01, y*0.015, 2, 0.02, 1.0, 0.5)
	wormY := sn.FractalNoise2D(x*0.015, y*0.01, 2, 0.02, 1.0, 0.5)
	wormCaves := sn.Noise2D(x+wormX*20, y+wormY*20) * 0.4 * depthFactor

	totalCaveNoise := largeCaves + mediumCaves + smallCaves + wormCaves

	// Threshold for cave generation - varies by depth
	threshold := 0.6
	if depthFromSurface > 50 {
		threshold = 0.5 // Easier to generate caves deep down
	}

	return totalCaveNoise > threshold
}

// UndergroundStructures generates special underground rooms and features
func (sn *SimplexNoise) GenerateUndergroundStructures(x, y float64, depthFromSurface int) int {
	// Only generate structures at certain depths
	if depthFromSurface < 20 || depthFromSurface > 100 {
		return 0 // No structure
	}

	// Rare structure generation
	structureNoise := sn.FractalNoise2D(x*0.005, y*0.005, 2, 0.01, 1.0, 0.5)

	if structureNoise > 0.95 {
		// Determine structure type based on depth and biome
		typeNoise := sn.Noise2D(x*0.1, y*0.1)

		if depthFromSurface > 70 {
			if typeNoise > 0.5 {
				return 3 // Deep crystal cavern
			} else {
				return 2 // Underground lake
			}
		} else {
			return 1 // Treasure room
		}
	}

	return 0 // No structure
}
