package physics

import (
	"math"

	"github.com/KdntNinja/webcraft/settings"
)

// IsSolid checks if a block at grid coordinates is solid, using grid offset
func IsSolid(blocks [][]int, x, y int, offsetX, offsetY int) bool {
	x -= offsetX
	y -= offsetY
	if y < 0 || x < 0 || y >= len(blocks) || x >= len(blocks[0]) {
		return false
	}
	return blocks[y][x] > 0
}

// CollideBlocks handles AABB collision using resolv physics library
func (a *AABB) CollideBlocks(blocks [][]int) {

	offsetX := a.GridOffsetX
	offsetY := a.GridOffsetY
	// Move horizontally first
	if a.VX != 0 {
		a.X += a.VX
		if a.VX > 0 { // Moving right
			// Subtract 1px to avoid extending into the next tile when precisely aligned
			rightEdge := int(math.Floor((a.X + float64(a.Width) - 1) / float64(settings.TileSize)))
			for y := int(math.Floor(a.Y / float64(settings.TileSize))); y <= int(math.Floor((a.Y+float64(a.Height)-1)/float64(settings.TileSize))); y++ {
				if IsSolid(blocks, rightEdge, y, offsetX, offsetY) {
					a.X = float64(rightEdge*settings.TileSize - a.Width)
					a.VX = 0
					break
				}
			}
		} else { // Moving left
			leftEdge := int(math.Floor(a.X / float64(settings.TileSize)))
			for y := int(math.Floor(a.Y / float64(settings.TileSize))); y <= int(math.Floor((a.Y+float64(a.Height)-1)/float64(settings.TileSize))); y++ {
				if IsSolid(blocks, leftEdge, y, offsetX, offsetY) {
					a.X = float64((leftEdge + 1) * settings.TileSize)
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
		if a.VY > 0 { // Moving down (falling)
			bottomEdge := int(math.Floor((a.Y + float64(a.Height)) / float64(settings.TileSize)))
			for x := int(math.Floor(a.X / float64(settings.TileSize))); x <= int(math.Floor((a.X+float64(a.Width)-1)/float64(settings.TileSize))); x++ {
				if IsSolid(blocks, x, bottomEdge, offsetX, offsetY) {
					a.Y = float64(bottomEdge*settings.TileSize - a.Height)
					a.VY = 0
					a.OnGround = true
					break
				}
			}
		} else { // Moving up (jumping)
			topEdge := int(math.Floor(a.Y / float64(settings.TileSize)))
			for x := int(math.Floor(a.X / float64(settings.TileSize))); x <= int(math.Floor((a.X+float64(a.Width)-1)/float64(settings.TileSize))); x++ {
				if IsSolid(blocks, x, topEdge, offsetX, offsetY) {
					a.Y = float64((topEdge + 1) * settings.TileSize)
					a.VY = 0
					break
				}
			}
		}
	}

	// Instant ground check for immediate settling using math.Floor
	if !a.OnGround && a.VY >= 0 {
		bottomY := a.Y + float64(a.Height)
		bottomEdge := int(math.Floor(bottomY / float64(settings.TileSize)))
		groundY := float64(bottomEdge * settings.TileSize)
		if bottomY >= groundY && bottomY <= groundY+2.0 {
			for x := int(math.Floor(a.X / float64(settings.TileSize))); x <= int(math.Floor((a.X+float64(a.Width)-1)/float64(settings.TileSize))); x++ {
				if IsSolid(blocks, x, bottomEdge, offsetX, offsetY) {
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
		bottomY := a.Y + float64(a.Height)
		bottomEdge := int(math.Floor(bottomY / float64(settings.TileSize)))
		groundY := float64(bottomEdge * settings.TileSize)
		if bottomY > groundY && bottomY < groundY+1.0 {
			a.Y = groundY - float64(a.Height)
		}
		if a.VY > -0.1 && a.VY < 0.1 {
			a.VY = 0
		}
	}
}
