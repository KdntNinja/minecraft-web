package player

// ApplyMovement handles horizontal movement physics with ground/air control differences
func (p *Player) ApplyMovement(isMoving bool, targetVX float64) {
	if isMoving {
		if p.OnGround {
			p.VX = targetVX // Instant ground acceleration
		} else {
			p.VX = p.VX*0.8 + targetVX*0.2 // Reduced air control (80% momentum + 20% input)
		}
	} else {
		if p.OnGround {
			p.VX *= GroundFriction // Apply ground friction
		} else {
			p.VX *= AirResistance // Minimal air resistance
		}

		// Stop micro-movements to prevent jitter
		if p.VX > -0.1 && p.VX < 0.1 {
			p.VX = 0
		}
	}
}

// HandleJump processes jump input with key state tracking to prevent spam
func (p *Player) HandleJump(jumpKeyPressed bool) {
	// Only jump on fresh key press while grounded
	if jumpKeyPressed && !p.jumpPressed && p.OnGround {
		p.VY = JumpSpeed
		p.OnGround = false
	}

	p.jumpPressed = jumpKeyPressed // Track key state for next frame
}

// ApplyGravity handles gravity, fall speed limiting, and ground sticking
func (p *Player) ApplyGravity() {
	// Update grounded time tracking
	if p.OnGround {
		p.lastGroundedTime = 0
	} else {
		p.lastGroundedTime++
	}

	if !p.OnGround {
		p.VY += Gravity // Apply gravity acceleration

		if p.VY > MaxFallSpeed {
			p.VY = MaxFallSpeed // Cap fall speed
		}
	} else {
		if p.VY > 0 {
			p.VY = 0 // Stick to ground when landed
		}
	}
}
