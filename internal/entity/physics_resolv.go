package entity

import (
	"github.com/KdntNinja/webcraft/internal/block"
	"github.com/solarlune/resolv"
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
	tileSize := block.TileSize

	// Create space with cell size matching tile size for optimal performance
	space := resolv.NewSpace(width*tileSize, height*tileSize, tileSize, tileSize)

	return &PhysicsWorld{
		Space:  space,
		Blocks: blocks,
	}
}

// CollideBlocksAdvanced uses improved collision detection with sub-pixel precision
func (a *AABB) CollideBlocksAdvanced(world *PhysicsWorld) {
	blocks := world.Blocks
	tileSize := float64(block.TileSize)

	// Move horizontally first with sub-stepping for better precision
	if a.VX != 0 {
		steps := 1
		if a.VX > tileSize/2 || a.VX < -tileSize/2 {
			steps = int(abs(a.VX)/(tileSize/2)) + 1 // Sub-step for fast movement
		}

		stepX := a.VX / float64(steps)

		for i := 0; i < steps; i++ {
			newX := a.X + stepX

			// Check collision at new position
			if a.VX > 0 { // Moving right
				rightEdge := int((newX + float64(a.Width)) / tileSize)
				for y := int(a.Y / tileSize); y <= int((a.Y+float64(a.Height)-1)/tileSize); y++ {
					if IsSolid(blocks, rightEdge, y) {
						a.X = float64(rightEdge)*tileSize - float64(a.Width) - 0.1 // Slight offset to prevent sticking
						a.VX = 0
						break
					}
				}
				if a.VX == 0 {
					break // Stop if collision found
				}
				a.X = newX
			} else { // Moving left
				leftEdge := int(newX / tileSize)
				for y := int(a.Y / tileSize); y <= int((a.Y+float64(a.Height)-1)/tileSize); y++ {
					if IsSolid(blocks, leftEdge, y) {
						a.X = float64(leftEdge+1)*tileSize + 0.1 // Slight offset to prevent sticking
						a.VX = 0
						break
					}
				}
				if a.VX == 0 {
					break // Stop if collision found
				}
				a.X = newX
			}
		}
	}

	// Move vertically with sub-stepping
	a.OnGround = false
	if a.VY != 0 {
		steps := 1
		if a.VY > tileSize/2 || a.VY < -tileSize/2 {
			steps = int(abs(a.VY)/(tileSize/2)) + 1 // Sub-step for fast movement
		}

		stepY := a.VY / float64(steps)

		for i := 0; i < steps; i++ {
			newY := a.Y + stepY

			// Check collision at new position
			if a.VY > 0 { // Moving down (falling)
				bottomEdge := int((newY + float64(a.Height)) / tileSize)
				for x := int(a.X / tileSize); x <= int((a.X+float64(a.Width)-1)/tileSize); x++ {
					if IsSolid(blocks, x, bottomEdge) {
						a.Y = float64(bottomEdge)*tileSize - float64(a.Height) - 0.1 // Slight offset
						a.VY = 0
						a.OnGround = true
						break
					}
				}
				if a.VY == 0 {
					break // Stop if collision found
				}
				a.Y = newY
			} else { // Moving up (jumping)
				topEdge := int(newY / tileSize)
				for x := int(a.X / tileSize); x <= int((a.X+float64(a.Width)-1)/tileSize); x++ {
					if IsSolid(blocks, x, topEdge) {
						a.Y = float64(topEdge+1)*tileSize + 0.1 // Slight offset
						a.VY = 0
						break
					}
				}
				if a.VY == 0 {
					break // Stop if collision found
				}
				a.Y = newY
			}
		}
	}

	// Ground check for when not moving vertically
	if !a.OnGround && a.VY >= 0 {
		bottomY := a.Y + float64(a.Height)
		bottomEdge := int((bottomY + 2.0) / tileSize) // Check slightly below

		for x := int(a.X / tileSize); x <= int((a.X+float64(a.Width)-1)/tileSize); x++ {
			if IsSolid(blocks, x, bottomEdge) {
				// Check if we're close enough to the ground
				groundY := float64(bottomEdge * block.TileSize)
				if bottomY >= groundY-2.0 && bottomY <= groundY+2.0 {
					a.OnGround = true
					break
				}
			}
		}
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
