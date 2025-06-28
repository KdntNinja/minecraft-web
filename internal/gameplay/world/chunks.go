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

// GetChunk retrieves a chunk from the world, generating it if necessary
func (w *World) GetChunk(chunkX, chunkY int) (block.Chunk, bool) {
	coord := ChunkCoord{X: chunkX, Y: chunkY}
	chunk, exists := w.Chunks[coord]
	if !exists {
		chunk = GenerateChunk(chunkX, chunkY)
		w.Chunks[coord] = chunk
		return chunk, false
	}
	return chunk, true
}

// LoadChunksAroundPlayer loads/generates chunks in a radius around the player
func (w *World) LoadChunksAroundPlayer(playerX, playerY float64, radius int) {
	playerChunkX := int(playerX) / (settings.ChunkWidth * settings.TileSize)
	playerChunkY := int(playerY) / (settings.ChunkHeight * settings.TileSize)

	for chunkY := playerChunkY - radius; chunkY <= playerChunkY+radius; chunkY++ {
		for chunkX := playerChunkX - radius; chunkX <= playerChunkX+radius; chunkX++ {
			coord := ChunkCoord{X: chunkX, Y: chunkY}
			if _, exists := w.Chunks[coord]; !exists {
				w.Chunks[coord] = GenerateChunk(chunkX, chunkY)
			}
		}
	}
}

// FindSurfaceHeight finds the Y coordinate of the surface at the given X coordinate in the world
func FindSurfaceHeight(worldX int, w *World) int {
	chunkX := worldX / settings.ChunkWidth
	inChunkX := worldX % settings.ChunkWidth

	// Search from top (lowest Y) to bottom (highest Y) for the first solid block
	for chunkY := 0; chunkY < 256; chunkY++ { // Arbitrary max height, can be adjusted
		coord := ChunkCoord{X: chunkX, Y: chunkY}
		chunk, exists := w.Chunks[coord]
		if !exists {
			chunk = GenerateChunk(chunkX, chunkY)
			w.Chunks[coord] = chunk
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
	return 50 // Default if no surface found
}
