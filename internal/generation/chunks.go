package generation

import (
	"fmt"
	"math/rand"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

var GetWorldSeedFunc func() int64

// GenerateChunk creates a chunk with Minecraft-like Perlin noise terrain generation
func GenerateChunk(chunkX, chunkY int) block.Chunk {
	fmt.Printf("CHUNK_GEN: Generating chunk at (%d, %d) with Perlin noise\n", chunkX, chunkY)
	var chunk block.Chunk

	// Calculate world coordinates for this chunk
	chunkWorldX := chunkX * settings.ChunkWidth
	chunkWorldY := chunkY * settings.ChunkHeight

	// Create random generator for this chunk (for trees, etc.)
	chunkSeed := GetWorldSeedFunc() + int64(chunkX*1000+chunkY)
	rng := rand.New(rand.NewSource(chunkSeed))

	// Initialize chunk with air
	for y := 0; y < settings.ChunkHeight; y++ {
		for x := 0; x < settings.ChunkWidth; x++ {
			chunk[y][x] = block.Air
		}
	}

	// Generate terrain for each column in the chunk
	for x := 0; x < settings.ChunkWidth; x++ {
		worldX := chunkWorldX + x
		surfaceHeight := GetHeightAt(worldX)

		// Generate each block in this column from top to bottom (Y=0 is top)
		for chunkLocalY := 0; chunkLocalY < settings.ChunkHeight; chunkLocalY++ {
			// Calculate world Y coordinate (Y=0 is top of world)
			worldY := chunkWorldY + chunkLocalY

			var blockType block.BlockType

			if worldY < surfaceHeight {
				// Above surface - air (already initialized)
				continue
			} else if worldY == surfaceHeight {
				// Surface layer - always grass on top of dirt/earth
				blockType = block.Grass
			} else if worldY <= surfaceHeight+4 {
				// Shallow underground - determine dirt/clay layers
				initTerrainNoise() // Ensure noise is initialized
				biomeNoise := noiseGen.Noise2D(float64(worldX)/80.0, 0)
				noiseVal := noiseGen.Noise2D(float64(worldX)/10.0, float64(worldY)/10.0)

				if biomeNoise > 0.4 && noiseVal > 0.2 {
					blockType = block.Clay // Clay deposits in certain biomes
				} else {
					blockType = block.Dirt // Normal dirt layer under grass
				}
			} else {
				// Underground - check for caves first
				if IsCave(worldX, worldY) {
					blockType = block.Air
				} else {
					// Determine underground block type
					depthFromSurface := worldY - surfaceHeight

					// Check for ore veins
					oreType := GetOreType(worldX, worldY)
					if oreType > 0 && rng.Float64() < settings.OreVeinChance {
						switch oreType {
						case 1:
							blockType = block.CopperOre
						case 2:
							blockType = block.IronOre
						case 3:
							blockType = block.GoldOre
						default:
							blockType = block.Stone
						}
					} else {
						// Normal underground blocks based on depth
						if depthFromSurface <= 4 {
							// Shallow underground - mostly dirt with some clay variation
							initTerrainNoise() // Ensure noise is initialized
							noiseVal := noiseGen.Noise2D(float64(worldX)/10.0, float64(worldY)/10.0)
							if noiseVal > 0.3 {
								blockType = block.Clay
							} else {
								blockType = block.Dirt
							}
						} else if worldY < settings.ChunkHeight*settings.WorldChunksY-15 {
							// Stone layer with variation
							initTerrainNoise()

							// Add stone variation using noise
							stoneVariation := noiseGen.Noise2D(float64(worldX)/15.0, float64(worldY)/15.0)
							clayPockets := noiseGen.Noise2D(float64(worldX)/8.0+500, float64(worldY)/8.0+500)

							// Create clay pockets in stone
							if clayPockets > 0.6 {
								blockType = block.Clay
							} else if stoneVariation > 0.4 && depthFromSurface > 15 {
								// Deeper stone areas can have ash pockets
								if rng.Float64() < 0.15 {
									blockType = block.Ash
								} else {
									blockType = block.Stone
								}
							} else {
								blockType = block.Stone
							}
						} else if worldY < settings.ChunkHeight*settings.WorldChunksY-5 {
							// Deep underground transition zone
							initTerrainNoise()
							deepNoise := noiseGen.Noise2D(float64(worldX)/20.0, float64(worldY)/20.0)

							if deepNoise > 0.5 {
								blockType = block.Ash
							} else if deepNoise < -0.3 && rng.Float64() < 0.2 {
								// Rare water pockets in deep areas
								blockType = block.Water
							} else {
								blockType = block.Stone
							}
						} else {
							// Deepest levels - hellstone with some variation
							initTerrainNoise()
							hellNoise := noiseGen.Noise2D(float64(worldX)/12.0, float64(worldY)/12.0)

							if hellNoise > 0.3 {
								blockType = block.Hellstone
							} else if hellNoise < -0.4 {
								blockType = block.Ash
							} else {
								blockType = block.Stone
							}
						}
					}
				}
			}

			// Set the block in the chunk
			chunk[chunkLocalY][x] = blockType
		}
	}

	// Generate trees in a separate pass to avoid coordinate confusion
	for x := 0; x < settings.ChunkWidth; x++ {
		worldX := chunkWorldX + x
		surfaceHeight := GetHeightAt(worldX)
		surfaceChunkY := surfaceHeight - chunkWorldY

		// Check if surface is in this chunk and is grass
		if surfaceChunkY >= 0 && surfaceChunkY < settings.ChunkHeight &&
			chunk[surfaceChunkY][x] == block.Grass &&
			rng.Float64() < settings.TreeChance &&
			x > 0 && x < settings.ChunkWidth-1 {

			// Determine tree height with weighted random selection
			// Heights: 2 (most common), 3, 1, 6 (rarest)
			var treeHeight int
			heightRoll := rng.Float64()
			if heightRoll < 0.5 {
				treeHeight = 2 // 50% chance - most common
			} else if heightRoll < 0.8 {
				treeHeight = 3 // 30% chance
			} else if heightRoll < 0.95 {
				treeHeight = 1 // 15% chance
			} else {
				treeHeight = 6 // 5% chance - rarest
			}

			fmt.Printf("TREE_DEBUG: Placing tree (height %d) at chunk (%d,%d) worldX=%d, surfaceHeight=%d, surfaceChunkY=%d\n",
				treeHeight, chunkX, chunkY, worldX, surfaceHeight, surfaceChunkY)

			// Place wood trunk blocks
			for trunkLevel := 0; trunkLevel < treeHeight; trunkLevel++ {
				trunkChunkY := surfaceChunkY - trunkLevel
				if trunkChunkY >= 0 && trunkChunkY < settings.ChunkHeight {
					chunk[trunkChunkY][x] = block.Wood
				}
			}

			// Place leaves above the trunk (going up means decreasing world Y, so decreasing chunk Y)
			// Leaves start from the top of the trunk and go up a few more levels
			leafStartLevel := treeHeight
			leafEndLevel := treeHeight + 2 // 2-3 levels of leaves above trunk

			for leafLevel := leafStartLevel; leafLevel <= leafEndLevel; leafLevel++ {
				leafChunkY := surfaceChunkY - leafLevel
				if leafChunkY >= 0 && leafChunkY < settings.ChunkHeight {
					// Place center leaves
					chunk[leafChunkY][x] = block.Leaves

					// Place side leaves for most leaf levels (except the very top)
					if leafLevel < leafEndLevel {
						if x > 0 {
							chunk[leafChunkY][x-1] = block.Leaves
						}
						if x < settings.ChunkWidth-1 {
							chunk[leafChunkY][x+1] = block.Leaves
						}
					}
				}
			}
		}
	}

	fmt.Printf("CHUNK_GEN: Completed chunk (%d, %d) with Perlin noise terrain\n", chunkX, chunkY)
	return chunk
}
