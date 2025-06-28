package terrain

import (
	"fmt"
	"math"

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
			} else if globalY < surfaceHeight {
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
			} else {
				// Above ground
				biome := GetBiomeAt(terrainNoise, float64(globalX))
				if globalY == surfaceHeight+1 && shouldPlaceTreeByID(globalX, biome) {
					chunk[y][x] = block.Wood
				} else if globalY == surfaceHeight+2 && shouldPlaceTreeByID(globalX, biome) {
					chunk[y][x] = block.Leaves
				} else {
					chunk[y][x] = block.Air
				}
			}
		}
	}

	return chunk
}

// Helper to map biome ID to surface block
func getSurfaceBlockByID(biome int) block.BlockType {
	switch biome {
	case 0:
		return block.Grass
	case 1:
		return block.Grass
	case 2:
		return block.Sand
	case 3:
		return block.Stone
	case 4:
		return block.Mud
	case 5:
		return block.Snow
	case 6:
		return block.Grass
	case 7:
		return block.Water
	default:
		return block.Grass
	}
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

	// World layer thresholds (Terraria-style)
	surfaceLayer := 60     // Average surface height
	caveBeltStart := 80    // Where caves become more common
	deepLayer := 150       // Deep underground layer
	underworldStart := 200 // Underworld begins
	worldBottom := 240     // Bottom of the world

	for y := 0; y < settings.ChunkHeight; y++ {
		globalY := chunkStartY + y
		for x := 0; x < settings.ChunkWidth; x++ {
			globalX := chunkX*settings.ChunkWidth + x

			// Get biome data for this position
			biome := terrainNoise.GetBiomeAt(float64(globalX))

			// Calculate dynamic surface height with extreme variation
			baseSurface := float64(surfaceLayer)
			biomeHeight := terrainNoise.GetBiomeTerrainHeight(float64(globalX), biome)
			surfaceNoise := terrainNoise.FractalNoise1D(float64(globalX)*0.012, 3, 0.015, 1.0, 0.6)
			detailNoise := terrainNoise.FractalNoise1D(float64(globalX)*0.045, 2, 0.04, 0.5, 0.5)

			surfaceHeight := int(baseSurface + biomeHeight*25.0 + surfaceNoise*15.0 + detailNoise*8.0)

			// Calculate depth from surface for this position
			depthFromSurface := globalY - surfaceHeight

			// --- BOUNDS CHECK: Prevent index out of range ---
			if y < 0 || y >= settings.ChunkHeight {
				continue
			}

			// Determine block based on depth and world layers
			if globalY < surfaceHeight {
				// Sky/Air layer
				chunk[y][x] = block.Air
			} else if globalY == surfaceHeight {
				// Surface layer - biome-specific blocks
				chunk[y][x] = getTerrariaStyleSurfaceBlock(biome, terrainNoise, globalX)

				// Add trees based on biome
				if shouldPlaceTree(globalX, biome) {
					// Add tree trunk above surface (if within chunk bounds)
					if y > 0 {
						chunk[y-1][x] = block.Wood
					}
					if y > 1 {
						chunk[y-2][x] = block.Leaves
					}
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
					if ore := generateTerrariaStyleOre(globalX, globalY, depthFromSurface, biome, terrainNoise); ore != block.Air {
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
					if ore := generateTerrariaStyleOre(globalX, globalY, depthFromSurface, biome, terrainNoise); ore != block.Air {
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
					if ore := generateTerrariaStyleOre(globalX, globalY, depthFromSurface, biome, terrainNoise); ore != block.Air {
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

// getSurfaceHeightRange calculates the min and max surface heights for a chunk column
func getSurfaceHeightRange(chunkX int, terrainNoise *noise.PerlinNoise) (int, int) {
	minHeight := 1000
	maxHeight := -1000

	// Sample a few points across the chunk width to get height range
	for x := 0; x < settings.ChunkWidth; x++ {
		globalX := chunkX*settings.ChunkWidth + x
		biome := terrainNoise.GetBiomeAt(float64(globalX))
		terrainHeight := terrainNoise.GetBiomeTerrainHeight(float64(globalX), biome)

		baseHeight := 50.0
		heightVariation := terrainHeight * 20.0
		surfaceHeight := int(baseHeight + heightVariation)

		if surfaceHeight < minHeight {
			minHeight = surfaceHeight
		}
		if surfaceHeight > maxHeight {
			maxHeight = surfaceHeight
		}
	}

	return minHeight, maxHeight
}

// Enhanced cave generation with realistic patterns
func shouldGenerateCave(globalX, globalY, depthFromSurface int, noise *noise.PerlinNoise) bool {
	if depthFromSurface < 5 {
		return false // No caves too close to surface
	}

	// Multiple cave systems at different scales
	x, y := float64(globalX), float64(globalY)

	// Large cave systems
	largeCaves := noise.FractalNoise2D(x*0.015, y*0.015, 3, 0.02, 1.0, 0.6)

	// Medium caves
	mediumCaves := noise.FractalNoise2D(x*0.03, y*0.025, 2, 0.04, 0.8, 0.5)

	// Cave tunnels using worm-like patterns
	wormX := noise.Noise2D(x*0.005, y*0.008) * 30
	wormY := noise.Noise2D(x*0.008, y*0.005) * 20
	tunnels := noise.Noise2D(x+wormX, y+wormY) * 0.4

	// Depth-based cave probability
	depthFactor := math.Min(float64(depthFromSurface)/30.0, 1.0)

	totalCaveNoise := (largeCaves + mediumCaves*0.7 + tunnels) * depthFactor

	// Dynamic threshold based on depth
	threshold := 0.6 - (depthFactor * 0.1) // Easier caves deeper down

	return totalCaveNoise > threshold
}

// Enhanced ore generation with realistic distribution
func generateOre(globalX, globalY, depthFromSurface int, biome noise.BiomeData, terrainNoise *noise.PerlinNoise) block.BlockType {
	x, y := float64(globalX), float64(globalY)

	// Biome-based ore modifiers
	biomeOreMultiplier := 1.0
	switch biome.Type {
	case 3: // MountainBiome
		biomeOreMultiplier = 1.4 // Mountains have more ores
	case 2: // DesertBiome
		biomeOreMultiplier = 0.7 // Deserts have fewer ores
	case 4: // SwampBiome
		biomeOreMultiplier = 0.8
	}

	// Depth-based ore distribution
	if depthFromSurface >= 5 && depthFromSurface < 20 {
		// Shallow ores: Copper and Iron
		copperNoise := terrainNoise.FractalNoise2D(x*0.08, y*0.09, 3, 0.12, 1.0, 0.6) * biomeOreMultiplier
		ironNoise := terrainNoise.FractalNoise2D(x*0.07+100, y*0.08, 3, 0.1, 1.0, 0.5) * biomeOreMultiplier

		if copperNoise > 0.75 {
			return block.CopperOre
		}
		if ironNoise > 0.8 {
			return block.IronOre
		}
	} else if depthFromSurface >= 20 && depthFromSurface < 60 {
		// Medium depth: Silver and Gold
		silverNoise := terrainNoise.FractalNoise2D(x*0.06+200, y*0.07, 4, 0.08, 1.2, 0.4) * biomeOreMultiplier
		goldNoise := terrainNoise.FractalNoise2D(x*0.05+300, y*0.06, 4, 0.06, 1.3, 0.3) * biomeOreMultiplier

		if goldNoise > 0.88 {
			return block.GoldOre
		}
		if silverNoise > 0.85 {
			return block.IronOre
		}
	} else if depthFromSurface >= 60 {
		// Deep ores: Platinum and rare ores
		platinumNoise := terrainNoise.FractalNoise2D(x*0.04+400, y*0.05, 4, 0.05, 1.4, 0.3) * biomeOreMultiplier

		if platinumNoise > 0.92 {
			return block.GoldOre
		}
	}

	return block.Air // No ore
}

// Get underground block based on depth, biome, and noise
func getUndergroundBlock(depthFromSurface int, biome noise.BiomeData, terrainNoise *noise.PerlinNoise, globalX, globalY int) block.BlockType {
	x, y := float64(globalX), float64(globalY)

	if depthFromSurface <= 3 {
		// Surface soil layer
		switch biome.Type {
		case 2: // DesertBiome
			return block.Sand
		case 5: // TundraBiome
			return block.Snow
		case 4: // SwampBiome
			// Mix of mud and dirt
			mudNoise := terrainNoise.Noise2D(x*0.1, y*0.1)
			if mudNoise > 0.3 {
				return block.Mud
			}
			return block.Dirt
		default:
			return block.Dirt
		}
	} else if depthFromSurface <= 8 {
		// Transition layer with some variation
		transitionNoise := terrainNoise.Noise2D(x*0.05, y*0.05)
		if transitionNoise > 0.4 {
			return block.Dirt
		}
		return block.Stone
	} else if depthFromSurface <= 20 {
		// Mixed stone and dirt layer
		mixNoise := terrainNoise.Noise2D(x*0.03, y*0.03)
		if mixNoise > 0.6 {
			return block.Dirt
		} else if mixNoise > 0.2 {
			return block.Stone
		} else {
			// Some variety in stone types
			switch biome.Type {
			case 2: // DesertBiome
				return block.Sand // Replace Sandstone with Sand
			case 3: // MountainBiome
				return block.Stone // Replace Granite with Stone
			default:
				return block.Stone
			}
		}
	} else {
		// Deep stone with occasional variation
		stoneVariation := terrainNoise.FractalNoise2D(x*0.03, y*0.03, 2, 0.02, 1.0, 0.5)

		if stoneVariation > 0.7 {
			// Stone variants based on biome
			switch biome.Type {
			case 3: // MountainBiome
				return block.Stone // Replace Granite with Stone
			case 7: // OceanBiome (if underground)
				return block.Stone // Replace Marble with Stone
			case 2: // DesertBiome
				return block.Sand // Replace Sandstone with Sand
			default:
				return block.Stone
			}
		} else if stoneVariation > 0.3 {
			return block.Stone
		} else {
			// Add some clay for variety instead of limestone
			return block.Clay
		}
	}
}

// Tree placement with biome-specific logic
func shouldPlaceTree(globalX int, biome noise.BiomeData) bool {
	treeChance := settings.TreeChance * 0.33 // Base tree chance (reduced from default)

	switch biome.Type {
	case 1: // ForestBiome
		treeChance = settings.TreeChance * 1.67 // 25% (0.15 * 1.67)
	case 6: // JungleBiome
		treeChance = settings.TreeChance * 2.0 // 30% (0.15 * 2.0)
	case 0: // PlainseBiome
		treeChance = settings.TreeChance * 0.53 // 8% (0.15 * 0.53)
	case 4: // SwampBiome
		treeChance = settings.TreeChance // Use default 15%
	case 2, 5: // DesertBiome, TundraBiome
		treeChance = settings.TreeChance * 0.067 // 1% (0.15 * 0.067)
	case 3: // MountainBiome
		if biome.Temperature > -0.3 {
			treeChance = settings.TreeChance * 0.8 // 12% (0.15 * 0.8)
		} else {
			treeChance = settings.TreeChance * 0.133 // 2% (0.15 * 0.133)
		}
	}

	// Use position-based deterministic randomness
	hash := float64(((globalX*73856093)^(globalX*19349663))%1000000) / 1000000.0
	return hash < treeChance
}

func shouldPlaceTreeByID(globalX int, biomeID int) bool {
	treeChance := settings.TreeChance * 0.33 // Base tree chance
	switch biomeID {
	case 1:
		treeChance = settings.TreeChance * 1.67 // 25%
	case 6:
		treeChance = settings.TreeChance * 2.0 // 30%
	case 0:
		treeChance = settings.TreeChance * 0.53 // 8%
	case 4:
		treeChance = settings.TreeChance // 15%
	case 2, 5:
		treeChance = settings.TreeChance * 0.067 // 1%
	case 3:
		treeChance = settings.TreeChance * 0.47 // 7%
	}
	hash := float64(((globalX*73856093)^(globalX*19349663))%1000000) / 1000000.0
	return hash < treeChance
}

// Surface block selection based on biome
func getSurfaceBlock(biome noise.BiomeData) block.BlockType {
	switch biome.Type {
	case 0: // PlainseBiome
		return block.Grass
	case 1: // ForestBiome
		return block.Grass
	case 2: // DesertBiome
		return block.Sand
	case 3: // MountainBiome
		if biome.Temperature < -0.2 {
			return block.Snow
		} else if biome.Elevation > 0.5 {
			return block.Stone
		}
		return block.Grass
	case 4: // SwampBiome
		return block.Mud
	case 5: // TundraBiome
		return block.Snow
	case 6: // JungleBiome
		return block.Grass
	case 7: // OceanBiome
		return block.Water
	default:
		return block.Grass
	}
}

// Helper function to generate chunk key
func getChunkKey(chunkX, chunkY int) string {
	// Use a string format to avoid collisions and support negative coordinates
	return fmt.Sprintf("%d,%d", chunkX, chunkY)
}

// Standalone biome and terrain height functions using PerlinNoise
func GetBiomeAt(noise *noise.PerlinNoise, x float64) int {
	temp := noise.Noise1D(x * 0.001)
	hum := noise.Noise1D(x*0.0015 + 1000)
	elev := noise.Noise1D(x*0.0008 + 2000)

	// Simple biome selection based on noise
	if elev > 0.6 {
		return 3 // MountainBiome
	} else if temp < -0.4 {
		return 5 // TundraBiome
	} else if hum < -0.4 {
		return 2 // DesertBiome
	} else if temp > 0.5 && hum > 0.3 {
		return 6 // JungleBiome
	} else if hum > 0.2 && temp > -0.2 {
		return 1 // ForestBiome
	} else if hum > 0.5 && elev < -0.2 {
		return 4 // SwampBiome
	} else if elev < -0.5 {
		return 7 // OceanBiome
	} else {
		return 0 // PlainsBiome
	}
}

func GetBiomeTerrainHeight(noise *noise.PerlinNoise, x float64, biome int) float64 {
	base := noise.Noise1D(x * 0.01)
	switch biome {
	case 3: // Mountain - Extremely tall and dramatic
		mountains := noise.Noise1D(x*0.008) * 2.5
		peaks := noise.RidgedNoise1D(x*0.02, 3, 0.02, 1.0) * 1.8
		return base*0.8 + mountains + peaks
	case 0: // Plains - Rolling hills with dramatic variations
		hills := noise.Noise1D(x*0.015) * 1.5
		rolling := noise.Noise1D(x*0.04) * 0.8
		return base*0.6 + hills + rolling
	case 2: // Desert - Sharp dunes and mesas
		dunes := noise.Noise1D(x*0.025) * 1.2
		mesas := noise.RidgedNoise1D(x*0.035, 2, 0.03, 0.8) * 1.0
		return base*0.7 + dunes + mesas
	case 1: // Forest - Varied mountainous forest terrain
		forest := noise.Noise1D(x*0.012) * 1.8
		ridges := noise.Noise1D(x*0.03) * 0.9
		return base*0.5 + forest + ridges
	case 4: // Swamp - Low valleys with occasional hills
		valleys := noise.Noise1D(x*0.02) * 0.4
		lowlands := noise.Noise1D(x*0.06) * 0.3
		return base*0.2 + valleys + lowlands - 0.3
	case 5: // Tundra - Jagged frozen peaks
		frozen := noise.Noise1D(x*0.018) * 1.6
		jagged := noise.RidgedNoise1D(x*0.04, 2, 0.04, 0.7) * 1.2
		return base*0.4 + frozen + jagged
	case 6: // Jungle - Dense mountainous jungle terrain
		jungle := noise.Noise1D(x*0.014) * 2.0
		canopy := noise.Noise1D(x*0.035) * 1.1
		return base*0.6 + jungle + canopy
	case 7: // Ocean - Deep underwater terrain
		depths := noise.Noise1D(x*0.025) * 0.6
		return base*0.3 + depths - 1.2
	default:
		return base * 1.5
	}
}

// Terraria-style helper functions

// getTerrariaStyleSurfaceBlock returns biome-appropriate surface blocks
func getTerrariaStyleSurfaceBlock(biome noise.BiomeData, terrainNoise *noise.PerlinNoise, globalX int) block.BlockType {
	// Add some surface variation within biomes
	surfaceVariation := terrainNoise.Noise1D(float64(globalX)*0.03 + 500)

	switch biome.Type {
	case 0: // PlainseBiome
		if surfaceVariation > 0.7 {
			return block.Dirt // Exposed dirt patches
		}
		return block.Grass
	case 1: // ForestBiome
		return block.Grass
	case 2: // DesertBiome
		return block.Sand
	case 3: // MountainBiome
		if biome.Temperature < -0.2 {
			return block.Snow
		} else if biome.Elevation > 0.5 || surfaceVariation > 0.4 {
			return block.Stone // Rocky mountain peaks
		}
		return block.Grass
	case 4: // SwampBiome
		if surfaceVariation > 0.3 {
			return block.Mud
		}
		return block.Grass
	case 5: // TundraBiome
		return block.Snow
	case 6: // JungleBiome
		if surfaceVariation > 0.6 {
			return block.Mud // Jungle mud patches
		}
		return block.Grass
	case 7: // OceanBiome
		return block.Sand
	default:
		return block.Grass
	}
}

// getTerrariaStyleSoilBlock returns biome-appropriate soil blocks
func getTerrariaStyleSoilBlock(biome noise.BiomeData, terrainNoise *noise.PerlinNoise, globalX, globalY int) block.BlockType {
	soilNoise := terrainNoise.FractalNoise2D(float64(globalX)*0.05, float64(globalY)*0.05, 2, 0.08, 1.0, 0.5)

	switch biome.Type {
	case 2: // DesertBiome
		if soilNoise > 0.6 {
			return block.Clay // Hardpan clay layers
		}
		return block.Sand
	case 4: // SwampBiome
		if soilNoise > 0.2 {
			return block.Mud
		}
		return block.Dirt
	case 5: // TundraBiome
		if soilNoise > 0.5 {
			return block.Clay // Frozen clay
		}
		return block.Dirt
	default:
		if soilNoise > 0.7 {
			return block.Clay
		}
		return block.Dirt
	}
}

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

// generateTerrariaStyleOre creates Terraria-like ore distribution
func generateTerrariaStyleOre(globalX, globalY, depthFromSurface int, biome noise.BiomeData, terrainNoise *noise.PerlinNoise) block.BlockType {
	x, y := float64(globalX), float64(globalY)

	// Biome-based ore probability modifiers
	biomeOreMultiplier := 1.0
	switch biome.Type {
	case 3: // MountainBiome - richer in ores
		biomeOreMultiplier = 1.5
	case 2: // DesertBiome - fewer ores
		biomeOreMultiplier = 0.7
	case 4: // SwampBiome - limited ores
		biomeOreMultiplier = 0.8
	case 6: // JungleBiome - different ore distribution
		biomeOreMultiplier = 1.2
	}

	// Terraria-like depth-based ore layers
	if depthFromSurface >= 8 && depthFromSurface < 25 {
		// Shallow layer: Copper and some Iron
		copperNoise := terrainNoise.FractalNoise2D(x*0.09, y*0.08, 3, 0.12, 1.0, 0.6) * biomeOreMultiplier
		ironNoise := terrainNoise.FractalNoise2D(x*0.07+100, y*0.07, 2, 0.08, 0.8, 0.5) * biomeOreMultiplier

		if copperNoise > 0.78 {
			return block.CopperOre
		}
		if ironNoise > 0.85 && depthFromSurface > 15 {
			return block.IronOre
		}
	} else if depthFromSurface >= 25 && depthFromSurface < 60 {
		// Medium depth: Iron and Gold
		ironNoise := terrainNoise.FractalNoise2D(x*0.06+200, y*0.06, 3, 0.07, 1.0, 0.5) * biomeOreMultiplier
		goldNoise := terrainNoise.FractalNoise2D(x*0.05+300, y*0.05, 4, 0.05, 1.2, 0.4) * biomeOreMultiplier

		if goldNoise > 0.90 {
			return block.GoldOre
		}
		if ironNoise > 0.80 {
			return block.IronOre
		}
		// Still some copper at this depth but rarer
		copperNoise := terrainNoise.FractalNoise2D(x*0.08+400, y*0.08, 2, 0.1, 0.6, 0.6) * biomeOreMultiplier
		if copperNoise > 0.88 {
			return block.CopperOre
		}
	} else if depthFromSurface >= 60 {
		// Deep layer: Mostly Gold with rare other ores
		goldNoise := terrainNoise.FractalNoise2D(x*0.04+500, y*0.04, 4, 0.04, 1.3, 0.3) * biomeOreMultiplier
		deepIronNoise := terrainNoise.FractalNoise2D(x*0.06+600, y*0.06, 3, 0.06, 1.0, 0.4) * biomeOreMultiplier

		if goldNoise > 0.87 {
			return block.GoldOre
		}
		if deepIronNoise > 0.85 {
			return block.IronOre
		}
	}

	return block.Air // No ore
}
