package noise

import (
	"math"
)

// SimplexNoise implements a simplified version of simplex noise for terrain generation
type SimplexNoise struct {
	perm []int
	seed int64
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

// Linear Congruential Generator for reproducible randomness
type lcg struct {
	state int64
}

func newLCG(seed int64) *lcg {
	return &lcg{state: seed}
}

func (l *lcg) next() int64 {
	l.state = (l.state*1664525 + 1013904223) & 0x7FFFFFFF
	return l.state
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

// Helper functions
func gradient1D(hash int, x float64) float64 {
	// Simple gradient function for 1D
	return float64((hash&1)*2-1) * x
}

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

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// Optimized terrain generation functions for low-end hardware

// FastTerrainNoise generates simple but fast terrain
func (sn *SimplexNoise) FastTerrainNoise(x float64) float64 {
	// Only use 2 octaves for speed
	base := sn.Noise1D(x * 0.01)
	detail := sn.Noise1D(x * 0.03) * 0.5
	return base + detail
}

// FastCaveNoise generates simple cave patterns
func (sn *SimplexNoise) FastCaveNoise(x, y float64) float64 {
	// Single octave for maximum performance
	return sn.Noise2D(x * 0.05, y * 0.08)
}

// FastOreNoise generates ore patterns with minimal computation
func (sn *SimplexNoise) FastOreNoise(x, y float64) float64 {
	// Very simple ore generation
	return sn.Noise2D(x * 0.1, y * 0.1)
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

// MinecraftTerrainNoise generates Minecraft-like terrain with smoother hills
func (sn *SimplexNoise) MinecraftTerrainNoise(x float64) float64 {
	// Minecraft-style terrain with gentler slopes
	continental := sn.Noise1D(x * 0.005) * 0.8 // Continental shape
	hills := sn.Noise1D(x * 0.02) * 0.3         // Rolling hills
	details := sn.Noise1D(x * 0.08) * 0.1       // Fine details

	return continental + hills + details
}

// HybridTerrainNoise combines Minecraft smoothness with Terraria variety
func (sn *SimplexNoise) HybridTerrainNoise(x float64) float64 {
	// Get both terrain types
	minecraft := sn.MinecraftTerrainNoise(x)
	terraria := sn.TerrariaTerrainNoise(x)

	// Blend them based on position for variety
	blendFactor := (sn.Noise1D(x * 0.003) + 1.0) * 0.5 // 0 to 1

	// More Minecraft-like in some areas, more Terraria-like in others
	return minecraft*(1.0-blendFactor) + terraria*blendFactor
}

// TerrariaTerrainNoise generates Terraria-like surface terrain with hills and valleys
func (sn *SimplexNoise) TerrariaTerrainNoise(x float64) float64 {
	// Base terrain using multiple octaves like Terraria
	largeTerrain := sn.Noise1D(x * 0.008)  // Large landmasses
	mediumTerrain := sn.Noise1D(x * 0.02) * 0.5 // Hills
	smallTerrain := sn.Noise1D(x * 0.05) * 0.25 // Small details

	return largeTerrain + mediumTerrain + smallTerrain
}

// HybridCaveNoise creates caves that blend Minecraft and Terraria styles
func (sn *SimplexNoise) HybridCaveNoise(x, y float64) float64 {
	// Minecraft-style caves - more horizontal tunnels
	minecraft := sn.Noise2D(x * 0.04, y * 0.02)

	// Terraria-style caves - more varied
	terraria := sn.TerrariaCaveNoise(x, y)

	// Combine them
	return (minecraft + terraria) * 0.5
}

// TerrariaCaveNoise generates cave patterns similar to Terraria
func (sn *SimplexNoise) TerrariaCaveNoise(x, y float64) float64 {
	// Primary cave tunnels - horizontal bias
	primaryCaves := sn.Noise2D(x * 0.03, y * 0.015)

	// Secondary cave systems - more chaotic
	secondaryCaves := sn.Noise2D(x * 0.06, y * 0.04) * 0.7

	// Large caverns - rare but spacious
	largeCaverns := sn.Noise2D(x * 0.01, y * 0.008) * 1.2

	// Combine cave systems
	return primaryCaves + secondaryCaves + largeCaverns*0.6
}

// TerrariaOreNoise generates ore vein patterns like Terraria
func (sn *SimplexNoise) TerrariaOreNoise(x, y float64, oreType int) float64 {
	// Different scales for different ore types
	var scale, threshold float64

	switch oreType {
	case 0: // Copper - common, small veins
		scale = 0.15
		threshold = 0.7
	case 1: // Iron - medium rarity, medium veins
		scale = 0.12
		threshold = 0.75
	case 2: // Silver - less common, smaller veins
		scale = 0.1
		threshold = 0.8
	case 3: // Gold - rare, small but rich veins
		scale = 0.08
		threshold = 0.85
	case 4: // Platinum - very rare, tiny veins
		scale = 0.06
		threshold = 0.9
	default:
		scale = 0.1
		threshold = 0.8
	}

	// Generate ore vein pattern
	oreNoise := sn.FractalNoise2D(x, y, 3, scale, 1.0, 0.5)

	// Add some randomness to vein shapes
	veinShape := sn.Noise2D(x*scale*2, y*scale*2) * 0.3

	return oreNoise + veinShape - threshold
}

// TerrariaUndergroundNoise generates underground layer transitions
func (sn *SimplexNoise) TerrariaUndergroundNoise(x, y float64) float64 {
	// Dirt to stone transition - should be somewhat uneven
	dirtStoneTransition := sn.FractalNoise2D(x, y, 3, 0.03, 1.0, 0.6)

	// Stone layer variations
	stoneVariation := sn.FractalNoise2D(x*1.2, y*0.8, 2, 0.02, 0.8, 0.5)

	return dirtStoneTransition + stoneVariation*0.5
}

// TerrariaUnderworldNoise generates hell/underworld terrain like Terraria
func (sn *SimplexNoise) TerrariaUnderworldNoise(x, y float64) float64 {
	// Jagged, chaotic terrain for underworld
	ashPockets := sn.FractalNoise2D(x, y, 4, 0.08, 1.0, 0.7)
	lavaPockets := sn.FractalNoise2D(x*0.7, y*1.3, 3, 0.05, 1.2, 0.6)
	hellstoneVeins := sn.RidgedNoise2D(x*2, y*0.5, 2, 0.1, 1.0)

	return ashPockets + lavaPockets*0.8 + hellstoneVeins*0.6
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

// TerrariaBiomeNoise generates biome boundaries similar to Terraria
func (sn *SimplexNoise) TerrariaBiomeNoise(x float64) float64 {
	// Large-scale biome distribution
	largeBiomes := sn.FractalNoise1D(x, 2, 0.002, 1.0, 0.5)

	// Medium-scale biome variations
	mediumBiomes := sn.FractalNoise1D(x*1.3, 3, 0.008, 0.6, 0.6)

	// Small transition zones
	transitions := sn.Noise1D(x*0.01) * 0.3

	return largeBiomes + mediumBiomes + transitions
}
