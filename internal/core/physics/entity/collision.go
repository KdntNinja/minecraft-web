package entity

import (
	"math"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

// IsSolid checks if a block at grid coordinates is solid
func IsSolid(blocks [][]int, x, y int) bool {
	if y < 0 || x < 0 || y >= len(blocks) || x >= len(blocks[0]) {
		return false
	}
	return blocks[y][x] > 0
}

// CollideBlocks handles AABB collision using resolv physics library
func (a *AABB) CollideBlocks(blocks [][]int) {
	// For now, let's use a simpler grid-based approach that works well
	// We can optimize this later with full resolv integration

	// Move horizontally first
	if a.VX != 0 {
		a.X += a.VX

		// Check for horizontal collision using math.Floor for precision
		if a.VX > 0 { // Moving right
			rightEdge := int(math.Floor((a.X + float64(a.Width)) / float64(block.TileSize)))
			for y := int(math.Floor(a.Y / float64(block.TileSize))); y <= int(math.Floor((a.Y+float64(a.Height)-1)/float64(block.TileSize))); y++ {
				if IsSolid(blocks, rightEdge, y) {
					a.X = float64(rightEdge*block.TileSize - a.Width)
					a.VX = 0
					break
				}
			}
		} else { // Moving left
			leftEdge := int(math.Floor(a.X / float64(block.TileSize)))
			for y := int(math.Floor(a.Y / float64(block.TileSize))); y <= int(math.Floor((a.Y+float64(a.Height)-1)/float64(block.TileSize))); y++ {
				if IsSolid(blocks, leftEdge, y) {
					a.X = float64((leftEdge + 1) * block.TileSize)
					a.VX = 0
					break
				}
			}
		}
	}

	// Move vertically second
	a.OnGround = false
	if a.VY != 0 {
		a.Y += a.VY

		// Check for vertical collision using math.Floor for precision
		if a.VY > 0 { // Moving down (falling)
			bottomEdge := int(math.Floor((a.Y + float64(a.Height)) / float64(block.TileSize)))
			for x := int(math.Floor(a.X / float64(block.TileSize))); x <= int(math.Floor((a.X+float64(a.Width)-1)/float64(block.TileSize))); x++ {
				if IsSolid(blocks, x, bottomEdge) {
					a.Y = float64(bottomEdge*block.TileSize - a.Height)
					a.VY = 0
					a.OnGround = true
					break
				}
			}
		} else { // Moving up (jumping)
			topEdge := int(math.Floor(a.Y / float64(block.TileSize)))
			for x := int(math.Floor(a.X / float64(block.TileSize))); x <= int(math.Floor((a.X+float64(a.Width)-1)/float64(block.TileSize))); x++ {
				if IsSolid(blocks, x, topEdge) {
					a.Y = float64((topEdge + 1) * block.TileSize)
					a.VY = 0
					break
				}
			}
		}
	}

	// Instant ground check for immediate settling using math.Floor
	if !a.OnGround && a.VY >= 0 {
		bottomY := a.Y + float64(a.Height)
		bottomEdge := int(math.Floor(bottomY / float64(block.TileSize)))

		// Instant settle if within 2 pixels of ground
		groundY := float64(bottomEdge * block.TileSize)
		if bottomY >= groundY && bottomY <= groundY+2.0 {
			for x := int(math.Floor(a.X / float64(block.TileSize))); x <= int(math.Floor((a.X+float64(a.Width)-1)/float64(block.TileSize))); x++ {
				if IsSolid(blocks, x, bottomEdge) {
					a.Y = groundY - float64(a.Height)
					a.VY = 0
					a.OnGround = true
					break
				}
			}
		}
	}

	// Instant stabilization when on ground using math.Floor for precision
	if a.OnGround {
		// Snap immediately to exact ground position
		bottomY := a.Y + float64(a.Height)
		bottomEdge := int(math.Floor(bottomY / float64(block.TileSize)))
		groundY := float64(bottomEdge * block.TileSize)

		// Snap if within 1 pixel for instant settling
		if bottomY > groundY && bottomY < groundY+1.0 {
			a.Y = groundY - float64(a.Height)
		}

		// Stop all vertical movement when grounded
		if a.VY > -0.1 && a.VY < 0.1 {
			a.VY = 0
		}
	}
}
