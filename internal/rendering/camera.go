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
	tileSize := settings.TileSize
	chunkWidth := settings.ChunkWidth
	chunkHeight := settings.ChunkHeight

	// Calculate visible tile bounds based on camera position
	startTileX := int(cameraX / float64(tileSize))
	endTileX := int((cameraX+float64(screenWidth))/float64(tileSize)) + 2 // +2 for safety margin
	startTileY := int(cameraY / float64(tileSize))
	endTileY := int((cameraY+float64(screenHeight))/float64(tileSize)) + 2 // +2 for safety margin

	// Calculate visible chunk bounds
	startChunkX := startTileX / chunkWidth
	endChunkX := (endTileX - 1) / chunkWidth
	startChunkY := startTileY / chunkHeight
	endChunkY := (endTileY - 1) / chunkHeight

	// Reuse draw options to reduce allocations
	drawOpts := getDrawOptions()

	// Precompute float versions for bounds
	fScreenWidth := float64(screenWidth)
	fScreenHeight := float64(screenHeight)
	fTileSize := float64(tileSize)

	for coord, chunk := range worldChunks {
		// Early culling: skip chunks outside the visible area
		if coord.X < startChunkX || coord.X > endChunkX || coord.Y < startChunkY || coord.Y > endChunkY {
			continue
		}
		baseTileX := coord.X * chunkWidth
		baseTileY := coord.Y * chunkHeight
		for y := 0; y < chunkHeight; y++ {
			globalTileY := baseTileY + y
			py := float64(globalTileY*tileSize) - cameraY
			// Skip entire row if off screen vertically
			if py+fTileSize < 0 || py >= fScreenHeight {
				continue
			}
			for x := 0; x < chunkWidth; x++ {
				globalTileX := baseTileX + x
				px := float64(globalTileX*tileSize) - cameraX
				// Skip tile if off screen horizontally
				if px+fTileSize < 0 || px >= fScreenWidth {
					continue
				}
				blockType := (*chunk)[y][x]
				if blockType == block.Air {
					continue // Skip air blocks
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
