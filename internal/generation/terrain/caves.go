package terrain

import (
	"math"

	"github.com/KdntNinja/webcraft/internal/generation/noise"
)

// shouldGenerateTerrariaStyleCave creates Terraria-like cave systems
func shouldGenerateTerrariaStyleCave(globalX, globalY, depthFromSurface int, terrainNoise *noise.PerlinNoise, difficulty float64) bool {
	if depthFromSurface < 8 {
		return false // No caves too close to surface
	}

	x, y := float64(globalX), float64(globalY)

	// Large cavern systems (like Terraria's big open areas)
	largeCaverns := terrainNoise.FractalNoise2D(x*0.012, y*0.015, 3, 0.025, 1.2, 0.6)

	// Winding tunnels (like Terraria's connecting passages)
	tunnels := terrainNoise.FractalNoise2D(x*0.03, y*0.025, 2, 0.04, 0.8, 0.5)
	tunnelWarp := terrainNoise.Noise2D(x*0.008, y*0.01) * 15.0
	warpedTunnels := terrainNoise.Noise2D(x+tunnelWarp, y*0.8) * 0.6

	// Vertical shafts (occasional deep connections)
	verticalShafts := terrainNoise.FractalNoise2D(x*0.005, y*0.08, 2, 0.02, 1.0, 0.4)

	// Depth-based cave probability (more caves deeper down)
	depthFactor := math.Min(float64(depthFromSurface)/40.0, 1.0)

	// Combine cave types
	totalCaveNoise := (largeCaverns*0.8 + tunnels*0.5 + warpedTunnels*0.7 + verticalShafts*0.5) * depthFactor

	// Adjust threshold based on difficulty (higher = fewer caves)
	threshold := 0.45 + (difficulty-0.5)*0.3

	return totalCaveNoise > threshold
}
