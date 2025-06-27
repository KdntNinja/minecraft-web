package world

import (
	"fmt"

	"github.com/KdntNinja/webcraft/engine/block"
)

func GenerateChunk(chunkX, chunkY int) block.Chunk {
	// Check cache first to avoid regenerating identical chunks
	cacheKey := fmt.Sprintf("%d,%d", chunkX, chunkY)
	if cached, exists := chunkCache[cacheKey]; exists {
		return cached
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
				// Above surface - air
				chunk[y][x] = block.Air
			} else if globalY == surface {
				// Surface block - depends on biome
				chunk[y][x] = getSurfaceBlock(globalX)
			} else {
				// Underground - complex generation
				depthFromSurface := globalY - surface

				// Check if we're in underworld first
				if isUnderworld(globalY, worldHeight) {
					// Underworld generation with more complex patterns
					underworldScale := 0.12
					underworldNoise2D := underworldNoise.Noise2D(float64(globalX)*underworldScale, float64(globalY)*underworldScale)

					// Add secondary noise for more varied underworld
					underworldNoise2 := underworldNoise.Noise2D(float64(globalX)*underworldScale*2, float64(globalY)*underworldScale*2) * 0.3
					combinedUnderworldNoise := underworldNoise2D + underworldNoise2

					if combinedUnderworldNoise < -0.3 {
						chunk[y][x] = block.Lava
					} else if combinedUnderworldNoise > 0.4 {
						chunk[y][x] = block.HellstoneOre
					} else {
						chunk[y][x] = block.Hellstone
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

				// Dirt layer - varies by biome
				var dirtThickness int
				switch biome {
				case DesertBiome:
					dirtThickness = 2 + int(dirtNoise.Noise2D(float64(globalX)*0.12, float64(globalY)*0.12)*1) // 1-3 blocks
				case SnowBiome:
					dirtThickness = 4 + int(dirtNoise.Noise2D(float64(globalX)*0.12, float64(globalY)*0.12)*2) // 2-6 blocks
				case ClayCanyonBiome:
					dirtThickness = 1 + int(dirtNoise.Noise2D(float64(globalX)*0.12, float64(globalY)*0.12)*1) // 0-2 blocks
				default: // Forest
					dirtThickness = 3 + int(dirtNoise.Noise2D(float64(globalX)*0.12, float64(globalY)*0.12)*2) // 1-5 blocks
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
