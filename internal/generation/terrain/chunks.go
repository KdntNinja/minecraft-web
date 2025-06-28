package terrain

import (
	"fmt"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation/noise"
)

var GetWorldSeedFunc func() int64

// GenerateVisibleTerrain ensures chunks have some visible terrain for testing
func GenerateVisibleTerrain(chunkX, chunkY int) block.Chunk {
	var chunk block.Chunk

	// Use Perlin noise for visible terrain
	terrainNoise := noise.NewPerlinNoise(settings.DefaultSeed)

	// Calculate the global Y start for this chunk
	chunkStartY := chunkY * settings.ChunkHeight

	for y := 0; y < settings.ChunkHeight; y++ {
		for x := 0; x < settings.ChunkWidth; x++ {
			globalX := chunkX*settings.ChunkWidth + x
			globalY := chunkStartY + y

			// --- FIX: Clamp surface height to always be inside the chunk ---
			surfaceNoise := terrainNoise.Noise1D(float64(globalX) * 0.02)
			surfaceHeight := chunkStartY + settings.ChunkHeight/2 + int(surfaceNoise*float64(settings.SurfaceHeightVar))
			if surfaceHeight < chunkStartY+2 {
				surfaceHeight = chunkStartY + 2
			}
			if surfaceHeight > chunkStartY+settings.ChunkHeight-3 {
				surfaceHeight = chunkStartY + settings.ChunkHeight - 3
			}

			// --- BOUNDS CHECK: Only write to chunk[y][x] if y and x are valid ---
			if y < 0 || y >= settings.ChunkHeight || x < 0 || x >= settings.ChunkWidth {
				continue
			}

			if globalY == surfaceHeight {
				biome := GetBiomeAt(terrainNoise, float64(globalX))
				chunk[y][x] = getSurfaceBlockByID(biome)
			} else if globalY > surfaceHeight { // Above ground
				biomeData := terrainNoise.GetBiomeAt(float64(globalX))
				if globalY == surfaceHeight+1 && shouldPlaceTreeByID(globalX, int(biomeData.Type)) {
					chunk[y][x] = block.Wood
				} else if globalY == surfaceHeight+2 && shouldPlaceTreeByID(globalX, int(biomeData.Type)) {
					chunk[y][x] = block.Leaves
				} else {
					chunk[y][x] = block.Air
				}
			} else {
				// Underground
				depth := surfaceHeight - globalY
				if depth <= 3 {
					chunk[y][x] = block.Dirt
				} else if depth <= 10 {
					stoneNoise := terrainNoise.Noise2D(float64(globalX)*0.1, float64(globalY)*0.1)
					if stoneNoise > 0.3 {
						chunk[y][x] = block.Stone
					} else {
						chunk[y][x] = block.Dirt
					}
				} else {
					chunk[y][x] = block.Stone
				}
			}
		}
	}

	return chunk
}

