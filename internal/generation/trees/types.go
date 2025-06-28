package trees

import "github.com/KdntNinja/webcraft/internal/core/engine/block"

// TreeType represents different tree varieties
type TreeType int

const (
	NormalTree TreeType = iota
	TallTree
	BushyTree
	WideTree
	TwinTree   // Rare: two trunks
	GiantTree  // Very rare: massive tree
	DeadTree   // Rare: only trunk, no leaves
	FlowerTree // Rare: leaves mixed with other blocks
)

// String returns the name of the tree type
func (t TreeType) String() string {
	switch t {
	case NormalTree:
		return "Normal"
	case TallTree:
		return "Tall"
	case BushyTree:
		return "Bushy"
	case WideTree:
		return "Wide"
	case TwinTree:
		return "Twin"
	case GiantTree:
		return "Giant"
	case DeadTree:
		return "Dead"
	case FlowerTree:
		return "Flower"
	default:
		return "Unknown"
	}
}

// GetRarity returns the rarity of the tree type (0.0-1.0, lower is rarer)
func (t TreeType) GetRarity() float64 {
	switch t {
	case NormalTree:
		return 0.45 // 45%
	case TallTree:
		return 0.25 // 25%
	case BushyTree:
		return 0.15 // 15%
	case WideTree:
		return 0.07 // 7%
	case TwinTree:
		return 0.05 // 5%
	case DeadTree:
		return 0.02 // 2%
	case GiantTree:
		return 0.005 // 0.5%
	case FlowerTree:
		return 0.005 // 0.5%
	default:
		return 0.0
	}
}

// GetLeafBlock returns the primary leaf block type for this tree
func (t TreeType) GetLeafBlock() block.BlockType {
	switch t {
	case DeadTree:
		return block.Air // Dead trees have no leaves
	default:
		return block.Leaves
	}
}

// GetTrunkBlock returns the trunk block type for this tree
func (t TreeType) GetTrunkBlock() block.BlockType {
	return block.Wood // All trees use wood for now
}
