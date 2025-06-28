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

	// Multiple cave types for variety

	// 1. Large caverns (Terraria-style)
	largeCaveNoise := caveNoise.Noise2D(x/60.0, y/60.0)

	// 2. Winding tunnels (horizontal emphasis)
	horizontalTunnels := caveNoise.Noise2D(x/25.0+1000, y/45.0+1000)

	// 3. Vertical shafts
	verticalShafts := caveNoise.Noise2D(x/45.0+2000, y/20.0+2000)

	// 4. Small pockets and connecting passages
	smallCaves := caveNoise.Noise2D(x/15.0+3000, y/15.0+3000)

	// 5. Tiny air bubbles
	airPockets := caveNoise.Noise2D(x/8.0+4000, y/8.0+4000)

	// Depth-based cave generation with different types
	if depth > 150 {
		// Very deep - large caverns and complex systems
		largeCaverns := largeCaveNoise*0.5 + horizontalTunnels*0.3 + verticalShafts*0.2
		if largeCaverns > 0.35 {
			return true
		}

		// Additional tunnel networks deep underground
		deepTunnels := smallCaves*0.6 + airPockets*0.4
		return deepTunnels > 0.55

	} else if depth > 100 {
		// Deep caves - mix of large and medium caves
		deepCaves := largeCaveNoise*0.4 + horizontalTunnels*0.4 + smallCaves*0.2
		if deepCaves > 0.4 {
			return true
		}

		// Vertical connections between levels
		return verticalShafts > 0.65

	} else if depth > 50 {
		// Medium depth - primarily horizontal tunnel systems
		mediumCaves := horizontalTunnels*0.5 + smallCaves*0.3 + airPockets*0.2
		if mediumCaves > 0.5 {
			return true
		}

		// Some vertical shafts connecting to surface
		return verticalShafts > 0.7

	} else {
		// Shallow caves - small pockets and rare tunnels
		shallowCaves := smallCaves*0.6 + airPockets*0.4
		return shallowCaves > 0.65
	}
}

// IsLargeCavern determines if a position should be part of a large underground cavern
func IsLargeCavern(worldX, worldY int) bool {
	surfaceHeight := GetHeightAt(worldX)
	depth := worldY - surfaceHeight

	// Only large caverns at significant depth
	if depth < 80 {
		return false
	}

	caveNoise := GetCaveNoise()
	x := float64(worldX)
	y := float64(worldY)

	// Large cavern noise with emphasis on creating big open spaces
	cavernNoise := caveNoise.Noise2D(x/80.0, y/80.0)

	// Add some variation to cavern shape
	shapeVariation := caveNoise.Noise2D(x/40.0+5000, y/40.0+5000)

	combinedNoise := cavernNoise*0.7 + shapeVariation*0.3

	// Threshold varies by depth - deeper = more likely to have large caverns
	threshold := 0.6 - float64(depth-80)*0.001 // Gets easier deeper down
	if threshold < 0.3 {
		threshold = 0.3
	}

	return combinedNoise > threshold
}

// GetCaveWaterLevel determines if a cave position should have water
func GetCaveWaterLevel(worldX, worldY int) bool {
	surfaceHeight := GetHeightAt(worldX)
	depth := worldY - surfaceHeight

	// Water only appears in deeper caves
	if depth < 60 {
		return false
	}

	// Use ore noise for water placement to avoid conflicts
	oreNoise := GetOreNoise()
	x := float64(worldX)
	y := float64(worldY)

	waterNoise := oreNoise.Noise2D(x/30.0+6000, y/30.0+6000)

	// Water pools in low-lying cave areas
	// More common in very deep areas
	threshold := 0.75 - float64(depth-60)*0.002
	if threshold < 0.65 {
		threshold = 0.65
	}

	return waterNoise > threshold
}
