package terrain

import (
	"math"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/generation/noise"
)

// GenerateVisibleTerrain ensures chunks have some visible terrain for testing
func GenerateVisibleTerrain(chunkX, chunkY int) block.Chunk {
	var chunk block.Chunk

	// Simple visible terrain generation for immediate feedback
	terrainNoise := noise.NewSimplexNoise(42)

	// Calculate the global Y start for this chunk
	chunkStartY := chunkY * block.ChunkHeight

	for y := 0; y < block.ChunkHeight; y++ {
		for x := 0; x < block.ChunkWidth; x++ {
			globalX := chunkX*block.ChunkWidth + x
			globalY := chunkStartY + y

			// Create a simple surface at different heights based on noise
			surfaceNoise := terrainNoise.Noise1D(float64(globalX) * 0.02)
			// Surface height in absolute world coordinates
			surfaceHeight := 50 + int(surfaceNoise*8) // Surface between y=42 and y=58

			if globalY == surfaceHeight {
				// Surface
				biome := terrainNoise.GetBiomeAt(float64(globalX))
				chunk[y][x] = getSurfaceBlock(biome)
			} else if globalY < surfaceHeight {
				// Underground
				depth := surfaceHeight - globalY
				if depth <= 3 {
					chunk[y][x] = block.Dirt
				} else if depth <= 10 {
					// Mix of stone and dirt
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
				if globalY == surfaceHeight+1 && shouldPlaceTree(globalX, terrainNoise.GetBiomeAt(float64(globalX))) {
					chunk[y][x] = block.Wood
				} else if globalY == surfaceHeight+2 && shouldPlaceTree(globalX, terrainNoise.GetBiomeAt(float64(globalX))) {
					chunk[y][x] = block.Leaves
				} else {
					chunk[y][x] = block.Air
				}
			}
		}
	}

	return chunk
}

// GenerateChunk creates a chunk with Terraria-like terrain layers
func GenerateChunk(chunkX, chunkY int) block.Chunk {
	var chunk block.Chunk

	chunkStartY := chunkY * block.ChunkHeight
	terrainNoise := noise.NewSimplexNoise(42)

	for y := 0; y < block.ChunkHeight; y++ {
		globalY := chunkStartY + y
		for x := 0; x < block.ChunkWidth; x++ {
			globalX := chunkX*block.ChunkWidth + x

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
			} else if globalY < block.ChunkHeight*chunkY+block.ChunkHeight-16 {
				// Stone layer, with caves and ores
				caveNoise := terrainNoise.Noise2D(float64(globalX)*0.08, float64(globalY)*0.08)
				if caveNoise > 0.55 {
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
			} else if globalY < block.ChunkHeight*chunkY+block.ChunkHeight-4 {
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
func getSurfaceHeightRange(chunkX int, terrainNoise *noise.SimplexNoise) (int, int) {
	minHeight := 1000
	maxHeight := -1000

	// Sample a few points across the chunk width to get height range
	for x := 0; x < block.ChunkWidth; x++ {
		globalX := chunkX*block.ChunkWidth + x
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
func shouldGenerateCave(globalX, globalY, depthFromSurface int, noise *noise.SimplexNoise) bool {
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
func generateOre(globalX, globalY, depthFromSurface int, biome noise.BiomeData, terrainNoise *noise.SimplexNoise) block.BlockType {
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
			return block.SilverOre
		}
	} else if depthFromSurface >= 60 {
		// Deep ores: Platinum and rare ores
		platinumNoise := terrainNoise.FractalNoise2D(x*0.04+400, y*0.05, 4, 0.05, 1.4, 0.3) * biomeOreMultiplier

		if platinumNoise > 0.92 {
			return block.PlatinumOre
		}
	}

	return block.Air // No ore
}

// Get underground block based on depth, biome, and noise
func getUndergroundBlock(depthFromSurface int, biome noise.BiomeData, terrainNoise *noise.SimplexNoise, globalX, globalY int) block.BlockType {
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
				return block.Sandstone
			case 3: // MountainBiome
				return block.Granite
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
				return block.Granite
			case 7: // OceanBiome (if underground)
				return block.Marble
			case 2: // DesertBiome
				return block.Sandstone
			default:
				return block.Stone
			}
		} else if stoneVariation > 0.3 {
			return block.Stone
		} else {
			// Add some limestone for variety
			return block.Limestone
		}
	}
}

// Tree placement with biome-specific logic
func shouldPlaceTree(globalX int, biome noise.BiomeData) bool {
	treeChance := 0.05 // Increased base 5% chance for more visible trees

	switch biome.Type {
	case 1: // ForestBiome
		treeChance = 0.25 // 25% chance in forests
	case 6: // JungleBiome
		treeChance = 0.3 // 30% chance in jungles
	case 0: // PlainseBiome
		treeChance = 0.08 // 8% chance in plains
	case 4: // SwampBiome
		treeChance = 0.15 // 15% chance in swamps
	case 2, 5: // DesertBiome, TundraBiome
		treeChance = 0.01 // Still rare but not invisible
	case 3: // MountainBiome
		if biome.Temperature > -0.3 {
			treeChance = 0.12 // Some trees on warmer mountains
		} else {
			treeChance = 0.02 // Rare on cold mountains
		}
	}

	// Use position-based deterministic randomness
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

// ProceduralChunkLoader manages dynamic chunk loading around the player
type ProceduralChunkLoader struct {
	loadedChunks map[string]block.Chunk
	loadRadius   int
	lastPlayerX  int
	lastPlayerY  int
}

// NewProceduralChunkLoader creates a new chunk loader
func NewProceduralChunkLoader(loadRadius int) *ProceduralChunkLoader {
	return &ProceduralChunkLoader{
		loadedChunks: make(map[string]block.Chunk),
		loadRadius:   loadRadius,
		lastPlayerX:  -999999, // Force initial load
		lastPlayerY:  -999999,
	}
}

// UpdateAroundPlayer loads/unloads chunks based on player position - optimized for viewport
func (pcl *ProceduralChunkLoader) UpdateAroundPlayer(playerX, playerY float64) {
	playerChunkX := int(playerX) / (block.ChunkWidth * block.TileSize)
	playerChunkY := int(playerY) / (block.ChunkHeight * block.TileSize)

	// Only update if player moved significantly (every few blocks for performance)
	chunkDiffX := abs(playerChunkX - pcl.lastPlayerX)
	chunkDiffY := abs(playerChunkY - pcl.lastPlayerY)
	if chunkDiffX == 0 && chunkDiffY == 0 {
		return
	}

	pcl.lastPlayerX = playerChunkX
	pcl.lastPlayerY = playerChunkY

	// Calculate visible area (screen size in chunks)
	screenWidthInChunks := (block.TilesX*block.TileSize)/(block.ChunkWidth*block.TileSize) + 2 // +2 for buffer
	screenHeightInChunks := 3                                                                  // Only load a few chunks vertically for performance

	// Load chunks in visible area around player
	newChunks := make(map[string]block.Chunk)

	for dy := -screenHeightInChunks; dy <= screenHeightInChunks; dy++ {
		for dx := -screenWidthInChunks; dx <= screenWidthInChunks; dx++ {
			chunkX := playerChunkX + dx
			chunkY := playerChunkY + dy

			chunkKey := getChunkKey(chunkX, chunkY)

			// Check if chunk is already loaded
			if chunk, exists := pcl.loadedChunks[chunkKey]; exists {
				newChunks[chunkKey] = chunk
			} else {
				// Generate new chunk
				newChunks[chunkKey] = GenerateChunk(chunkX, chunkY)
			}
		}
	}

	pcl.loadedChunks = newChunks
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetLoadedChunks returns all currently loaded chunks
func (pcl *ProceduralChunkLoader) GetLoadedChunks() map[string]block.Chunk {
	return pcl.loadedChunks
}

// GetChunkAt returns the chunk at specific coordinates
func (pcl *ProceduralChunkLoader) GetChunkAt(chunkX, chunkY int) (block.Chunk, bool) {
	chunkKey := getChunkKey(chunkX, chunkY)
	chunk, exists := pcl.loadedChunks[chunkKey]
	return chunk, exists
}

// Helper function to generate chunk key
func getChunkKey(chunkX, chunkY int) string {
	return string(rune(chunkX*10000 + chunkY)) // Simple key generation
}
