package generation

// GetOreType determines what type of ore (if any) should be at a position
func GetOreType(worldX, worldY int) int {
	// Only generate ores underground
	surfaceHeight := GetHeightAt(worldX)
	if worldY <= surfaceHeight+10 {
		return 0 // No ore near surface
	}

	oreNoise := GetOreNoise()

	x := float64(worldX)
	y := float64(worldY)
	depth := worldY - surfaceHeight

	// Multiple ore noise layers for different ore types
	copperNoise := oreNoise.Noise2D(x/20.0, y/20.0)
	ironNoise := oreNoise.Noise2D(x/25.0+1000, y/25.0+1000)
	goldNoise := oreNoise.Noise2D(x/35.0+2000, y/35.0+2000)

	// Copper ore - shallow, most common
	if depth > 15 && depth < 80 && copperNoise < -0.6 {
		return 1 // Copper ore
	}

	// Iron ore - medium depth, common
	if depth > 30 && depth < 120 && ironNoise < -0.65 {
		return 2 // Iron ore
	}

	// Gold ore - deep, rare
	if depth > 60 && goldNoise < -0.75 {
		return 3 // Gold ore
	}

	return 0 // No ore
}

// IsLiquid determines if a position should contain liquid (water or lava)
func IsLiquid(worldX, worldY int) int {
	surfaceHeight := GetHeightAt(worldX)
	depth := worldY - surfaceHeight

	// No liquids near surface
	if depth < 20 {
		return 0
	}

	oreNoise := GetOreNoise()

	x := float64(worldX)
	y := float64(worldY)

	// Water pools in medium depths
	if depth > 30 && depth < 100 {
		waterNoise := oreNoise.Noise2D(x/30.0+3000, y/30.0+3000)
		if waterNoise < -0.8 {
			return 1 // Water
		}
	}

	// Lava pools in deep areas
	if depth > 80 {
		lavaNoise := oreNoise.Noise2D(x/25.0+4000, y/25.0+4000)
		if lavaNoise < -0.82 {
			return 2 // Lava (we'll represent as Water for now since we don't have lava texture)
		}
	}

	return 0 // No liquid
}
