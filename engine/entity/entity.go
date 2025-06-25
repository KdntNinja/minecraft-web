package entity

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
	return blocks[y][x] != 0 // 0 = Air, nonzero = solid
}

// Generic AABB collision logic for entities colliding with blocks
type AABB struct {
	X, Y          float64
	Width, Height int
	VX, VY        float64
	OnGround      bool
}

func (a *AABB) CollideBlocks(blocks [][]int) {
	a.OnGround = false
	px0 := int(a.X) / a.Width
	py0 := int(a.Y) / a.Height
	px1 := int(a.X+float64(a.Width)-1) / a.Width
	py1 := int(a.Y+float64(a.Height)-1) / a.Height

	for y := py0; y <= py1; y++ {
		for x := px0; x <= px1; x++ {
			if IsSolid(blocks, x, y) {
				if a.VY > 0 && int(a.Y+float64(a.Height)) > y*a.Height && int(a.Y) < (y+1)*a.Height {
					a.Y = float64(y*a.Height - a.Height)
					a.VY = 0
					a.OnGround = true
				}
				if a.VY < 0 && int(a.Y) < (y+1)*a.Height && int(a.Y+float64(a.Height)) > y*a.Height {
					a.Y = float64((y + 1) * a.Height)
					a.VY = 0
				}
				if a.VX > 0 && int(a.X+float64(a.Width)) > x*a.Width && int(a.X) < (x+1)*a.Width {
					a.X = float64(x*a.Width - a.Width)
					a.VX = 0
				}
				if a.VX < 0 && int(a.X) < (x+1)*a.Width && int(a.X+float64(a.Width)) > x*a.Width {
					a.X = float64((x + 1) * a.Width)
					a.VX = 0
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
