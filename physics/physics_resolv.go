package physics

import (
	"github.com/solarlune/resolv"

	"github.com/KdntNinja/webcraft/settings"
)

// PhysicsWorld manages the physics simulation using resolv
type PhysicsWorld struct {
	Space  *resolv.Space
	Blocks [][]int // Keep reference to blocks for simple collision detection
}

// NewPhysicsWorld creates a new physics world from block data
func NewPhysicsWorld(blocks [][]int) *PhysicsWorld {
	if len(blocks) == 0 || len(blocks[0]) == 0 {
		return &PhysicsWorld{
			Space:  resolv.NewSpace(800, 600, 16, 16),
			Blocks: blocks,
		}
	}

	height := len(blocks)
	width := len(blocks[0])
	tileSize := settings.TileSize

	// Create space with cell size matching tile size for optimal performance
	space := resolv.NewSpace(width*tileSize, height*tileSize, tileSize, tileSize)

	return &PhysicsWorld{
		Space:  space,
		Blocks: blocks,
	}
}
