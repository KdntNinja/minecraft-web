package terrain

import (
	"fmt"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

var GetWorldSeedFunc func() int64

// GenerateChunk creates a chunk with authentic Terraria-like terrain layers
func GenerateChunk(chunkX, chunkY int) block.Chunk {
	fmt.Printf("CHUNK_GEN: Generating chunk at (%d, %d)\n", chunkX, chunkY)
	var chunk block.Chunk

	// Flat surface at mid-chunk
	localSurfaceY := settings.ChunkHeight / 2
	for y := 0; y < settings.ChunkHeight; y++ {
		invY := settings.ChunkHeight - 1 - y
		for x := 0; x < settings.ChunkWidth; x++ {
			// Simple block layering: grass, dirt, stone, bedrock
			if y > localSurfaceY {
				chunk[invY][x] = block.Air
			} else if y == localSurfaceY {
				chunk[invY][x] = block.Grass
			} else if y <= localSurfaceY+3 {
				chunk[invY][x] = block.Dirt
			} else if y < settings.ChunkHeight-1 {
				chunk[invY][x] = block.Stone
			} else {
				chunk[invY][x] = block.Hellstone
			}
		}
	}

	return chunk
}
