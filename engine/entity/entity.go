package entity

import "github.com/KdntNinja/webcraft/engine/block"

// Entity is the interface for all game entities (player, mobs, etc.)
type Entity interface {
	Update()
	CollideBlocks(blocks [][]int)
	ClampX(min, max float64)
	GetPosition() (float64, float64)
	SetPosition(x, y float64)
}

// IsSolid is a generic helper for all entities to check block solidity
func IsSolid(blocks [][]int, x, y int) bool {
	if y < 0 || x < 0 || y >= len(blocks) || x >= len(blocks[0]) {
		return false
	}
	// Only treat nonzero as solid, zero (air) is not a block at all
	return blocks[y][x] > 0
}

// Generic AABB collision logic for entities colliding with blocks
type AABB struct {
	X, Y          float64
	Width, Height int
	VX, VY        float64
	OnGround      bool
}

func (a *AABB) CollideBlocks(blocks [][]int) {
	// Move horizontally first
	if a.VX != 0 {
		a.X += a.VX

		// Check for horizontal collision and resolve
		if a.VX > 0 { // Moving right
			rightEdge := int((a.X + float64(a.Width) - 1) / float64(block.TileSize))
			for y := int(a.Y / float64(block.TileSize)); y <= int((a.Y+float64(a.Height)-1)/float64(block.TileSize)); y++ {
				if IsSolid(blocks, rightEdge, y) {
					a.X = float64(rightEdge*block.TileSize - a.Width)
					break
				}
			}
		} else { // Moving left
			leftEdge := int(a.X / float64(block.TileSize))
			for y := int(a.Y / float64(block.TileSize)); y <= int((a.Y+float64(a.Height)-1)/float64(block.TileSize)); y++ {
				if IsSolid(blocks, leftEdge, y) {
					a.X = float64((leftEdge + 1) * block.TileSize)
					break
				}
			}
		}
	}

	// Move vertically second
	a.OnGround = false
	if a.VY != 0 {
		a.Y += a.VY

		// Check for vertical collision and resolve
		if a.VY > 0 { // Moving down (falling)
			bottomEdge := int((a.Y + float64(a.Height) - 1) / float64(block.TileSize))
			for x := int(a.X / float64(block.TileSize)); x <= int((a.X+float64(a.Width)-1)/float64(block.TileSize)); x++ {
				if IsSolid(blocks, x, bottomEdge) {
					a.Y = float64(bottomEdge*block.TileSize - a.Height)
					a.VY = 0
					a.OnGround = true
					break
				}
			}
		} else { // Moving up (jumping)
			topEdge := int(a.Y / float64(block.TileSize))
			for x := int(a.X / float64(block.TileSize)); x <= int((a.X+float64(a.Width)-1)/float64(block.TileSize)); x++ {
				if IsSolid(blocks, x, topEdge) {
					a.Y = float64((topEdge + 1) * block.TileSize)
					a.VY = 0
					break
				}
			}
		}
	}
}

// Generic AABB helpers for all entities
func ClampX(x *float64, min, max float64) {
	if *x < min {
		*x = min
	}
	if *x > max {
		*x = max
	}
}

func GetPosition(x, y float64) (float64, float64) {
	return x, y
}

func SetPosition(x, y float64, nx, ny float64) (float64, float64) {
	return nx, ny
}

// ApplyGravity applies gravity to the entity's vertical velocity
func ApplyGravity(vy *float64, gravity float64) {
	*vy += gravity
}
