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
	baseHeight := noise.Noise2D(x/settings.TerrainBaseScale, 0)

	// Add hills and valleys with medium-frequency noise
	hillHeight := noise.Noise2D(x/settings.TerrainHillScale, settings.TerrainHillOffset)

	// Add small details with high-frequency noise
	detailHeight := noise.Noise2D(x/settings.TerrainDetailScale, settings.TerrainDetailOffset)

	// Combine noise layers
	combinedNoise := baseHeight*settings.TerrainBaseWeight + hillHeight*settings.TerrainHillWeight + detailHeight*settings.TerrainDetailWeight

	// Scale and offset to world coordinates
	height := int(float64(settings.SurfaceBaseHeight) + combinedNoise*float64(settings.SurfaceHeightVar))

	// Ensure height is within reasonable bounds
	minHeight := settings.TerrainMinHeight
	maxHeight := settings.ChunkHeight*settings.WorldChunksY - settings.TerrainMaxHeightBuffer
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
