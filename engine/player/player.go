package player

import (
	"github.com/KdntNinja/webcraft/engine/entity"
)

type Player struct {
	entity.AABB
	lastGroundedTime int  // Frames since last grounded (for coyote time)
	jumpPressed      bool // Track jump key state to prevent spam
	wasOnGround      bool // Previous frame ground state
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		AABB: entity.AABB{
			X: x, Y: y, Width: Width, Height: Height,
		},
	}
}

func (p *Player) Update() {
	p.wasOnGround = p.OnGround // Store previous ground state

	// Process input and update movement
	isMoving, targetVX, jumpKeyPressed := p.HandleInput()
	p.ApplyMovement(isMoving, targetVX)
	p.HandleJump(jumpKeyPressed)
	p.ApplyGravity()
}

// Entity interface implementations (delegate to AABB)
func (p *Player) CollideBlocks(blocks [][]int) {
	p.AABB.CollideBlocks(blocks)
}

func (p *Player) ClampX(min, max float64) {
	p.AABB.ClampX(min, max)
}

func (p *Player) GetPosition() (float64, float64) {
	return p.AABB.GetPosition()
}

func (p *Player) SetPosition(x, y float64) {
	p.AABB.SetPosition(x, y)
}
