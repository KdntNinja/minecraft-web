package generation

import "github.com/KdntNinja/webcraft/internal/core/settings"

var terrainHeights = make(map[int]int)

// GetHeightAt calculates the terrain height at a given world X coordinate
func GetHeightAt(worldX int) int {
	// Check cache first
	if height, exists := terrainHeights[worldX]; exists {
		return height
	}

	noise := GetTerrainNoise()

	// Generate height using multiple noise layers for realistic terrain
	x := float64(worldX)

	// Base terrain height using low-frequency noise
	baseHeight := noise.Noise2D(x/100.0, 0)

	// Add hills and valleys with medium-frequency noise
	hillHeight := noise.Noise2D(x/50.0, 100)

	// Add small details with high-frequency noise
	detailHeight := noise.Noise2D(x/20.0, 200)

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

// ResetHeightCache clears the terrain height cache
func ResetHeightCache() {
	terrainHeights = make(map[int]int)
}
