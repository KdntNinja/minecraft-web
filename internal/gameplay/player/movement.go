package player

import (
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// ApplyMovement handles horizontal movement physics using entity movement system
func (p *Player) ApplyMovement(isMoving bool, targetVX float64) {
	p.AABB.ApplyHorizontalMovement(targetVX, settings.PlayerGroundFriction, settings.PlayerAirResistance, isMoving)

	// Apply velocity damping to prevent jitter when not moving
	if !isMoving {
		entity.DampVelocity(&p.VX, &p.VY, 0.1)
	}
}

// HandleJump processes jump input using entity input state tracking
func (p *Player) HandleJump() {
	// Only jump on fresh key press while grounded
	if p.InputState.CanJump() {
		if p.AABB.Jump(settings.PlayerJumpSpeed) {
			// Jump was successful
		}
	}
}

// ApplyGravity handles gravity and fall physics with instant settling
func (p *Player) ApplyGravity() {
	// Stronger gravity near ground for instant settling
	gravityToApply := settings.PlayerGravity
	if p.VY > 0 && p.VY < 3.0 { // When falling slowly
		gravityToApply = settings.PlayerGravity * 2.0 // Double gravity for quick settling
	}

	// Update grounded time tracking
	p.AABB.ApplyVerticalMovement(gravityToApply, settings.PlayerMaxFallSpeed)

	// Instant stabilization when on ground
	if p.OnGround {
		entity.StabilizePosition(&p.Y, settings.TileSize) // Snap to tile grid
		entity.DampVelocity(&p.VX, &p.VY, 0.3)            // Strong damping for instant settling
	}
}
