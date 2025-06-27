package noise

// Helper functions for noise generation

// gradient1D generates a simple gradient function for 1D noise
func gradient1D(hash int, x float64) float64 {
	// Simple gradient function for 1D
	return float64((hash&1)*2-1) * x
}

// gradient2D generates a simple gradient function for 2D noise using predefined vectors
func gradient2D(hash int, x, y float64) float64 {
	// Simple gradient function for 2D using predefined vectors
	h := hash & 3
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

// lerp performs linear interpolation between two values
func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// Linear Congruential Generator for reproducible randomness
type lcg struct {
	state int64
}

// newLCG creates a new Linear Congruential Generator with the given seed
func newLCG(seed int64) *lcg {
	return &lcg{state: seed}
}

// next generates the next random number in the sequence
func (l *lcg) next() int64 {
	l.state = (l.state*1664525 + 1013904223) & 0x7FFFFFFF
	return l.state
}
