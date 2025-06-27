package entity

import (
	"math"

	"github.com/KdntNinja/webcraft/engine/block"
)

// IsSolid is a generic helper for all entities to check block solidity
func IsSolid(blocks [][]int, x, y int) bool {
	if y < 0 || x < 0 || y >= len(blocks) || x >= len(blocks[0]) {
		return false
	}
	// Only treat nonzero as solid, zero (air) is not a block at all
	return blocks[y][x] > 0
}

// CollideBlocks handles AABB collision logic for entities colliding with blocks
func (a *AABB) CollideBlocks(blocks [][]int) {
	// Move horizontally first with precise collision detection
	if a.VX != 0 {
		a.X += a.VX

		// Check for horizontal collision and resolve
		if a.VX > 0 { // Moving right
			rightEdge := int(math.Floor((a.X + float64(a.Width)) / float64(block.TileSize)))
			for y := int(math.Floor(a.Y / float64(block.TileSize))); y <= int(math.Floor((a.Y+float64(a.Height)-1)/float64(block.TileSize))); y++ {
				if IsSolid(blocks, rightEdge, y) {
					// Snap precisely to wall surface
					a.X = float64(rightEdge*block.TileSize - a.Width)
					a.VX = 0
					break
				}
			}
		} else { // Moving left
			leftEdge := int(math.Floor(a.X / float64(block.TileSize)))
			for y := int(math.Floor(a.Y / float64(block.TileSize))); y <= int(math.Floor((a.Y+float64(a.Height)-1)/float64(block.TileSize))); y++ {
				if IsSolid(blocks, leftEdge, y) {
					// Snap precisely to wall surface
					a.X = float64((leftEdge + 1) * block.TileSize)
					a.VX = 0
					break
				}
			}
		}
	}

	// Move vertically second with improved ground detection
	a.OnGround = false
	if a.VY != 0 {
		a.Y += a.VY

		// Check for vertical collision and resolve
		if a.VY > 0 { // Moving down (falling)
			bottomEdge := int(math.Floor((a.Y + float64(a.Height)) / float64(block.TileSize)))
			for x := int(math.Floor(a.X / float64(block.TileSize))); x <= int(math.Floor((a.X+float64(a.Width)-1)/float64(block.TileSize))); x++ {
				if IsSolid(blocks, x, bottomEdge) {
					// Precise ground positioning - snap exactly to block surface
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
					// Precise ceiling collision - snap exactly to block surface
					a.Y = float64((topEdge + 1) * block.TileSize)
					a.VY = 0
					break
				}
			}
		}
	}

	// Additional ground check for precise detection (prevents floaty jumping)
	if !a.OnGround && a.VY >= 0 {
		bottomY := a.Y + float64(a.Height)
		bottomEdge := int(math.Floor(bottomY / float64(block.TileSize)))

		// Check if we're touching or very close to the ground (within 1 pixel)
		groundY := float64(bottomEdge * block.TileSize)
		if bottomY >= groundY && bottomY <= groundY+1.0 { // Within 1 pixel of ground
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
}
