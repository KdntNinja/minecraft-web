package noise

// Utility functions for noise generation algorithms

// Helper function to clamp values (inlined for performance)
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func gradient1D(hash int, x float64) float64 {
	return float64((hash&1)*2-1) * x // Convert hash to -1 or 1, multiply by distance
}

func gradient2D(hash int, x, y float64) float64 {
	h := hash & 3 // Use bottom 2 bits for 4 directions
	switch h {
	case 0:
		return x + y
	case 1:
		return -x + y
	case 2:
		return x - y
	default:
		return -x - y
	}
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a) // Linear interpolation between a and b
}

// Linear Congruential Generator for deterministic random numbers
// Used for deterministic noise-based randomness

type lcg struct {
	state int64
}

func newLCG(seed int64) *lcg {
	return &lcg{state: seed}
}

func (l *lcg) next() int64 {
	l.state = (l.state*1664525 + 1013904223) & 0x7FFFFFFF // Standard LCG formula
	return l.state
}
