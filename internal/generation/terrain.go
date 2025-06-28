package generation

import (
	"github.com/aquilax/go-perlin"

	"github.com/KdntNinja/webcraft/internal/core/settings"
)

var (
	terrainHeights = make(map[int]int)
	terrainSeed    int64
	noiseGen       *perlin.Perlin
	caveNoise      *perlin.Perlin
	oreNoise       *perlin.Perlin
)

func init() {
	// Set the function variable in chunks.go to avoid import cycle
	GetWorldSeedFunc = GetWorldSeed
}

func initTerrainNoise() {
	if noiseGen == nil {
		noiseGen = perlin.NewPerlin(settings.PerlinAlpha, settings.PerlinBeta, settings.PerlinOctaves, terrainSeed)
		caveNoise = perlin.NewPerlin(settings.PerlinAlpha, settings.PerlinBeta, settings.PerlinOctaves, terrainSeed+1000)
		oreNoise = perlin.NewPerlin(settings.PerlinAlpha, settings.PerlinBeta, settings.PerlinOctaves, terrainSeed+2000)
	}
}

// ResetWorldGeneration forces regeneration with a new provided seed
func ResetWorldGeneration(seed int64) {
	terrainSeed = seed
	terrainHeights = make(map[int]int)
	noiseGen = nil
	caveNoise = nil
	oreNoise = nil
	initTerrainNoise()
}

// GetWorldSeed returns the current world seed for display or saving
func GetWorldSeed() int64 {
	initTerrainNoise()
	return terrainSeed
}

// GetHeightAt calculates the terrain height at a given world X coordinate
func GetHeightAt(worldX int) int {
	// Check cache first
	if height, exists := terrainHeights[worldX]; exists {
		return height
	}

	initTerrainNoise()

	// Generate height using multiple noise layers for realistic terrain
	x := float64(worldX)

	// Base terrain height using low-frequency noise
	baseHeight := noiseGen.Noise2D(x/100.0, 0)

	// Add hills and valleys with medium-frequency noise
	hillHeight := noiseGen.Noise2D(x/50.0, 100)

	// Add small details with high-frequency noise
	detailHeight := noiseGen.Noise2D(x/20.0, 200)

	// Combine noise layers
	combinedNoise := baseHeight*0.6 + hillHeight*0.3 + detailHeight*0.1

	// Scale and offset to world coordinates
	height := int(float64(settings.SurfaceBaseHeight) + combinedNoise*float64(settings.SurfaceHeightVar))

	// Ensure height is within reasonable bounds
	minHeight := 20                                              // Keep surface well above bedrock
	maxHeight := settings.ChunkHeight*settings.WorldChunksY - 50 // Leave room above
	if height < minHeight {
		height = minHeight
	}
	if height > maxHeight {
		height = maxHeight
	}

	// Cache the result
	terrainHeights[worldX] = height

	return height
}

// IsCave determines if a position should be a cave using 3D-like noise
func IsCave(worldX, worldY int) bool {
	// Only generate caves underground (below surface + some buffer)
	surfaceHeight := GetHeightAt(worldX)
	if worldY <= surfaceHeight+8 {
		return false
	}

	initTerrainNoise()

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

// GetOreType determines what type of ore (if any) should be at a position
func GetOreType(worldX, worldY int) int {
	// Only generate ores underground
	surfaceHeight := GetHeightAt(worldX)
	if worldY <= surfaceHeight+10 {
		return 0 // No ore near surface
	}

	initTerrainNoise()

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

	initTerrainNoise()

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
