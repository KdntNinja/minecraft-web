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

	// Minecraft-style ore vein generation with different noise patterns for each ore

	// Copper ore - shallow, clustered veins
	copperVein := oreNoise.Noise2D(x/15.0, y/15.0)
	copperCluster := oreNoise.Noise2D(x/8.0+100, y/8.0+100)
	if depth > 15 && depth < 80 {
		copperCombined := copperVein*0.6 + copperCluster*0.4
		if copperCombined < -0.3 { // Made much more common (was -0.55)
			return 1 // Copper ore
		}
	}

	// Iron ore - medium depth, medium-sized veins
	ironVein := oreNoise.Noise2D(x/20.0+1000, y/20.0+1000)
	ironBranch := oreNoise.Noise2D(x/12.0+1100, y/12.0+1100)
	if depth > 25 && depth < 120 {
		ironCombined := ironVein*0.7 + ironBranch*0.3
		if ironCombined < -0.35 { // Made more common (was -0.6)
			return 2 // Iron ore
		}
	}

	// Gold ore - deep, rare, small but rich veins
	goldVein := oreNoise.Noise2D(x/30.0+2000, y/30.0+2000)
	goldPocket := oreNoise.Noise2D(x/6.0+2100, y/6.0+2100)
	if depth > 60 {
		goldCombined := goldVein*0.8 + goldPocket*0.2
		if goldCombined < -0.5 { // Made more common (was -0.72)
			return 3 // Gold ore
		}
	}

	return 0 // No ore
}

// GetOreVeinDensity returns how dense an ore vein should be (for clustering)
func GetOreVeinDensity(worldX, worldY, oreType int) float64 {
	oreNoise := GetOreNoise()
	x := float64(worldX)
	y := float64(worldY)

	// Different vein patterns for different ores
	switch oreType {
	case 1: // Copper - tight clusters
		return oreNoise.Noise2D(x/5.0+500, y/5.0+500)
	case 2: // Iron - medium spread
		return oreNoise.Noise2D(x/8.0+1500, y/8.0+1500)
	case 3: // Gold - small, concentrated pockets
		return oreNoise.Noise2D(x/4.0+2500, y/4.0+2500)
	default:
		return 0
	}
}

// IsOreVeinExtension checks if a position should extend an existing ore vein
func IsOreVeinExtension(worldX, worldY, oreType int) bool {
	density := GetOreVeinDensity(worldX, worldY, oreType)

	// Different thresholds for different ore types
	switch oreType {
	case 1: // Copper - extends easily
		return density > 0.3
	case 2: // Iron - moderate extension
		return density > 0.4
	case 3: // Gold - tight, concentrated veins
		return density > 0.5
	default:
		return false
	}
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

	// Water pools in medium depths - larger areas
	if depth > 30 && depth < 100 {
		waterNoise := oreNoise.Noise2D(x/25.0+3000, y/25.0+3000)
		waterSpread := oreNoise.Noise2D(x/15.0+3100, y/15.0+3100)
		waterCombined := waterNoise*0.7 + waterSpread*0.3
		if waterCombined < -0.75 {
			return 1 // Water
		}
	}

	// Lava pools in deep areas - smaller but more intense
	if depth > 80 {
		lavaNoise := oreNoise.Noise2D(x/20.0+4000, y/20.0+4000)
		lavaHeat := oreNoise.Noise2D(x/10.0+4100, y/10.0+4100)
		lavaCombined := lavaNoise*0.8 + lavaHeat*0.2
		if lavaCombined < -0.8 {
			return 2 // Lava (represented as Water for now)
		}
	}

	return 0 // No liquid
}
