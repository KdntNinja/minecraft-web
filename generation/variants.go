package generation

import (
	"math/rand"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/settings"
)

// GenerateTwinTree creates a tree with two trunks
func GenerateTwinTree(chunk *coretypes.Chunk, x, surfaceChunkY int, rng *rand.Rand, shape TreeShape) {
	// Place two trunks side by side
	trunks := []int{x, x + 1}
	if x == 0 {
		trunks = []int{x, x + 1}
	} else if x >= settings.ChunkWidth-1 {
		trunks = []int{x - 1, x}
	}

	// Generate both trunks
	for _, trunkX := range trunks {
		if trunkX >= 0 && trunkX < settings.ChunkWidth {
			for trunkLevel := 0; trunkLevel < shape.TrunkHeight; trunkLevel++ {
				trunkChunkY := surfaceChunkY - trunkLevel
				if trunkChunkY >= 0 && trunkChunkY < settings.ChunkHeight {
					chunk.Blocks[trunkChunkY][trunkX] = coretypes.Wood
				}
			}
		}
	}

	// Generate shared canopy
	for leafLevel := shape.TrunkHeight; leafLevel <= shape.TrunkHeight+shape.LeafLayers; leafLevel++ {
		leafChunkY := surfaceChunkY - leafLevel
		if leafChunkY >= 0 && leafChunkY < settings.ChunkHeight {
			for dx := -2; dx <= 2; dx++ {
				leafX := x + dx
				if leafX >= 0 && leafX < settings.ChunkWidth && rng.Float64() < 0.7 {
					chunk.Blocks[leafChunkY][leafX] = coretypes.Leaves
				}
			}
		}
	}
}

// GenerateGiantTree creates a massive tree structure
func GenerateGiantTree(chunk *coretypes.Chunk, x, surfaceChunkY int, rng *rand.Rand, shape TreeShape) {
	// Thick trunk (2x2 if possible)
	trunkPositions := []int{x}
	if x < settings.ChunkWidth-1 {
		trunkPositions = append(trunkPositions, x+1)
	}

	// Place thick trunk
	for _, trunkX := range trunkPositions {
		for trunkLevel := 0; trunkLevel < shape.TrunkHeight; trunkLevel++ {
			trunkChunkY := surfaceChunkY - trunkLevel
			if trunkChunkY >= 0 && trunkChunkY < settings.ChunkHeight {
				chunk.Blocks[trunkChunkY][trunkX] = coretypes.Wood
			}
		}
	}

	// Multiple branch layers
	for branchLayer := 0; branchLayer < 3; branchLayer++ {
		branchY := surfaceChunkY - (shape.TrunkHeight * 2 / 3) - branchLayer*2
		if branchY >= 0 && branchY < settings.ChunkHeight {
			for dx := -shape.BranchLength; dx <= shape.BranchLength; dx++ {
				branchX := x + dx
				if branchX >= 0 && branchX < settings.ChunkWidth && abs(dx) > 0 && rng.Float64() < 0.8 {
					chunk.Blocks[branchY][branchX] = coretypes.Wood
				}
			}
		}
	}

	// Massive canopy
	for leafLevel := shape.TrunkHeight; leafLevel <= shape.TrunkHeight+shape.LeafLayers; leafLevel++ {
		leafChunkY := surfaceChunkY - leafLevel
		if leafChunkY >= 0 && leafChunkY < settings.ChunkHeight {
			currentWidth := shape.LeafWidth - (leafLevel-shape.TrunkHeight)/2 // Taper towards top
			if currentWidth < 1 {
				currentWidth = 1
			}

			for dx := -currentWidth; dx <= currentWidth; dx++ {
				leafX := x + dx
				if leafX >= 0 && leafX < settings.ChunkWidth {
					leafProb := 0.8
					if abs(dx) == currentWidth {
						leafProb = 0.5
					}
					if rng.Float64() < leafProb {
						chunk.Blocks[leafChunkY][leafX] = coretypes.Leaves
					}
				}
			}
		}
	}
}

// GenerateSpookyTree creates a dead tree with twisted branches
func GenerateSpookyTree(chunk *coretypes.Chunk, x, surfaceChunkY int, rng *rand.Rand, shape TreeShape) {
	// Main trunk
	for trunkLevel := 0; trunkLevel < shape.TrunkHeight; trunkLevel++ {
		trunkChunkY := surfaceChunkY - trunkLevel
		if trunkChunkY >= 0 && trunkChunkY < settings.ChunkHeight {
			chunk.Blocks[trunkChunkY][x] = coretypes.Wood
		}
	}

	// Twisted branches at various heights
	for i := 0; i < shape.TrunkHeight/2; i++ {
		branchY := surfaceChunkY - shape.TrunkHeight + i*2
		if branchY >= 0 && branchY < settings.ChunkHeight {
			// Random twisted branches
			if rng.Float64() < 0.6 {
				direction := 1
				if rng.Float64() < 0.5 {
					direction = -1
				}
				branchX := x + direction
				if branchX >= 0 && branchX < settings.ChunkWidth {
					chunk.Blocks[branchY][branchX] = coretypes.Wood
					// Chance for branch extension
					if rng.Float64() < 0.4 {
						branchX += direction
						if branchX >= 0 && branchX < settings.ChunkWidth {
							chunk.Blocks[branchY][branchX] = coretypes.Wood
						}
					}
				}
			}
		}
	}
}

// GeneratePalmTree creates a palm-like tree with leaves only at the top
func GeneratePalmTree(chunk *coretypes.Chunk, x, surfaceChunkY int, rng *rand.Rand, shape TreeShape) {
	// Tall thin trunk
	for trunkLevel := 0; trunkLevel < shape.TrunkHeight; trunkLevel++ {
		trunkChunkY := surfaceChunkY - trunkLevel
		if trunkChunkY >= 0 && trunkChunkY < settings.ChunkHeight {
			chunk.Blocks[trunkChunkY][x] = coretypes.Wood
		}
	}

	// Leaves only at the very top in a palm frond pattern
	topY := surfaceChunkY - shape.TrunkHeight
	if topY >= 0 && topY < settings.ChunkHeight {
		// Center leaves
		chunk.Blocks[topY][x] = coretypes.Leaves

		// Frond-like leaves extending outward
		for dx := -2; dx <= 2; dx++ {
			leafX := x + dx
			if leafX >= 0 && leafX < settings.ChunkWidth && dx != 0 {
				if rng.Float64() < 0.8 {
					chunk.Blocks[topY][leafX] = coretypes.Leaves
				}
			}
		}

		// Some leaves one level up
		if topY > 0 {
			for dx := -1; dx <= 1; dx++ {
				leafX := x + dx
				if leafX >= 0 && leafX < settings.ChunkWidth {
					if rng.Float64() < 0.6 {
						chunk.Blocks[topY-1][leafX] = coretypes.Leaves
					}
				}
			}
		}
	}
}
