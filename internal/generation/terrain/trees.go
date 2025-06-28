package terrain

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation/noise"
)

// generateTree creates a simple tree structure at the given position
func generateTree(chunk *block.Chunk, surfaceX, surfaceY int, biome noise.BiomeData, terrainNoise *noise.PerlinNoise, globalX int) {
	// Simple deterministic tree height based on biome and position
	heightNoise := terrainNoise.Noise1D(float64(globalX) * 0.1)
	var baseHeight int

	switch biome.Type {
	case 1: // ForestBiome - medium trees
		baseHeight = 4 + int(heightNoise*3) // 4-7 blocks
	case 6: // JungleBiome - tall trees
		baseHeight = 6 + int(heightNoise*4) // 6-10 blocks
	case 0: // PlainsBiome - small trees
		baseHeight = 3 + int(heightNoise*2) // 3-5 blocks
	case 4: // SwampBiome - medium trees
		baseHeight = 4 + int(heightNoise*2) // 4-6 blocks
	case 3, 5: // MountainBiome, TundraBiome - small hardy trees
		baseHeight = 3 + int(heightNoise*2) // 3-5 blocks
	default:
		baseHeight = 4 + int(heightNoise*2) // 4-6 blocks
	}

	if baseHeight < 3 {
		baseHeight = 3
	}
	if baseHeight > 10 {
		baseHeight = 10
	}

	// Generate trunk
	for i := 1; i <= baseHeight; i++ {
		trunkY := surfaceY - i
		if trunkY >= 0 && trunkY < settings.ChunkHeight {
			(*chunk)[trunkY][surfaceX] = block.Wood
		}
	}

	// Generate simple canopy
	canopyY := surfaceY - baseHeight
	canopySize := 2 // Simple 5x5 canopy

	for layer := 0; layer < 3; layer++ {
		layerY := canopyY + layer
		if layerY >= 0 && layerY < settings.ChunkHeight {
			for dx := -canopySize; dx <= canopySize; dx++ {
				leafX := surfaceX + dx
				if leafX >= 0 && leafX < settings.ChunkWidth {
					// Skip corners and center (where trunk is)
					if layer == 0 && (dx == 0) {
						continue // Don't overwrite trunk
					}
					if (dx == -canopySize || dx == canopySize) && layer == 2 {
						// Randomly skip some edge leaves for natural look
						edgeNoise := terrainNoise.Noise2D(float64(globalX+dx)*0.3, float64(layerY)*0.3)
						if edgeNoise < 0.3 {
							continue
						}
					}
					(*chunk)[layerY][leafX] = block.Leaves
				}
			}
		}
	}
}

// Tree placement logic
func shouldPlaceTree(globalX int, biome noise.BiomeData) bool {
	// Calculate tree chance based on biome - updated for 5 biomes
	var treeChance float64
	switch biome.Type {
	case 0: // PlainsBiome
		treeChance = 0.02 // Sparse trees
	case 1: // ForestBiome
		treeChance = 0.15 // Dense forest
	case 2: // MountainBiome (renumbered)
		treeChance = 0.05 // Sparse mountain trees
	case 3: // JungleBiome (renumbered)
		treeChance = 0.20 // Very dense jungle
	case 4: // OceanBiome (renumbered)
		treeChance = 0.0 // No trees in ocean
	default:
		treeChance = 0.02
	}

	// Use simple hash-based tree placement for deterministic results
	hash := float64(((globalX*73856093)^(globalX*19349663))%1000000) / 1000000.0
	return hash < treeChance
}

func shouldPlaceTreeByID(globalX int, biomeID int) bool {
	// Calculate tree chance based on biome ID - updated for 5 biomes
	var treeChance float64
	switch biomeID {
	case 0: // PlainsBiome
		treeChance = 0.02
	case 1: // ForestBiome
		treeChance = 0.15
	case 2: // MountainBiome (renumbered)
		treeChance = 0.05
	case 3: // JungleBiome (renumbered)
		treeChance = 0.20
	case 4: // OceanBiome (renumbered)
		treeChance = 0.0 // No trees in ocean
	default:
		treeChance = 0.02
	}

	// Use simple hash-based tree placement for deterministic results
	hash := float64(((globalX*73856093)^(globalX*19349663))%1000000) / 1000000.0
	return hash < treeChance
}
