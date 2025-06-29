package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/world/chunks"
)

// DrawWithCamera renders the world with camera offset for following player, using the chunk map
func DrawWithCamera(worldChunks map[chunks.ChunkCoord]*block.Chunk, screen *ebiten.Image, cameraX, cameraY float64) {
	if !isInitialized {
		initTileImages()
	}

	// Fill sky background
	screen.Fill(color.RGBA{135, 206, 250, 255})

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Calculate visible tile bounds based on camera position
	startTileX := int(cameraX / float64(settings.TileSize))
	endTileX := int((cameraX+float64(screenWidth))/float64(settings.TileSize)) + 2 // +2 for safety margin
	startTileY := int(cameraY / float64(settings.TileSize))
	endTileY := int((cameraY+float64(screenHeight))/float64(settings.TileSize)) + 2 // +2 for safety margin

	// Calculate visible chunk bounds horizontally
	startChunkX := startTileX / settings.ChunkWidth
	endChunkX := (endTileX - 1) / settings.ChunkWidth
	startChunkY := startTileY / settings.ChunkHeight
	endChunkY := (endTileY - 1) / settings.ChunkHeight

	// Reuse draw options to reduce allocations
	drawOpts := getDrawOptions()

	for coord, chunk := range worldChunks {
		// If any part of the chunk is in the visible horizontal window, render the whole chunk horizontally
		if coord.X < startChunkX || coord.X > endChunkX || coord.Y < startChunkY || coord.Y > endChunkY {
			continue
		}
		for y := 0; y < settings.ChunkHeight; y++ {
			for x := 0; x < settings.ChunkWidth; x++ {
				globalTileX := coord.X*settings.ChunkWidth + x
				globalTileY := coord.Y*settings.ChunkHeight + y
				blockType := (*chunk)[y][x]
				if blockType == block.Air {
					continue // Skip air blocks
				}
				// Calculate screen position with camera offset
				px := float64(globalTileX*settings.TileSize) - cameraX
				py := float64(globalTileY*settings.TileSize) - cameraY
				// Final bounds check to ensure we're drawing on screen
				if px+float64(settings.TileSize) < 0 || px >= float64(screenWidth) ||
					py+float64(settings.TileSize) < 0 || py >= float64(screenHeight) {
					continue
				}
				tile := tileImages[blockType]
				if tile == nil {
					continue
				}
				drawOpts.GeoM.Reset()
				drawOpts.GeoM.Translate(px, py)
				screen.DrawImage(tile, drawOpts)
			}
		}
	}
}
