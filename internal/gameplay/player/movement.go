package player

import (
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// ApplyMovement handles horizontal movement physics using entity movement system
func (p *Player) ApplyMovement(isMoving bool, targetVX float64) {
	// Improved: Use acceleration and smoothing for more responsive and natural movement
	acceleration := 0.7
	maxSpeed := settings.PlayerMoveSpeed * 1.2
	if isMoving {
		if p.OnGround {
			// Accelerate towards targetVX
			p.VX += (targetVX - p.VX) * acceleration
		} else {
			// Air control: less acceleration
			p.VX += (targetVX - p.VX) * (acceleration * 0.4)
		}
	} else {
		// Decelerate smoothly to zero
		if p.OnGround {
			p.VX *= settings.PlayerGroundFriction
		} else {
			p.VX *= settings.PlayerAirResistance
		}
		if p.VX > -0.1 && p.VX < 0.1 {
			p.VX = 0
		}
	}
	// Clamp max speed
	if p.VX > maxSpeed {
		p.VX = maxSpeed
	} else if p.VX < -maxSpeed {
		p.VX = -maxSpeed
	}
	// Dampen micro-movements
	if !isMoving {
		entity.DampVelocity(&p.VX, &p.VY, 0.05)
	}
}

// HandleJump processes jump input using entity input state tracking
func (p *Player) HandleJump() {
	// --- Jump buffering and coyote time ---
	// Allow jump for a few frames after leaving ground (coyote time)
	coyoteFrames := 6
	jumpBufferFrames := 6
	if p.OnGround {
		p.InputState.LastGroundedTime = 0
	} else {
		p.InputState.LastGroundedTime++
	}
	// Buffer jump input for a few frames
	if p.InputState.JumpPressed {
		p.InputState.WasJumpPressed = true
		p.InputState.LastJumpPressed = 0
	} else if p.InputState.WasJumpPressed {
		p.InputState.LastJumpPressed++
		if p.InputState.LastJumpPressed > jumpBufferFrames {
			p.InputState.WasJumpPressed = false
		}
	}
	// Allow jump if within coyote time and jump was buffered
	if (p.OnGround || p.InputState.LastGroundedTime <= coyoteFrames) && p.InputState.WasJumpPressed {
		if p.AABB.Jump(settings.PlayerJumpSpeed) {
			p.InputState.WasJumpPressed = false
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
