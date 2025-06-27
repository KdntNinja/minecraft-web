package world

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/generation/terrain"
)

// GenerateChunk generates a chunk at the specified coordinates using the terrain generation system
func GenerateChunk(chunkX, chunkY int) block.Chunk {
	return terrain.GenerateChunk(chunkX, chunkY)
}

// GetChunk retrieves a chunk from the world, generating it if necessary
func (w *World) GetChunk(chunkX, chunkY int) *block.Chunk {
	// Check if chunk exists in world
	if chunkY >= 0 && chunkY < len(w.Blocks) &&
		chunkX >= 0 && chunkX < len(w.Blocks[chunkY]) {
		chunk := GenerateChunk(chunkX, chunkY)
		w.Blocks[chunkY][chunkX] = chunk
		return &w.Blocks[chunkY][chunkX]
	}

	return nil
}

// LoadChunksAroundPlayer loads chunks in a radius around the player
func (w *World) LoadChunksAroundPlayer(playerX, playerY float64, radius int) {
	playerChunkX := int(playerX) / (block.ChunkWidth * block.TileSize)
	playerChunkY := int(playerY) / (block.ChunkHeight * block.TileSize)

	for chunkY := playerChunkY - radius; chunkY <= playerChunkY+radius; chunkY++ {
		for chunkX := playerChunkX - radius; chunkX <= playerChunkX+radius; chunkX++ {
			if chunkY >= 0 && chunkY < len(w.Blocks) &&
				chunkX >= 0 && chunkX < len(w.Blocks[chunkY]) {
				// Generate chunk if not already generated
				w.Blocks[chunkY][chunkX] = GenerateChunk(chunkX, chunkY)
			}
		}
	}
}
