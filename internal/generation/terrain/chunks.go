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

// GenerateChunk creates a chunk with Terraria-like terrain layers
func GenerateChunk(chunkX, chunkY int) block.Chunk {
	var chunk block.Chunk

	chunkStartY := chunkY * settings.ChunkHeight
	var seed int64 = settings.DefaultSeed
	if GetWorldSeedFunc != nil {
		seed = GetWorldSeedFunc()
	}
	terrainNoise := noise.NewPerlinNoise(seed)

	for y := 0; y < settings.ChunkHeight; y++ {
		globalY := chunkStartY + y
		for x := 0; x < settings.ChunkWidth; x++ {
			globalX := chunkX*settings.ChunkWidth + x

			// Add more noise for surface height and dirt thickness
			surfaceBase := 36
			surfaceNoise := terrainNoise.Noise1D(float64(globalX) * 0.015)
			surfaceDetail := terrainNoise.Noise1D(float64(globalX)*0.07) * 2.5
			surfaceHeight := surfaceBase + int(surfaceNoise*6+surfaceDetail)

			dirtBase := 6
			dirtNoise := terrainNoise.Noise1D(float64(globalX)*0.09+1000) * 2.5
			dirtThickness := dirtBase + int(dirtNoise)
			if dirtThickness < 4 {
				dirtThickness = 4
			}

			// --- BOUNDS CHECK: Prevent index out of range ---
			if y < 0 || y >= settings.ChunkHeight {
				continue
			}

			if globalY < surfaceHeight {
				chunk[y][x] = block.Air
			} else if globalY == surfaceHeight {
				// Surface block: add some variety (grass, sand, clay, snow)
				surfTypeNoise := terrainNoise.Noise1D(float64(globalX)*0.03 + 500)
				if surfTypeNoise > 0.6 {
					chunk[y][x] = block.Snow
				} else if surfTypeNoise < -0.6 {
					chunk[y][x] = block.Sand
				} else if surfTypeNoise > 0.2 {
					chunk[y][x] = block.Clay
				} else {
					chunk[y][x] = block.Grass
				}
			} else if globalY <= surfaceHeight+dirtThickness {
				// Dirt layer, with some clay and sand patches
				dirtTypeNoise := terrainNoise.Noise2D(float64(globalX)*0.05, float64(globalY)*0.05)
				if dirtTypeNoise > 0.7 {
					chunk[y][x] = block.Clay
				} else if dirtTypeNoise < -0.7 {
					chunk[y][x] = block.Sand
				} else {
					chunk[y][x] = block.Dirt
				}
			} else if globalY < settings.ChunkHeight*chunkY+settings.ChunkHeight-16 {
				// Stone layer, with caves and ores
				caveNoise := terrainNoise.Noise2D(float64(globalX)*settings.CaveFrequency, float64(globalY)*settings.CaveFrequency)
				if caveNoise > settings.CaveThreshold+0.05 { // Slightly higher than base threshold
					chunk[y][x] = block.Air
				} else {
					oreNoise := terrainNoise.Noise2D(float64(globalX)*0.13, float64(globalY)*0.13)
					if oreNoise > 0.82 {
						chunk[y][x] = block.IronOre
					} else if oreNoise < -0.82 {
						chunk[y][x] = block.CopperOre
					} else {
						chunk[y][x] = block.Stone
					}
				}
			} else if globalY < settings.ChunkHeight*chunkY+settings.ChunkHeight-4 {
				// Underworld transition (ash)
				chunk[y][x] = block.Ash
			} else {
				// Very bottom: hellstone
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
	case 3: // Mountain
		return base*0.3 + noise.Noise1D(x*0.02)*0.7
	case 0: // Plains
		return base*0.2 + noise.Noise1D(x*0.03)*0.8
	case 2: // Desert
		return base*0.4 + noise.Noise1D(x*0.04)
	case 1: // Forest
		return base + noise.Noise1D(x*0.015)
	case 4: // Swamp
		return base*0.1 + noise.Noise1D(x*0.02)*0.2
	case 5: // Tundra
		return base*0.2 + noise.Noise1D(x*0.025)*0.5
	case 6: // Jungle
		return base*0.3 + noise.Noise1D(x*0.018)*0.7
	case 7: // Ocean
		return base*0.1 - 0.5
	default:
		return base
	}
}
