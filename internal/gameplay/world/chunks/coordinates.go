package chunks

import (
	"math"

	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// WorldToChunk converts world coordinates to chunk coordinates
func WorldToChunk(worldX, worldY float64) (int, int) {
	chunkX := int(math.Floor(worldX / float64(settings.ChunkWidth*settings.TileSize)))
	chunkY := int(math.Floor(worldY / float64(settings.ChunkHeight*settings.TileSize)))
	return chunkX, chunkY
}

// ChunkToWorld converts chunk coordinates to world coordinates (top-left corner)
func ChunkToWorld(chunkX, chunkY int) (float64, float64) {
	worldX := float64(chunkX * settings.ChunkWidth * settings.TileSize)
	worldY := float64(chunkY * settings.ChunkHeight * settings.TileSize)
	return worldX, worldY
}

// GetChunkCenter returns the center world coordinates of a chunk
func GetChunkCenter(chunkX, chunkY int) (float64, float64) {
	worldX, worldY := ChunkToWorld(chunkX, chunkY)
	centerX := worldX + float64(settings.ChunkWidth*settings.TileSize/2)
	centerY := worldY + float64(settings.ChunkHeight*settings.TileSize/2)
	return centerX, centerY
}

// BlockToChunk converts block coordinates to chunk coordinates
func BlockToChunk(blockX, blockY int) (int, int) {
	chunkX := int(math.Floor(float64(blockX) / float64(settings.ChunkWidth)))
	chunkY := int(math.Floor(float64(blockY) / float64(settings.ChunkHeight)))
	return chunkX, chunkY
}

// GetChunksInRadius returns all chunk coordinates within a given radius
func GetChunksInRadius(centerChunkX, centerChunkY, radius int) []ChunkCoord {
	var chunks []ChunkCoord

	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			// Optional: Use circular radius instead of square
			distance := math.Sqrt(float64(dx*dx + dy*dy))
			if distance <= float64(radius) {
				chunks = append(chunks, ChunkCoord{
					X: centerChunkX + dx,
					Y: centerChunkY + dy,
				})
			}
		}
	}

	return chunks
}

// GetChunkDistance calculates the distance between two chunks
func GetChunkDistance(chunk1, chunk2 ChunkCoord) float64 {
	dx := float64(chunk1.X - chunk2.X)
	dy := float64(chunk1.Y - chunk2.Y)
	return math.Sqrt(dx*dx + dy*dy)
}
