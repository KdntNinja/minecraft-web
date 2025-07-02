package player

import (
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// ApplyMovement handles the player's horizontal movement, applying acceleration
// and friction to create a responsive, Minecraft-like feel. It distinguishes
// between ground and air movement, and handles sprinting and sneaking.
func (p *Player) ApplyMovement(isMoving bool, targetVX float64) {
	// --- Minecraft-like Movement Physics ---

	// Retrieve movement parameters from the settings package. This allows for
	// easy tuning of the player's movement characteristics.
	walkAccel := settings.PlayerWalkAccel
	airAccel := settings.PlayerAirAccel
	groundFriction := settings.PlayerGroundFriction
	airFriction := settings.PlayerAirFriction
	sneakSpeed := settings.PlayerSneakSpeed
	sneakAccel := settings.PlayerSneakAccel
	maxWalkSpeed := settings.PlayerMoveSpeed
	maxSprintSpeed := settings.PlayerSprintSpeed
	maxSneakSpeed := sneakSpeed
	maxSpeed := maxWalkSpeed
	accel := walkAccel

	// Adjust speed and acceleration for sneaking (holding Shift) or sprinting.
	if p.InputState.SneakPressed {
		maxSpeed = maxSneakSpeed
		accel = sneakAccel
	} else if p.IsSprinting {
		maxSpeed = maxSprintSpeed
		accel = settings.PlayerSprintAccel
	}

	// Apply acceleration or friction based on whether the player is trying to move.
	if isMoving {
		// When on the ground, acceleration is higher for a snappier feel.
		if p.OnGround {
			if targetVX > 0 {
				p.VX += accel
				if p.VX > maxSpeed {
					p.VX = maxSpeed
				}
			} else if targetVX < 0 {
				p.VX -= accel
				if p.VX < -maxSpeed {
					p.VX = -maxSpeed
				}
			}
		} else {
			// In the air, the player has less control (lower acceleration).
			if targetVX > 0 {
				p.VX += airAccel
				if p.VX > maxSpeed {
					p.VX = maxSpeed
				}
			} else if targetVX < 0 {
				p.VX -= airAccel
				if p.VX < -maxSpeed {
					p.VX = -maxSpeed
				}
			}
		}
	} else {
		// When no movement input is given, apply friction to slow the player down.
		if p.OnGround {
			p.VX *= groundFriction
		} else {
			p.VX *= airFriction
		}
		// If velocity is very close to zero, set it to zero to prevent "sliding".
		if p.VX > -0.02 && p.VX < 0.02 {
			p.VX = 0
		}
	}
}

// HandleJump manages the player's jump action. It incorporates several mechanics
// common in platformers to make jumping feel fair and responsive:
// - Coyote Time: Allows jumping for a few frames after leaving a ledge.
// - Jump Buffering: Registers a jump press just before landing, so it executes on touchdown.
// - Variable Jump Height: Allows for shorter hops or full jumps based on how long the jump button is held.
func (p *Player) HandleJump() {
	// --- Advanced Jump Mechanics ---
	coyoteFrames := settings.PlayerCoyoteFrames
	jumpBufferFrames := settings.PlayerJumpBufferFrames
	jumpHoldMax := settings.PlayerJumpHoldMax

	// When on the ground, reset coyote time and jump hold duration.
	// Otherwise, increment the time since the player was last grounded.
	if p.OnGround {
		p.InputState.LastGroundedTime = 0
		p.InputState.JumpHoldTime = 0
	} else {
		p.InputState.LastGroundedTime++
	}

	// Buffer the jump input. If the jump button is pressed, we register it and
	// keep it active for a few frames, even if the button is released.
	if p.InputState.JumpPressed {
		p.InputState.WasJumpPressed = true
		p.InputState.LastJumpPressed = 0
	} else if p.InputState.WasJumpPressed {
		p.InputState.LastJumpPressed++
		if p.InputState.LastJumpPressed > jumpBufferFrames {
			p.InputState.WasJumpPressed = false
		}
	}

	// A jump is initiated if the player is on the ground OR within the coyote time
	// window, and a jump has been buffered.
	if (p.OnGround || p.InputState.LastGroundedTime <= coyoteFrames) && p.InputState.WasJumpPressed {
		if p.AABB.Jump(settings.PlayerJumpSpeed * 1.15) { // Higher initial jump velocity for a better feel
			p.InputState.WasJumpPressed = false
			p.InputState.JumpHoldTime = 1
		}
	}

	// Implement variable jump height. As long as the jump button is held
	// (up to a maximum duration), a small upward force is continuously applied.
	if !p.OnGround && p.InputState.JumpPressed && p.InputState.JumpHoldTime > 0 && p.InputState.JumpHoldTime < jumpHoldMax {
		p.VY -= settings.PlayerJumpHoldForce // Apply upward force to extend the jump
		p.InputState.JumpHoldTime++
	}

	// If the jump button is released, stop applying the upward force.
	if !p.InputState.JumpPressed {
		p.InputState.JumpHoldTime = 0
	}
}

// ApplyGravity handles the vertical physics for the player, including gravity
// and terminal velocity. It also ensures the player position is stabilized when
// they are on the ground.
func (p *Player) ApplyGravity() {
	// Use a slightly stronger gravity when the player is near the apex of their
	// jump to make the arc feel more natural and less "floaty".
	gravityToApply := settings.PlayerGravity
	if p.VY > 0 && p.VY < 2.0 {
		gravityToApply = settings.PlayerGravity * 1.5 // Stronger gravity at jump apex
	}

	// Apply gravity, respecting the maximum fall speed (terminal velocity).
	p.AABB.ApplyVerticalMovement(gravityToApply, settings.PlayerMaxFallSpeed)

	// When the player is on the ground, stabilize their Y position to prevent
	// bouncing or jittering, and apply damping to slow them down smoothly.
	if p.OnGround {
		entity.StabilizePosition(&p.Y, settings.TileSize)
		entity.DampVelocity(&p.VX, &p.VY, 0.2) // Dampen velocity for a smooth landing
	}
}