// GenerateChunk creates a chunk with authentic Terraria-like terrain layers
func GenerateChunk(chunkX, chunkY int) block.Chunk {
	fmt.Printf("CHUNK_GEN: Generating Terraria-like chunk at (%d, %d)\n", chunkX, chunkY)
	var chunk block.Chunk

	chunkStartY := chunkY * settings.ChunkHeight
	var seed int64 = settings.DefaultSeed
	if GetWorldSeedFunc != nil {
		seed = GetWorldSeedFunc()
	}
	terrainNoise := noise.NewPerlinNoise(seed)

	// World layer thresholds (Terraria-style) - adjusted for proper coordinate system
	baseSurfaceHeight := 60.0 // Average surface height in global coordinates
	caveBeltStart := 80       // Where caves become more common
	deepLayer := 150          // Deep underground layer
	underworldStart := 200    // Underworld begins
	worldBottom := 240        // Bottom of the world

	// Debug info for the first chunk
	if chunkX == 0 && chunkY == 0 {
		fmt.Printf("CHUNK_DEBUG: Chunk (0,0) - chunkStartY=%d, baseSurfaceHeight=%.1f\n", chunkStartY, baseSurfaceHeight)
	}

	for y := 0; y < settings.ChunkHeight; y++ {
		globalY := chunkStartY + y
		for x := 0; x < settings.ChunkWidth; x++ {
			globalX := chunkX*settings.ChunkWidth + x

			// Get biome data for this position
			biome := terrainNoise.GetBiomeAt(float64(globalX))

			// Calculate dynamic surface height with extreme variation
			biomeHeight := terrainNoise.GetBiomeTerrainHeight(float64(globalX), biome)
			surfaceNoise := terrainNoise.FractalNoise1D(float64(globalX)*0.012, 3, 0.015, 1.0, 0.6)
			detailNoise := terrainNoise.FractalNoise1D(float64(globalX)*0.045, 2, 0.04, 0.5, 0.5)

			// Calculate the actual surface height for this X coordinate
			surfaceHeight := int(baseSurfaceHeight + biomeHeight*25.0 + surfaceNoise*15.0 + detailNoise*8.0)

			// Clamp surface height to reasonable bounds
			if surfaceHeight < 10 {
				surfaceHeight = 10
			}
			if surfaceHeight > 120 {
				surfaceHeight = 120
			}

			// Debug surface height for specific positions
			if chunkX == 0 && chunkY == 0 && x == 0 && y == 0 {
				fmt.Printf("SURFACE_DEBUG: globalX=%d, globalY=%d, biomeHeight=%.2f, surfaceHeight=%d\n",
					globalX, globalY, biomeHeight, surfaceHeight)
			}

			// Calculate depth from surface for this position
			depthFromSurface := surfaceHeight - globalY

			// --- BOUNDS CHECK: Prevent index out of range ---
			if y < 0 || y >= settings.ChunkHeight {
				continue
			}

			// Determine block based on depth and world layers
			if globalY > surfaceHeight {
				// Above surface - should be air
				chunk[y][x] = block.Air
			} else if globalY == surfaceHeight {
				// Surface layer - biome-specific blocks
				chunk[y][x] = getTerrariaStyleSurfaceBlock(biome, terrainNoise, globalX)

				// Add trees based on biome (improved tree generation)
				if shouldPlaceTree(globalX, biome) {
					generateTree(&chunk, x, y, biome, terrainNoise, globalX)
				}
			} else if depthFromSurface <= 3 {
				// Shallow soil layer (dirt/sand based on biome)
				chunk[y][x] = getTerrariaStyleSoilBlock(biome, terrainNoise, globalX, globalY)
			} else if depthFromSurface <= 15 {
				// Dirt-to-stone transition layer with variety
				transitionNoise := terrainNoise.FractalNoise2D(float64(globalX)*0.08, float64(globalY)*0.08, 2, 0.1, 1.0, 0.5)
				if transitionNoise > 0.3 {
					chunk[y][x] = block.Stone
				} else if transitionNoise > -0.2 {
					chunk[y][x] = block.Dirt
				} else {
					chunk[y][x] = getTerrariaStyleSoilBlock(biome, terrainNoise, globalX, globalY)
				}
			} else if globalY < caveBeltStart || depthFromSurface < 20 {
				// Early stone layer with minimal caves
				if shouldGenerateTerrariaStyleCave(globalX, globalY, depthFromSurface, terrainNoise, 0.7) {
					chunk[y][x] = block.Air
				} else {
					// Check for shallow ores
					if ore := generateOre(globalX, globalY, depthFromSurface, biome, terrainNoise); ore != block.Air {
						chunk[y][x] = ore
					} else {
						chunk[y][x] = block.Stone
					}
				}
			} else if globalY < deepLayer {
				// Cave belt - lots of caves and tunnels (like Terraria's main cave layer)
				if shouldGenerateTerrariaStyleCave(globalX, globalY, depthFromSurface, terrainNoise, 0.5) {
					chunk[y][x] = block.Air
				} else {
					// Stone with good ore distribution
					if ore := generateOre(globalX, globalY, depthFromSurface, biome, terrainNoise); ore != block.Air {
						chunk[y][x] = ore
					} else {
						chunk[y][x] = block.Stone
					}
				}
			} else if globalY < underworldStart {
				// Deep underground - denser stone, rare caves, precious ores
				if shouldGenerateTerrariaStyleCave(globalX, globalY, depthFromSurface, terrainNoise, 0.8) {
					chunk[y][x] = block.Air
				} else {
					// Check for deep ores first
					if ore := generateOre(globalX, globalY, depthFromSurface, biome, terrainNoise); ore != block.Air {
						chunk[y][x] = ore
					} else {
						// Mostly stone with some variety
						stoneVariation := terrainNoise.FractalNoise2D(float64(globalX)*0.03, float64(globalY)*0.03, 2, 0.04, 1.0, 0.5)
						if stoneVariation > 0.6 {
							chunk[y][x] = block.Clay // Clay pockets
						} else {
							chunk[y][x] = block.Stone
						}
					}
				}
			} else if globalY < worldBottom-8 {
				// Underworld layer - ash with hellstone veins
				underworldNoise := terrainNoise.FractalNoise2D(float64(globalX)*0.06, float64(globalY)*0.06, 3, 0.08, 1.0, 0.6)
				hellstoneVeins := terrainNoise.RidgedNoise2D(float64(globalX)*0.15, float64(globalY)*0.12, 2, 0.1, 1.0)

				if hellstoneVeins > 0.4 {
					chunk[y][x] = block.Hellstone
				} else if underworldNoise > 0.3 {
					chunk[y][x] = block.Ash
				} else {
					chunk[y][x] = block.Stone
				}
			} else {
				// Bottom bedrock layer - solid hellstone
				chunk[y][x] = block.Hellstone
			}
		}
	}

	return chunk
}
