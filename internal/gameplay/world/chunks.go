package world

import (
	"fmt"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation"
)

// GenerateChunk generates a chunk at the specified coordinates using the terrain generation system
func GenerateChunk(chunkX, chunkY int) block.Chunk {
	fmt.Printf("WORLD_CHUNK_GEN: Requesting chunk generation for (%d, %d)\n", chunkX, chunkY)
	return generation.GenerateChunk(chunkX, chunkY)
}

// FindSurfaceHeight finds the Y coordinate of the surface at the given X coordinate in the world
func FindSurfaceHeight(worldX int, w *World) int {
	fmt.Printf("SURFACE_DEBUG: FindSurfaceHeight for worldX=%d\n", worldX)

	// Use the terrain generator's height calculation for consistency
	terrainHeight := generation.GetHeightAt(worldX)
	fmt.Printf("SURFACE_DEBUG: Terrain generator height for worldX=%d is %d\n", worldX, terrainHeight)

	// Verify by checking actual blocks in the world if needed
	chunkX := worldX / settings.ChunkWidth
	inChunkX := worldX % settings.ChunkWidth

	// Search around the calculated height to find the actual surface
	searchStart := terrainHeight - 5
	searchEnd := terrainHeight + 5

	if searchStart < 0 {
		searchStart = 0
	}

	for globalY := searchStart; globalY <= searchEnd; globalY++ {
		chunkY := globalY / settings.ChunkHeight
		localY := globalY % settings.ChunkHeight

		coord := ChunkCoord{X: chunkX, Y: chunkY}
		chunk, exists := w.Chunks[coord]
		if !exists {
			continue
		}

		if localY >= 0 && localY < settings.ChunkHeight && inChunkX >= 0 && inChunkX < settings.ChunkWidth {
			// Convert to chunk coordinate system (inverted Y)
			chunkLocalY := settings.ChunkHeight - 1 - localY
			blockType := chunk[chunkLocalY][inChunkX]

			// Check if this is the surface (solid block with air above)
			if blockType != block.Air {
				// Check if there's air above
				if globalY > 0 {
					aboveGlobalY := globalY - 1
					aboveChunkY := aboveGlobalY / settings.ChunkHeight
					aboveLocalY := aboveGlobalY % settings.ChunkHeight
					aboveCoord := ChunkCoord{X: chunkX, Y: aboveChunkY}

					if aboveChunk, aboveExists := w.Chunks[aboveCoord]; aboveExists {
						aboveChunkLocalY := settings.ChunkHeight - 1 - aboveLocalY
						if aboveLocalY >= 0 && aboveLocalY < settings.ChunkHeight {
							aboveBlockType := aboveChunk[aboveChunkLocalY][inChunkX]
							if aboveBlockType == block.Air {
								fmt.Printf("SURFACE_DEBUG: Found actual surface at globalY=%d, blockType=%d\n", globalY, int(blockType))
								return globalY
							}
						}
					}
				} else {
					// Top of world
					fmt.Printf("SURFACE_DEBUG: Found surface at top of world, globalY=%d, blockType=%d\n", globalY, int(blockType))
					return globalY
				}
			}
		}
	}

	fmt.Printf("SURFACE_DEBUG: Using terrain generator height=%d\n", terrainHeight)
	return terrainHeight
}
