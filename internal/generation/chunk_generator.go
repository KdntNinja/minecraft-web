package generation

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation/trees"
)

// GetWorldSeedFunc is set by the init function to avoid import cycles
var GetWorldSeedFunc func() int64

func init() {
	// Set the function variable to avoid import cycle
	GetWorldSeedFunc = GetSeed
}

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

	// Generate terrain for each column in the chunk (parallelized)
	var wg1 sync.WaitGroup
	for x := 0; x < settings.ChunkWidth; x++ {
		wg1.Add(1)
		go func(x int) {
			defer wg1.Done()
			worldX := chunkWorldX + x
			surfaceHeight := GetHeightAt(worldX)

			// Generate each block in this column from top to bottom (Y=0 is top)
			for chunkLocalY := 0; chunkLocalY < settings.ChunkHeight; chunkLocalY++ {
				worldY := chunkWorldY + chunkLocalY
				var blockType block.BlockType
				if worldY < surfaceHeight {
					// Above surface - air (already initialized)
					continue
				} else if worldY == surfaceHeight {
					if IsSurfaceCaveEntrance(worldX, worldY) {
						blockType = block.Air
					} else {
						blockType = GetSurfaceBlockType(worldX)
					}
				} else if worldY <= surfaceHeight+4 {
					if IsSurfaceCaveEntrance(worldX, worldY) {
						blockType = block.Air
					} else {
						blockType = GetShallowUndergroundBlock(worldX, worldY)
					}
				} else {
					if IsCave(worldX, worldY) {
						if IsLargeCavern(worldX, worldY) {
							if GetCaveWaterLevel(worldX, worldY) {
								blockType = block.Water
							} else {
								blockType = block.Air
							}
						} else {
							liquidType := IsLiquid(worldX, worldY)
							if liquidType > 0 {
								blockType = block.Water
							} else {
								blockType = block.Air
							}
						}
					} else {
						blockType = GetUndergroundBlock(worldX, worldY, surfaceHeight, rng)
					}
				}
				chunk[chunkLocalY][x] = blockType
			}
		}(x)
	}
	wg1.Wait()

	// Generate trees in a separate pass to avoid coordinate confusion (parallelized)
	var wg2 sync.WaitGroup
	for x := 0; x < settings.ChunkWidth; x++ {
		wg2.Add(1)
		go func(x int) {
			defer wg2.Done()
			worldX := chunkWorldX + x
			surfaceHeight := GetHeightAt(worldX)
			surfaceChunkY := surfaceHeight - chunkWorldY

			// Check if surface is in this chunk and is grass
			if surfaceChunkY >= 0 && surfaceChunkY < settings.ChunkHeight &&
				chunk[surfaceChunkY][x] == block.Grass &&
				rng.Float64() < settings.TreeChance &&
				x > 0 && x < settings.ChunkWidth-1 {

				trees.GenerateTreeAtPosition(&chunk, x, surfaceChunkY, rng)
			}
		}(x)
	}
	wg2.Wait()

	fmt.Printf("CHUNK_GEN: Completed chunk (%d, %d) with Perlin noise terrain\n", chunkX, chunkY)
	return chunk
}
