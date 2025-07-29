package physics

// Physics utility functions for entity movement and collision

func ApplyGravity(vy *float64, gravity float64) {
	*vy += gravity
}

func ApplyFriction(vx *float64, friction float64) {
	*vx *= friction
	// Zero out tiny movements to prevent floating point jitter (tighter threshold for quicker settling)
	if *vx > -0.05 && *vx < 0.05 {
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

// DampVelocity reduces velocity to prevent micro-movements and jitter
func DampVelocity(vx, vy *float64, threshold float64) {
	// Stop very small horizontal movements
	if *vx > -threshold && *vx < threshold {
		*vx = 0
	}
	// Stop very small vertical movements
	if *vy > -threshold && *vy < threshold {
		*vy = 0
	}
}

// StabilizePosition snaps position to prevent floating point drift
func StabilizePosition(pos *float64, gridSize int) {
	// Round to nearest 0.1 pixel to prevent drift
	rounded := float64(int(*pos*10)) / 10
	if *pos-rounded < 0.05 && *pos-rounded > -0.05 {
		*pos = rounded
	}
}
