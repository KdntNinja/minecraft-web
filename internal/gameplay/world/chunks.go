package world

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation/terrain"
)

// GenerateChunk generates a chunk at the specified coordinates using the terrain generation system
func GenerateChunk(chunkX, chunkY int) block.Chunk {
	return terrain.GenerateChunk(chunkX, chunkY)
}

// FindSurfaceHeight finds the Y coordinate of the surface at the given X coordinate in the world
func FindSurfaceHeight(worldX int, w *World) int {
	chunkX := worldX / settings.ChunkWidth
	inChunkX := worldX % settings.ChunkWidth

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
				if blockType != block.Air {
					return globalY
				}
			}
		}
	}
	return settings.SurfaceBaseHeight // Use settings default if no surface found
}
