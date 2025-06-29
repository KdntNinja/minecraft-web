package generation

import "github.com/KdntNinja/webcraft/internal/core/settings"

// IsCave determines if a position should be a cave using 3D-like noise
func IsCave(worldX, worldY int) bool {
	surfaceHeight := GetHeightAt(worldX)

	// Allow caves to generate closer to surface, including surface entrances
	if worldY < surfaceHeight-2 {
		return false // Above ground
	}

	caveNoise := GetCaveNoise()
	x := float64(worldX)
	y := float64(worldY)
	depth := worldY - surfaceHeight

	// Surface cave entrances (more common and visible)
	if depth >= settings.CaveSurfaceEntranceMinDepth && depth <= settings.CaveSurfaceEntranceMaxDepth {
		entranceNoise := caveNoise.Noise2D(x/settings.CaveSurfaceEntranceScale+settings.CaveSurfaceEntranceOffset, y/settings.CaveSurfaceEntranceScale+settings.CaveSurfaceEntranceOffset)
		if entranceNoise > settings.CaveSurfaceEntranceThresh {
			return true // Surface cave entrance
		}
	}

	// More interconnected caves: blend horizontal and vertical tunnels, increase their weights
	largeCaveNoise := caveNoise.Noise2D(x/settings.CaveLargeScale, y/settings.CaveLargeScale)
	horizontalTunnels := caveNoise.Noise2D(x/settings.CaveHorizontalScale+settings.CaveHorizontalYOffset, y/settings.CaveHorizontalYScale+settings.CaveHorizontalYOffset)
	verticalShafts := caveNoise.Noise2D(x/settings.CaveVerticalScale+settings.CaveVerticalYOffset, y/settings.CaveVerticalYScale+settings.CaveVerticalYOffset)
	smallCaves := caveNoise.Noise2D(x/settings.CaveSmallScale+settings.CaveSmallYOffset, y/settings.CaveSmallScale+settings.CaveSmallYOffset)
	airPockets := caveNoise.Noise2D(x/settings.CaveAirPocketScale+settings.CaveAirPocketYOffset, y/settings.CaveAirPocketScale+settings.CaveAirPocketYOffset)

	// Blend horizontal and vertical tunnels for more cross-connections
	tunnelBlend := (horizontalTunnels + verticalShafts) * 0.5

	// Depth-based cave generation with different types
	if depth > settings.CaveVeryDeepDepth {
		// Very deep - large caverns and complex systems
		largeCaverns := largeCaveNoise*0.25 + tunnelBlend*0.35 + smallCaves*0.25 + airPockets*0.15
		if largeCaverns > 0.08 {
			return true
		}

		// Additional tunnel networks deep underground
		deepTunnels := tunnelBlend*0.5 + smallCaves*0.3 + airPockets*0.2
		return deepTunnels > 0.18

	} else if depth > settings.CaveDeepDepth {
		// Deep caves - mix of large and medium caves
		deepCaves := largeCaveNoise*0.18 + tunnelBlend*0.45 + smallCaves*0.25 + airPockets*0.12
		if deepCaves > 0.10 {
			return true
		}

		// Vertical connections between levels
		return tunnelBlend > 0.18

	} else if depth > settings.CaveMediumDepth {
		// Medium depth - primarily horizontal tunnel systems
		mediumCaves := tunnelBlend*0.5 + smallCaves*0.35 + airPockets*0.15
		if mediumCaves > 0.13 {
			return true
		}

		// Some vertical shafts connecting to surface
		return tunnelBlend > 0.22

	} else if depth > settings.CaveShallowDepth { // more shallow caves
		// Medium-shallow caves - tunnel systems closer to surface
		mediumShallowCaves := tunnelBlend*0.45 + smallCaves*0.4 + airPockets*0.15
		if mediumShallowCaves > 0.10 {
			return true
		}

		// More vertical connections to surface
		return tunnelBlend > 0.18

	} else if depth > settings.CaveMinShallowDepth {
		// Shallow caves - small pockets and tunnels near surface
		shallowCaves := tunnelBlend*0.4 + smallCaves*0.45 + airPockets*0.15
		return shallowCaves > 0.09 // much more shallow caves
	}

	return false
}

// IsLargeCavern determines if a position should be part of a large underground cavern
func IsLargeCavern(worldX, worldY int) bool {
	surfaceHeight := GetHeightAt(worldX)
	depth := worldY - surfaceHeight

	// Allow large caverns closer to surface
	if depth < 40 { // Reduced from 80
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
	threshold := 0.6 - float64(depth-40)*0.001 // Adjusted for new minimum depth
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

// IsSurfaceCaveEntrance specifically creates visible cave entrances at the surface
func IsSurfaceCaveEntrance(worldX, worldY int) bool {
	surfaceHeight := GetHeightAt(worldX)

	// Only check positions at or just below surface
	if worldY < surfaceHeight || worldY > surfaceHeight+3 {
		return false
	}

	caveNoise := GetCaveNoise()
	x := float64(worldX)

	// Use a different noise pattern for entrance placement
	entranceNoise := caveNoise.Noise2D(x/40.0+9000, 0)

	// Make entrances more likely in hilly areas
	terrainNoise := GetTerrainNoise()
	hilliness := terrainNoise.Noise2D(x/30.0, 0)

	// Combine entrance noise with terrain variation
	combinedNoise := entranceNoise*0.7 + hilliness*0.3

	// More surface entrances
	return combinedNoise > 0.6 // Reduced from 0.7
}
