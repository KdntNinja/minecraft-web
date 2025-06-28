package trees

import (
	"fmt"
	"math/rand"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// GenerateTreeAtPosition generates a tree at the specified position in a chunk
func GenerateTreeAtPosition(chunk *block.Chunk, x, surfaceChunkY int, rng *rand.Rand) {
	treeType, shape := GetTreeTypeAndShape(rng)
	shape.ValidateShape()

	fmt.Printf("TREE_DEBUG: Placing %v tree (height %d) at x=%d, surfaceChunkY=%d\n",
		treeType, shape.TrunkHeight, x, surfaceChunkY)

	// Use custom pattern if available
	if shape.CustomPattern != nil {
		shape.CustomPattern(chunk, x, surfaceChunkY, rng, shape)
		return
	}

	// Generate standard tree pattern
	generateStandardTree(chunk, x, surfaceChunkY, rng, treeType, shape)
}

// generateStandardTree creates the basic tree structure
func generateStandardTree(chunk *block.Chunk, x, surfaceChunkY int, rng *rand.Rand, treeType TreeType, shape TreeShape) {
	// Place trunk
	trunkBlock := treeType.GetTrunkBlock()
	for trunkLevel := 0; trunkLevel < shape.TrunkHeight; trunkLevel++ {
		trunkChunkY := surfaceChunkY - trunkLevel
		if trunkChunkY >= 0 && trunkChunkY < settings.ChunkHeight {
			chunk[trunkChunkY][x] = trunkBlock
		}
	}

	// Place branches if the tree has them
	if shape.HasBranches {
		generateBranches(chunk, x, surfaceChunkY, shape, rng)
	}

	// Place leaves (only if not a dead tree)
	if treeType != DeadTree && shape.LeafLayers > 0 {
		generateLeaves(chunk, x, surfaceChunkY, shape, treeType, rng)
	}
}

// generateBranches creates branch structures
func generateBranches(chunk *block.Chunk, x, surfaceChunkY int, shape TreeShape, rng *rand.Rand) {
	// Add branches at different heights
	branchHeight := shape.TrunkHeight / 2
	for i := 0; i < 2; i++ {
		branchY := surfaceChunkY - branchHeight - i
		if branchY >= 0 && branchY < settings.ChunkHeight {
			// Left branch
			if x > 0 && rng.Float64() < 0.7 {
				chunk[branchY][x-1] = block.Wood
				if shape.BranchLength > 1 && x > 1 && rng.Float64() < 0.5 {
					chunk[branchY][x-2] = block.Wood
				}
			}
			// Right branch
			if x < settings.ChunkWidth-1 && rng.Float64() < 0.7 {
				chunk[branchY][x+1] = block.Wood
				if shape.BranchLength > 1 && x < settings.ChunkWidth-2 && rng.Float64() < 0.5 {
					chunk[branchY][x+2] = block.Wood
				}
			}
		}
	}
}

// generateLeaves creates natural-looking leaf canopies with proper distribution
func generateLeaves(chunk *block.Chunk, x, surfaceChunkY int, shape TreeShape, treeType TreeType, rng *rand.Rand) {
	leafStartLevel := shape.TrunkHeight
	leafEndLevel := shape.TrunkHeight + shape.LeafLayers
	primaryLeafBlock := treeType.GetLeafBlock()

	// Generate leaves layer by layer from bottom to top
	for leafLevel := leafStartLevel; leafLevel <= leafEndLevel; leafLevel++ {
		leafChunkY := surfaceChunkY - leafLevel
		if leafChunkY < 0 || leafChunkY >= settings.ChunkHeight {
			continue
		}

		// Calculate layer properties
		layerFromBottom := leafLevel - leafStartLevel
		totalLayers := shape.LeafLayers

		// Generate natural leaf pattern based on tree type and layer
		switch treeType {
		case NormalTree, TallTree:
			generateNaturalLeafLayer(chunk, x, leafChunkY, layerFromBottom, totalLayers, shape, primaryLeafBlock, rng)
		case BushyTree:
			generateBushyLeafLayer(chunk, x, leafChunkY, layerFromBottom, totalLayers, shape, primaryLeafBlock, rng)
		case WideTree:
			generateWideLeafLayer(chunk, x, leafChunkY, layerFromBottom, totalLayers, shape, primaryLeafBlock, rng)
		case TwinTree:
			generateTwinLeafLayer(chunk, x, leafChunkY, layerFromBottom, totalLayers, shape, primaryLeafBlock, rng)
		case FlowerTree:
			generateFlowerLeafLayer(chunk, x, leafChunkY, layerFromBottom, totalLayers, shape, primaryLeafBlock, rng)
		default:
			generateNaturalLeafLayer(chunk, x, leafChunkY, layerFromBottom, totalLayers, shape, primaryLeafBlock, rng)
		}
	}

	// Add branch-end leaves if tree has branches
	if shape.HasBranches {
		generateBranchLeaves(chunk, x, surfaceChunkY, shape, primaryLeafBlock, rng)
	}
}

// generateNaturalLeafLayer creates a natural circular/oval leaf pattern
func generateNaturalLeafLayer(chunk *block.Chunk, centerX, chunkY, layerFromBottom, totalLayers int, shape TreeShape, leafBlock block.BlockType, rng *rand.Rand) {
	// Calculate radius based on layer position (wider in middle)
	progress := float64(layerFromBottom) / float64(totalLayers)

	var baseRadius float64
	if progress < 0.5 {
		// Growing phase (bottom half)
		baseRadius = progress * 2.0 * float64(shape.LeafWidth+1)
	} else {
		// Shrinking phase (top half)
		baseRadius = (2.0 - progress*2.0) * float64(shape.LeafWidth+1)
	}

	// Add randomness to radius
	radiusVariation := rng.Float64()*0.5 + 0.75 // 0.75 to 1.25 multiplier
	radius := baseRadius * radiusVariation

	// Generate leaves in circular pattern
	maxRadius := int(radius) + 1
	for dx := -maxRadius; dx <= maxRadius; dx++ {
		leafX := centerX + dx
		if leafX < 0 || leafX >= settings.ChunkWidth {
			continue
		}

		// Calculate distance from center
		distance := float64(abs(dx))

		// Use probability based on distance and some noise
		leafProbability := 1.0 - (distance / (radius + 1.0))

		// Add noise for natural edge variation
		noiseValue := rng.Float64() * 0.4   // 0 to 0.4
		leafProbability += noiseValue - 0.2 // -0.2 to +0.2 adjustment

		// Higher density near center
		if distance <= 1.0 {
			leafProbability += 0.3
		}

		// Apply sparseness for certain tree types
		if shape.IsSparse && rng.Float64() < 0.25 {
			leafProbability -= 0.4
		}

		if leafProbability > 0.4 && rng.Float64() < leafProbability {
			chunk[chunkY][leafX] = leafBlock
		}
	}
}

// generateBushyLeafLayer creates dense, compact leaf patterns
func generateBushyLeafLayer(chunk *block.Chunk, centerX, chunkY, layerFromBottom, totalLayers int, shape TreeShape, leafBlock block.BlockType, rng *rand.Rand) {
	// Bushy trees have consistent width with high density
	width := shape.LeafWidth + 1

	for dx := -width; dx <= width; dx++ {
		leafX := centerX + dx
		if leafX < 0 || leafX >= settings.ChunkWidth {
			continue
		}

		// High density everywhere except very edges
		leafProbability := 0.9
		if abs(dx) == width {
			leafProbability = 0.6
		}

		// Add some randomness
		leafProbability += rng.Float64()*0.2 - 0.1

		if rng.Float64() < leafProbability {
			chunk[chunkY][leafX] = leafBlock
		}
	}
}

// generateWideLeafLayer creates spreading, wide leaf patterns
func generateWideLeafLayer(chunk *block.Chunk, centerX, chunkY, layerFromBottom, totalLayers int, shape TreeShape, leafBlock block.BlockType, rng *rand.Rand) {
	// Wide trees get progressively wider, then narrow at top
	progress := float64(layerFromBottom) / float64(totalLayers)

	var width int
	if progress < 0.7 {
		// Expanding phase
		width = int(float64(shape.LeafWidth+2) * (progress / 0.7))
	} else {
		// Contracting phase
		width = int(float64(shape.LeafWidth+2) * (1.3 - progress) / 0.3)
	}

	if width < 1 {
		width = 1
	}

	for dx := -width; dx <= width; dx++ {
		leafX := centerX + dx
		if leafX < 0 || leafX >= settings.ChunkWidth {
			continue
		}

		// Probability decreases with distance from center
		distance := float64(abs(dx))
		leafProbability := 1.0 - (distance / float64(width+1))

		// Add variation
		leafProbability += rng.Float64()*0.3 - 0.15

		if rng.Float64() < leafProbability {
			chunk[chunkY][leafX] = leafBlock
		}
	}
}

// generateTwinLeafLayer creates two separate leaf clusters
func generateTwinLeafLayer(chunk *block.Chunk, centerX, chunkY, layerFromBottom, totalLayers int, shape TreeShape, leafBlock block.BlockType, rng *rand.Rand) {
	// Two clusters separated by 2-3 blocks
	separation := 2 + rng.Intn(2) // 2 or 3 blocks apart

	// Left cluster
	leftCenter := centerX - separation/2 - 1
	generateSmallCluster(chunk, leftCenter, chunkY, shape.LeafWidth, leafBlock, rng)

	// Right cluster
	rightCenter := centerX + separation/2 + 1
	generateSmallCluster(chunk, rightCenter, chunkY, shape.LeafWidth, leafBlock, rng)
}

// generateFlowerLeafLayer creates varied leaf patterns with flower blocks
func generateFlowerLeafLayer(chunk *block.Chunk, centerX, chunkY, layerFromBottom, totalLayers int, shape TreeShape, leafBlock block.BlockType, rng *rand.Rand) {
	// Similar to natural but with flower block variations
	generateNaturalLeafLayer(chunk, centerX, chunkY, layerFromBottom, totalLayers, shape, leafBlock, rng)

	// Add flower blocks (clay) scattered throughout
	width := shape.LeafWidth + 1
	for dx := -width; dx <= width; dx++ {
		leafX := centerX + dx
		if leafX < 0 || leafX >= settings.ChunkWidth {
			continue
		}

		// 15% chance to replace leaf with flower block
		if chunk[chunkY][leafX] == leafBlock && rng.Float64() < 0.15 {
			chunk[chunkY][leafX] = block.Clay // "Flowers"
		}
	}
}

// generateSmallCluster creates a small compact leaf cluster
func generateSmallCluster(chunk *block.Chunk, centerX, chunkY, maxWidth int, leafBlock block.BlockType, rng *rand.Rand) {
	width := 1 + rng.Intn(maxWidth)

	for dx := -width; dx <= width; dx++ {
		leafX := centerX + dx
		if leafX < 0 || leafX >= settings.ChunkWidth {
			continue
		}

		// High probability for small clusters
		if rng.Float64() < 0.8 {
			chunk[chunkY][leafX] = leafBlock
		}
	}
}

// generateBranchLeaves adds leaves at the end of branches
func generateBranchLeaves(chunk *block.Chunk, trunkX, surfaceChunkY int, shape TreeShape, leafBlock block.BlockType, rng *rand.Rand) {
	branchHeight := shape.TrunkHeight / 2

	for i := 0; i < 2; i++ {
		branchY := surfaceChunkY - branchHeight - i
		if branchY < 0 || branchY >= settings.ChunkHeight {
			continue
		}

		// Check for branches and add leaves at their ends
		for dx := -shape.BranchLength - 1; dx <= shape.BranchLength+1; dx++ {
			branchX := trunkX + dx
			if branchX < 0 || branchX >= settings.ChunkWidth {
				continue
			}

			// If there's a branch block, potentially add leaves around it
			if chunk[branchY][branchX] == block.Wood && abs(dx) > 1 {
				// Add leaves above branch
				if branchY > 0 && rng.Float64() < 0.7 {
					chunk[branchY-1][branchX] = leafBlock
				}
				// Add leaves beside branch end
				if abs(dx) >= shape.BranchLength && rng.Float64() < 0.5 {
					if branchX > 0 && branchX < settings.ChunkWidth-1 {
						if rng.Float64() < 0.5 {
							chunk[branchY][branchX-1] = leafBlock
						} else {
							chunk[branchY][branchX+1] = leafBlock
						}
					}
				}
			}
		}
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
