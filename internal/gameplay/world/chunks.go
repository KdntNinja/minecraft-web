package world

import (
	"fmt"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation/terrain"
)

// GenerateChunk generates a chunk at the specified coordinates using the terrain generation system
func GenerateChunk(chunkX, chunkY int) block.Chunk {
	fmt.Printf("WORLD_CHUNK_GEN: Requesting chunk generation for (%d, %d)\n", chunkX, chunkY)
	return terrain.GenerateChunk(chunkX, chunkY)
}

// FindSurfaceHeight finds the Y coordinate of the surface at the given X coordinate in the world
func FindSurfaceHeight(worldX int, w *World) int {
	chunkX := worldX / settings.ChunkWidth
	inChunkX := worldX % settings.ChunkWidth

	fmt.Printf("SURFACE_DEBUG: FindSurfaceHeight for worldX=%d, chunkX=%d, inChunkX=%d\n", worldX, chunkX, inChunkX)

	// Search from top (lowest Y) to bottom (highest Y) for the first solid block
	// Only search in existing chunks since world is pre-generated
	for chunkY := 0; chunkY < settings.WorldChunksY; chunkY++ {
		coord := ChunkCoord{X: chunkX, Y: chunkY}
		chunk, exists := w.Chunks[coord]
		if !exists {
			continue // Skip non-existent chunks in fixed world
		}
		for y := 0; y < settings.ChunkHeight; y++ {
			globalY := chunkY*settings.ChunkHeight + y
			if inChunkX >= 0 && inChunkX < settings.ChunkWidth {
				blockType := chunk[y][inChunkX]

				// Debug: print first few blocks we encounter
				if chunkY <= 2 && y <= 5 {
					fmt.Printf("SURFACE_DEBUG: chunkY=%d, y=%d, globalY=%d, blockType=%d\n", chunkY, y, globalY, int(blockType))
				}

				if blockType != block.Air {
					fmt.Printf("SURFACE_DEBUG: Found surface at globalY=%d, blockType=%d\n", globalY, int(blockType))
					return globalY
				}
			}
		}
	}
	fmt.Printf("SURFACE_DEBUG: No surface found, using default=%d\n", settings.SurfaceBaseHeight)
	return settings.SurfaceBaseHeight // Use settings default if no surface found
}
