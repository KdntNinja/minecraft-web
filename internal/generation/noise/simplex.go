package noise

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"math"
	"time"
)

// SimplexNoise implements a simplified version of simplex noise for terrain generation
type SimplexNoise struct {
	perm []int
	seed int64
}

// EnhancedNoiseGenerator provides multiple noise algorithms for better terrain generation
type EnhancedNoiseGenerator struct {
	Primary   *SimplexNoise
	Secondary *SimplexNoise
	Detail    *SimplexNoise
	Cave      *SimplexNoise
	Ore       *SimplexNoise
	Biome     *SimplexNoise
	Structure *SimplexNoise
	Seed      int64
}

// NewEnhancedNoiseGenerator creates a new enhanced noise generator with crypto-random seed
func NewEnhancedNoiseGenerator() *EnhancedNoiseGenerator {
	seed := generateCryptoRandomSeed()
	log.Printf("Generating enhanced noise with seed: %d", seed)
	return NewEnhancedNoiseGeneratorWithSeed(seed)
}

// NewEnhancedNoiseGeneratorWithSeed creates enhanced noise with specific seed
func NewEnhancedNoiseGeneratorWithSeed(seed int64) *EnhancedNoiseGenerator {
	return &EnhancedNoiseGenerator{
		Primary:   NewSimplexNoise(seed),
		Secondary: NewSimplexNoise(seed + 1000),
		Detail:    NewSimplexNoise(seed + 2000),
		Cave:      NewSimplexNoise(seed + 3000),
		Ore:       NewSimplexNoise(seed + 4000),
		Biome:     NewSimplexNoise(seed + 5000),
		Structure: NewSimplexNoise(seed + 6000),
		Seed:      seed,
	}
}

// generateCryptoRandomSeed creates a cryptographically secure random seed
func generateCryptoRandomSeed() int64 {
	var seedBytes [8]byte
	_, err := rand.Read(seedBytes[:])
	if err != nil {
		// Fallback to time-based seed if crypto/rand fails
		log.Printf("Warning: Failed to generate crypto random seed, using time-based seed: %v", err)
		return time.Now().UnixNano()
	}

	seed := int64(binary.LittleEndian.Uint64(seedBytes[:]))
	// Ensure positive seed
	if seed < 0 {
		seed = -seed
	}
	return seed
}

// NewSimplexNoise creates a new SimplexNoise generator with the given seed
func NewSimplexNoise(seed int64) *SimplexNoise {
	sn := &SimplexNoise{
		seed: seed,
		perm: make([]int, 512),
	}

	// Initialize permutation table based on seed
	for i := 0; i < 256; i++ {
		sn.perm[i] = i
	}

	// Shuffle the permutation table using the seed
	rng := newLCG(seed)
	for i := 255; i > 0; i-- {
		j := int(rng.next() % int64(i+1))
		sn.perm[i], sn.perm[j] = sn.perm[j], sn.perm[i]
	}

	// Duplicate the permutation table
	for i := 0; i < 256; i++ {
		sn.perm[256+i] = sn.perm[i]
	}

	return sn
}

// Noise1D generates 1D noise value between -1 and 1
func (sn *SimplexNoise) Noise1D(x float64) float64 {
	// Scale the input
	x *= 0.5

	// Get the integer part
	i := int(math.Floor(x))

	// Get the fractional part
	f := x - float64(i)

	// Smooth the fractional part using smoothstep function
	u := f * f * (3.0 - 2.0*f)

	// Get permutation indices
	a := sn.perm[i&255]
	b := sn.perm[(i+1)&255]

	// Generate gradient values
	ga := gradient1D(a, f)
	gb := gradient1D(b, f-1.0)

	// Interpolate between the two gradient values
	return lerp(ga, gb, u)
}

// Noise2D generates 2D noise value between -1 and 1
func (sn *SimplexNoise) Noise2D(x, y float64) float64 {
	// Scale the input
	x *= 0.5
	y *= 0.5

	// Get the integer parts
	ix := int(math.Floor(x))
	iy := int(math.Floor(y))

	// Get the fractional parts
	fx := x - float64(ix)
	fy := y - float64(iy)

	// Smooth the fractional parts
	ux := fx * fx * (3.0 - 2.0*fx)
	uy := fy * fy * (3.0 - 2.0*fy)

	// Get permutation indices for the four corners
	a := sn.perm[ix&255] + iy
	b := sn.perm[(ix+1)&255] + iy

	aa := sn.perm[a&255]
	ab := sn.perm[(a+1)&255]
	ba := sn.perm[b&255]
	bb := sn.perm[(b+1)&255]

	// Generate gradient values for the four corners
	g1 := gradient2D(aa, fx, fy)
	g2 := gradient2D(ba, fx-1.0, fy)
	g3 := gradient2D(ab, fx, fy-1.0)
	g4 := gradient2D(bb, fx-1.0, fy-1.0)

	// Interpolate
	i1 := lerp(g1, g2, ux)
	i2 := lerp(g3, g4, ux)

	return lerp(i1, i2, uy)
}

// FractalNoise1D generates 1D fractal noise by combining multiple octaves
func (sn *SimplexNoise) FractalNoise1D(x float64, octaves int, frequency, amplitude, persistence float64) float64 {
	value := 0.0
	maxValue := 0.0

	for i := 0; i < octaves; i++ {
		value += sn.Noise1D(x*frequency) * amplitude
		maxValue += amplitude

		frequency *= 2.0
		amplitude *= persistence
	}

	return value / maxValue
}

// FractalNoise2D generates 2D fractal noise - optimized for performance
func (sn *SimplexNoise) FractalNoise2D(x, y float64, octaves int, frequency, amplitude, persistence float64) float64 {
	value := 0.0
	maxValue := 0.0

	// Limit octaves for performance on low-end hardware
	if octaves > 3 {
		octaves = 3
	}

	for i := 0; i < octaves; i++ {
		value += sn.Noise2D(x*frequency, y*frequency) * amplitude
		maxValue += amplitude

		frequency *= 2.0
		amplitude *= persistence
	}

	return value / maxValue
}

// RidgedNoise1D creates ridged noise (useful for mountain ridges)
func (sn *SimplexNoise) RidgedNoise1D(x float64, octaves int, frequency, amplitude float64) float64 {
	value := 0.0
	maxValue := 0.0

	// Limit octaves for performance
	if octaves > 2 {
		octaves = 2
	}

	for i := 0; i < octaves; i++ {
		n := math.Abs(sn.Noise1D(x * frequency))
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
func (sn *SimplexNoise) RidgedNoise2D(x, y float64, octaves int, frequency, amplitude float64) float64 {
	value := 0.0
	maxValue := 0.0

	for i := 0; i < octaves; i++ {
		n := math.Abs(sn.Noise2D(x*frequency, y*frequency))
		n = 1.0 - n // Invert to create ridges
		n = n * n   // Square to sharpen ridges

		value += n * amplitude
		maxValue += amplitude

		frequency *= 2.0
		amplitude *= 0.5
	}

	return (value/maxValue)*2.0 - 1.0
}
