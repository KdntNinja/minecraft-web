package physics

// Movement utility functions for entities

// ApplyHorizontalMovement handles horizontal movement with ground/air control differences
func (a *AABB) ApplyHorizontalMovement(targetVX, groundFriction, airResistance float64, isMoving bool) {
	if isMoving {
		if a.OnGround {
			a.VX = targetVX // Instant ground acceleration
		} else {
			a.VX = a.VX*0.8 + targetVX*0.2 // Reduced air control (80% momentum + 20% input)
		}
	} else {
		if a.OnGround {
			a.VX *= groundFriction // Apply ground friction
		} else {
			a.VX *= airResistance // Minimal air resistance
		}

		// Stop micro-movements to prevent jitter
		if a.VX > -0.1 && a.VX < 0.1 {
			a.VX = 0
		}
	}
}

// ApplyVerticalMovement handles gravity, fall speed limiting, and ground sticking
func (a *AABB) ApplyVerticalMovement(gravity, maxFallSpeed float64) {
	if !a.OnGround {
		a.VY += gravity // Apply gravity acceleration

		if a.VY > maxFallSpeed {
			a.VY = maxFallSpeed // Cap fall speed
		}
	} else {
		if a.VY > 0 {
			a.VY = 0 // Stick to ground when landed
		}
	}
}

// Jump applies jump velocity if conditions are met
func (a *AABB) Jump(jumpSpeed float64) bool {
	if a.OnGround {
		a.VY = jumpSpeed
		a.OnGround = false
		return true
	}
	return false
}
