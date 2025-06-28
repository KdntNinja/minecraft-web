package terrain

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

var (
	terrainHeights = make(map[int]int)
	terrainSeed    int64
)

func init() {
	// Set the function variable in chunks.go to avoid import cycle
	GetWorldSeedFunc = GetWorldSeed
}

func initTerrainNoise() {
}

func getSurfaceHeight(x int) int {
	return 0
}

// getOreType determines what ore should be placed using simple noise
func getOreType(x, y, depthFromSurface int) block.BlockType {
	return block.Stone
}

// ResetWorldGeneration forces regeneration with a new provided seed
func ResetWorldGeneration(seed int64) {
	terrainSeed = seed
	terrainHeights = make(map[int]int)
}

// GetWorldSeed returns the current world seed for display or saving
func GetWorldSeed() int64 {
	initTerrainNoise()
	return terrainSeed
}

// getBlockType returns the block type for a given world position, generating Minecraft-like layers
func getBlockType(x, y int) block.BlockType {
	return block.Stone
}
