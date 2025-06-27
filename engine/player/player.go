package player

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/entity"
)

const (
	Width        = block.TileSize
	Height       = block.TileSize * 2
	MoveSpeed    = 4.0   // More controlled Minecraft-like speed
	JumpSpeed    = -10.0 // Minecraft-like jump strength
	Gravity      = 0.5   // Gentler gravity like Minecraft
	MaxFallSpeed = 12.0  // Reasonable terminal velocity

	// Collision tolerance to prevent bouncing
	CollisionTolerance = 0.1
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

	// Jump - only when grounded and key pressed (with better input handling)
	if (ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeySpace)) && p.OnGround {
		p.VY = JumpSpeed
		p.OnGround = false
	}

	// Apply gravity with terminal velocity (more Minecraft-like)
	if !p.OnGround {
		p.VY += Gravity
		if p.VY > MaxFallSpeed {
			p.VY = MaxFallSpeed
		}
	} else {
		// When on ground, ensure minimal downward velocity to prevent bouncing
		if p.VY > 0 {
			p.VY = 0
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
