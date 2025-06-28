package trees

import (
	"math/rand"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

// TreeShape defines the shape and size parameters for a tree
type TreeShape struct {
	TrunkHeight   int
	LeafLayers    int
	LeafWidth     int
	HasBranches   bool
	BranchLength  int
	IsSparse      bool
	CustomPattern func(*block.Chunk, int, int, *rand.Rand, TreeShape)
}

// GetTreeTypeAndShape determines what type of tree to generate with weighted randomness
func GetTreeTypeAndShape(rng *rand.Rand) (TreeType, TreeShape) {
	roll := rng.Float64()
	cumulative := 0.0

	// Check each tree type in order of rarity (most common first)
	treeTypes := []TreeType{NormalTree, TallTree, BushyTree, WideTree, TwinTree, DeadTree, GiantTree, FlowerTree}

	for _, treeType := range treeTypes {
		cumulative += treeType.GetRarity()
		if roll < cumulative {
			return treeType, generateShapeForType(treeType, rng)
		}
	}

	// Fallback to normal tree
	return NormalTree, generateShapeForType(NormalTree, rng)
}

// generateShapeForType creates appropriate shape parameters for a given tree type
func generateShapeForType(treeType TreeType, rng *rand.Rand) TreeShape {
	switch treeType {
	case NormalTree:
		return TreeShape{
			TrunkHeight: 2 + rng.Intn(3), // 2-4 blocks
			LeafLayers:  2 + rng.Intn(2), // 2-3 layers
			LeafWidth:   1,
			HasBranches: false,
		}
	case TallTree:
		return TreeShape{
			TrunkHeight:  4 + rng.Intn(4), // 4-7 blocks
			LeafLayers:   3 + rng.Intn(2), // 3-4 layers
			LeafWidth:    1,
			HasBranches:  rng.Float64() < 0.3, // 30% chance of branches
			BranchLength: 1,
		}
	case BushyTree:
		return TreeShape{
			TrunkHeight: 2 + rng.Intn(2), // 2-3 blocks
			LeafLayers:  3 + rng.Intn(2), // 3-4 layers
			LeafWidth:   2,               // Wider leaves
			HasBranches: false,
		}
	case WideTree:
		return TreeShape{
			TrunkHeight:  3 + rng.Intn(2), // 3-4 blocks
			LeafLayers:   2,
			LeafWidth:    2,
			HasBranches:  true,
			BranchLength: 2,
		}
	case TwinTree:
		return TreeShape{
			TrunkHeight:   3 + rng.Intn(3), // 3-5 blocks
			LeafLayers:    3,
			LeafWidth:     1,
			HasBranches:   false,
			CustomPattern: GenerateTwinTree,
		}
	case DeadTree:
		return TreeShape{
			TrunkHeight:  3 + rng.Intn(4),     // 3-6 blocks
			LeafLayers:   0,                   // No leaves
			HasBranches:  rng.Float64() < 0.7, // 70% chance of bare branches
			BranchLength: 1 + rng.Intn(2),
		}
	case GiantTree:
		return TreeShape{
			TrunkHeight:   8 + rng.Intn(4), // 8-11 blocks
			LeafLayers:    5 + rng.Intn(2), // 5-6 layers
			LeafWidth:     3,               // Very wide
			HasBranches:   true,
			BranchLength:  2 + rng.Intn(2),
			CustomPattern: GenerateGiantTree,
		}
	case FlowerTree:
		return TreeShape{
			TrunkHeight: 2 + rng.Intn(3), // 2-4 blocks
			LeafLayers:  2 + rng.Intn(2), // 2-3 layers
			LeafWidth:   1,
			HasBranches: false,
			IsSparse:    true, // Some leaves replaced with other blocks
		}
	default:
		return generateShapeForType(NormalTree, rng)
	}
}

// ValidateShape ensures tree shape parameters are within reasonable bounds
func (s *TreeShape) ValidateShape() {
	if s.TrunkHeight < 1 {
		s.TrunkHeight = 1
	}
	if s.TrunkHeight > 15 {
		s.TrunkHeight = 15
	}
	if s.LeafLayers < 0 {
		s.LeafLayers = 0
	}
	if s.LeafLayers > 8 {
		s.LeafLayers = 8
	}
	if s.LeafWidth < 0 {
		s.LeafWidth = 0
	}
	if s.LeafWidth > 4 {
		s.LeafWidth = 4
	}
	if s.BranchLength < 0 {
		s.BranchLength = 0
	}
	if s.BranchLength > 3 {
		s.BranchLength = 3
	}
}
