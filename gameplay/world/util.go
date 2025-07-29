package world

import (
	"github.com/KdntNinja/webcraft/settings"
)

// FindSurfaceHeight finds the surface height at the given X coordinate
func FindSurfaceHeight(blockX int, w *World) int {
	// Find which chunk this block belongs to
	chunkX := blockX / settings.ChunkWidth // TODO: Move BlockToChunk to a neutral package

	// We need to look through multiple chunks vertically to find surface
	for chunkY := 0; chunkY < 10; chunkY++ { // Search downward through chunks
		chunk := w.ChunkManager.GetChunk(chunkX, chunkY)

		// Look through this chunk for the surface
		localX := blockX - (chunkX * settings.ChunkWidth)
		if localX < 0 || localX >= settings.ChunkWidth {
			continue
		}

		for localY := 0; localY < settings.ChunkHeight; localY++ {
			if chunk.Blocks[localY][localX] != 0 {
				// Found the first non-air block - this is the surface
				return (chunkY * settings.ChunkHeight) + localY
			}
		}
	}

	// Default surface height if not found
	return settings.SurfaceBaseHeight
}
