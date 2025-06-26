package player

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/entity"
)

const (
	Width        = block.TileSize
	Height       = block.TileSize * 2
	MoveSpeed    = 4.0
	JumpSpeed    = -10.0
	Gravity      = 0.6
	MaxFallSpeed = 10.0
)

type Player struct {
	entity.AABB
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		AABB: entity.AABB{
			X: x, Y: y, Width: Width, Height: Height,
		},
	}
}

func (p *Player) Update() {
	// Handle horizontal movement - smooth and responsive
	p.VX = 0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.VX = -MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.VX = MoveSpeed
	}

	// Jump - only when grounded and key pressed
	if (ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeySpace)) && p.OnGround {
		p.VY = JumpSpeed
		p.OnGround = false
	}

	// Apply gravity with terminal velocity
	if !p.OnGround {
		p.VY += Gravity
		if p.VY > MaxFallSpeed {
			p.VY = MaxFallSpeed
		}
	}
}

func (p *Player) CollideBlocks(blocks [][]int) {
	p.AABB.CollideBlocks(blocks)
}

func (p *Player) ClampX(min, max float64) {
	entity.ClampX(&p.X, min, max)
}

func (p *Player) GetPosition() (float64, float64) {
	return entity.GetPosition(p.X, p.Y)
}

func (p *Player) SetPosition(x, y float64) {
	p.X, p.Y = entity.SetPosition(p.X, p.Y, x, y)
}
