package noise

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"time"

	"github.com/aquilax/go-perlin"

	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// PerlinNoise implements a wrapper around the aquilax/go-perlin package for Perlin noise generation
type PerlinNoise struct {
	p    *perlin.Perlin
	seed int64
}

// NewRandomPerlinNoise creates a new PerlinNoise generator with a cryptographically random seed
func NewRandomPerlinNoise() *PerlinNoise {
	seed := generateCryptoRandomSeed()
	return NewPerlinNoise(seed)
}

// generateCryptoRandomSeed creates a cryptographically secure random seed
func generateCryptoRandomSeed() int64 {
	var seedBytes [8]byte
	_, err := rand.Read(seedBytes[:])
	if err != nil {
		// fallback to time-based seed if crypto fails
		return int64(binary.LittleEndian.Uint64(seedBytes[:])) ^ time.Now().UnixNano()
	}
	return int64(binary.LittleEndian.Uint64(seedBytes[:]))
}

// NewPerlinNoise creates a new PerlinNoise generator with the given seed
func NewPerlinNoise(seed int64) *PerlinNoise {
	if seed == 0 {
		seed = generateCryptoRandomSeed()
	}
	return &PerlinNoise{
		p:    perlin.NewPerlin(settings.PerlinAlpha, settings.PerlinBeta, settings.PerlinOctaves, seed),
		seed: seed,
	}
}

// Noise1D generates 1D noise value between -1 and 1
func (pn *PerlinNoise) Noise1D(x float64) float64 {
	return pn.p.Noise1D(x)
}

// Noise2D generates 2D noise value between -1 and 1
func (pn *PerlinNoise) Noise2D(x, y float64) float64 {
	return pn.p.Noise2D(x, y)
}

// FractalNoise1D generates 1D fractal noise by combining multiple octaves
func (pn *PerlinNoise) FractalNoise1D(x float64, octaves int, frequency, amplitude, persistence float64) float64 {
	value := 0.0
	maxValue := 0.0

	for i := 0; i < octaves; i++ {
		value += pn.Noise1D(x*frequency) * amplitude
		maxValue += amplitude

		frequency *= 2.0
		amplitude *= persistence
	}

	return value / maxValue
}

// FractalNoise2D generates 2D fractal noise - optimized for performance
func (pn *PerlinNoise) FractalNoise2D(x, y float64, octaves int, frequency, amplitude, persistence float64) float64 {
	value := 0.0
	maxValue := 0.0

	// Limit octaves for performance on low-end hardware
	if octaves > 3 {
		octaves = 3
	}

	for i := 0; i < octaves; i++ {
		value += pn.Noise2D(x*frequency, y*frequency) * amplitude
		maxValue += amplitude

		frequency *= 2.0
		amplitude *= persistence
	}

	return value / maxValue
}

// RidgedNoise1D creates ridged noise (useful for mountain ridges)
func (pn *PerlinNoise) RidgedNoise1D(x float64, octaves int, frequency, amplitude float64) float64 {
	value := 0.0
	maxValue := 0.0

	// Limit octaves for performance
	if octaves > 2 {
		octaves = 2
	}

	for i := 0; i < octaves; i++ {
		n := math.Abs(pn.Noise1D(x * frequency))
		n = 1.0 - n // Invert to create ridges
		n = n * n   // Square to sharpen ridges

		value += n * amplitude
		maxValue += amplitude

		frequency *= 2.0
		amplitude *= 0.5
	}

	return (value/maxValue)*2.0 - 1.0
}

// RidgedNoise2D creates 2D ridged noise for mountain ridges and hellstone veins
func (pn *PerlinNoise) RidgedNoise2D(x, y float64, octaves int, frequency, amplitude float64) float64 {
	value := 0.0
	maxValue := 0.0

	for i := 0; i < octaves; i++ {
		n := math.Abs(pn.Noise2D(x*frequency, y*frequency))
		n = 1.0 - n // Invert to create ridges
		n = n * n   // Square to sharpen ridges

		value += n * amplitude
		maxValue += amplitude

		frequency *= 2.0
		amplitude *= 0.5
	}

	return (value/maxValue)*2.0 - 1.0
}
