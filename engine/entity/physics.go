package entity

// Physics utility functions for entity movement and collision

func ApplyGravity(vy *float64, gravity float64) {
	*vy += gravity
}

func ApplyFriction(vx *float64, friction float64) {
	*vx *= friction
	// Zero out tiny movements to prevent floating point jitter
	if *vx > -0.1 && *vx < 0.1 {
		*vx = 0
	}
}

func ClampVelocity(vy *float64, maxFallSpeed float64) {
	if *vy > maxFallSpeed {
		*vy = maxFallSpeed
	}
}

func ClampPosition(x *float64, min, max float64) {
	if *x < min {
		*x = min
	}
	if *x > max {
		*x = max
	}
}
