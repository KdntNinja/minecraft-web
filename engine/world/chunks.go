package world

import (
	"fmt"

	"github.com/KdntNinja/webcraft/engine/block"
)

var (
	// Increase cache size for better performance
	maxCacheSize = 100
)

func GenerateChunk(chunkX, chunkY int) block.Chunk {
	// Check cache first to avoid regenerating identical chunks
	cacheKey := fmt.Sprintf("%d,%d", chunkX, chunkY)
	if cached, exists := chunkCache[cacheKey]; exists {
		return cached
	}

	// Limit cache size to prevent memory bloat
	if len(chunkCache) >= maxCacheSize {
		// Clear old entries (simple approach - clear half the cache)
		for k := range chunkCache {
			delete(chunkCache, k)
			if len(chunkCache) <= maxCacheSize/2 {
				break
			}
		}
	}

	var chunk block.Chunk
	// Optimize: calculate surface heights for entire chunk width at once
	surfaces := make([]int, block.ChunkWidth)
	for x := 0; x < block.ChunkWidth; x++ {
		globalX := chunkX*block.ChunkWidth + x
		surfaces[x] = getSurfaceHeight(globalX)
	}

	// Calculate world height for underworld detection
	worldHeight := 24 * 20 // Approximate world height

	// Generate chunk blocks using pre-calculated surface heights and multiple noise layers
	for y := 0; y < block.ChunkHeight; y++ {
		globalY := chunkY*block.ChunkHeight + y
		for x := 0; x < block.ChunkWidth; x++ {
			globalX := chunkX*block.ChunkWidth + x
			surface := surfaces[x]

			if globalY < surface {
				// Above surface - air, but check for trees
				if globalY == surface-1 && shouldPlaceTree(globalX) {
					// Place tree trunk
					chunk[y][x] = block.Wood
				} else if globalY == surface-2 && shouldPlaceTree(globalX) {
					// Place tree leaves
					chunk[y][x] = block.Leaves
				} else {
					chunk[y][x] = block.Air
				}
			} else if globalY == surface {
				// Surface block - depends on biome
				chunk[y][x] = getSurfaceBlock(globalX)
			} else {
				// Underground - complex generation
				depthFromSurface := globalY - surface // Check if we're in underworld first
				if isUnderworld(globalY, worldHeight) {
					// Terraria-style underworld generation
					underworldNoise := underworldNoise.TerrariaUnderworldNoise(float64(globalX), float64(globalY))

					if underworldNoise < -0.4 {
						chunk[y][x] = block.Lava
					} else if underworldNoise > 0.5 {
						chunk[y][x] = block.HellstoneOre
					} else if underworldNoise > 0.2 {
						chunk[y][x] = block.Hellstone
					} else {
						chunk[y][x] = block.Ash // Add ash blocks in underworld
					}
					continue
				}

				// Check for caves first
				if isInCave(globalX, globalY) {
					chunk[y][x] = block.Air
					continue
				}

				// Different underground layers based on biome and depth
				biome := getBiome(globalX)

				// Terraria-style underground layer generation
				undergroundNoise := dirtNoise.TerrariaUndergroundNoise(float64(globalX), float64(globalY))

				// Dirt layer - varies by biome and underground noise
				var dirtThickness int
				switch biome {
				case DesertBiome:
					// Sand extends deeper in desert
					dirtThickness = 3 + int(undergroundNoise*2) // 1-5 blocks of sand
				case SnowBiome:
					// Ice and snow layers
					dirtThickness = 5 + int(undergroundNoise*3) // 2-8 blocks
				case ClayCanyonBiome:
					// Thin soil layer over rock
					dirtThickness = 2 + int(undergroundNoise*1) // 1-3 blocks
				default: // Forest
					// Standard dirt layer
					dirtThickness = 4 + int(undergroundNoise*2) // 2-6 blocks
				}

				if depthFromSurface <= dirtThickness {
					// Dirt/mud layer based on biome
					switch biome {
					case DesertBiome:
						chunk[y][x] = block.Sand
					case SnowBiome:
						chunk[y][x] = block.Ice
					case ClayCanyonBiome:
						chunk[y][x] = block.Clay
					default:
						chunk[y][x] = block.Dirt
					}
				} else {
					// Stone layer with ore generation
					stoneBlock := getOreType(globalX, globalY, depthFromSurface)

					// Add some stone variation in deeper layers
					if depthFromSurface > 20 {
						stoneScale := 0.15
						stoneNoise2D := stoneNoise.Noise2D(float64(globalX)*stoneScale, float64(globalY)*stoneScale)

						if stoneNoise2D > 0.6 {
							stoneBlock = block.Granite
						} else if stoneNoise2D < -0.6 {
							stoneBlock = block.Marble
						}
					}

					// Very deep - add obsidian
					if depthFromSurface > 35 {
						if stoneNoise.Noise2D(float64(globalX)*0.2, float64(globalY)*0.2) > 0.7 {
							stoneBlock = block.Obsidian
						}
					}

					chunk[y][x] = stoneBlock
				}
			}
		}
	}

	// Cache the generated chunk
	chunkCache[cacheKey] = chunk
	return chunk
}
