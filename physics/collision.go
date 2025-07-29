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

// CollideBlocks handles AABB collision via raycasting at the object edges
func (a *AABB) CollideBlocks(blocks [][]int) {
	offsetX := a.GridOffsetX
	offsetY := a.GridOffsetY

	// Horizontal stepping collision
	if a.VX != 0 {
		dx := a.VX
		signX := math.Copysign(1, dx)
		steps := int(math.Abs(dx))
		// step one pixel at a time
		for i := 0; i < steps; i++ {
			a.X += signX
			// check both top and bottom edge
			tileX := int(math.Floor((a.X + math.Max(0, signX*(float64(a.Width)-1))) / float64(settings.TileSize)))
			minY := int(math.Floor((a.Y + 1) / float64(settings.TileSize)))
			maxY := int(math.Floor((a.Y + float64(a.Height) - 1) / float64(settings.TileSize)))
			hit := false
			for y := minY; y <= maxY; y++ {
				if IsSolid(blocks, tileX, y, offsetX, offsetY) {
					hit = true
					break
				}
			}
			if hit {
				a.X -= signX
				a.VX = 0
				break
			}
		}
		// remaining fraction
		remX := dx - signX*float64(steps)
		a.X += remX
		// fractional collision check
		tileX := int(math.Floor((a.X + math.Max(0, remX+math.Max(0, float64(a.Width)-1))) / float64(settings.TileSize)))
		minY := int(math.Floor((a.Y + 1) / float64(settings.TileSize)))
		maxY := int(math.Floor((a.Y + float64(a.Height) - 1) / float64(settings.TileSize)))
		for y := minY; y <= maxY; y++ {
			if IsSolid(blocks, tileX, y, offsetX, offsetY) {
				a.X -= remX
				a.VX = 0
				break
			}
		}
	}

	// Vertical stepping collision
	a.OnGround = false
	if a.VY != 0 {
		dy := a.VY
		signY := math.Copysign(1, dy)
		steps := int(math.Abs(dy))
		for i := 0; i < steps; i++ {
			a.Y += signY
			// check both left and right edge
			tileY := int(math.Floor((a.Y + math.Max(0, signY*(float64(a.Height)-1))) / float64(settings.TileSize)))
			minX := int(math.Floor((a.X + 1) / float64(settings.TileSize)))
			maxX := int(math.Floor((a.X + float64(a.Width) - 1) / float64(settings.TileSize)))
			hit := false
			for x := minX; x <= maxX; x++ {
				if IsSolid(blocks, x, tileY, offsetX, offsetY) {
					hit = true
					break
				}
			}
			if hit {
				a.Y -= signY
				a.VY = 0
				if signY > 0 {
					a.OnGround = true
				}
				break
			}
		}
		remY := dy - signY*float64(steps)
		a.Y += remY
		tileY := int(math.Floor((a.Y + math.Max(0, remY+math.Max(0, float64(a.Height)-1))) / float64(settings.TileSize)))
		minX := int(math.Floor((a.X + 1) / float64(settings.TileSize)))
		maxX := int(math.Floor((a.X + float64(a.Width) - 1) / float64(settings.TileSize)))
		for x := minX; x <= maxX; x++ {
			if IsSolid(blocks, x, tileY, offsetX, offsetY) {
				a.Y -= remY
				a.VY = 0
				if signY > 0 {
					a.OnGround = true
				}
				break
			}
		}
	}
}
