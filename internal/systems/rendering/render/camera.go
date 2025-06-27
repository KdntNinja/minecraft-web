package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

// DrawWithCamera renders the world with camera offset for following player
func DrawWithCamera(g *[][]block.Chunk, screen *ebiten.Image, cameraX, cameraY float64) {
	if !isInitialized {
		initTileImages()
	}

	// Fill sky background
	screen.Fill(color.RGBA{135, 206, 250, 255})

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Calculate visible tile bounds based on camera position
	startTileX := int(cameraX / float64(block.TileSize))
	endTileX := int((cameraX+float64(screenWidth))/float64(block.TileSize)) + 2 // +2 for safety margin
	startTileY := int(cameraY / float64(block.TileSize))
	endTileY := int((cameraY+float64(screenHeight))/float64(block.TileSize)) + 2 // +2 for safety margin

	// Ensure bounds are not negative
	if startTileX < 0 {
		startTileX = 0
	}
	if startTileY < 0 {
		startTileY = 0
	}

	// Pre-calculate maximum bounds to avoid recalculating
	maxChunksY := len(*g)
	if maxChunksY == 0 {
		return
	}
	maxChunksX := len((*g)[0])
	if maxChunksX == 0 {
		return
	}

	maxTileX := maxChunksX * block.ChunkWidth
	maxTileY := maxChunksY * block.ChunkHeight

	// Clamp end bounds
	if endTileX > maxTileX {
		endTileX = maxTileX
	}
	if endTileY > maxTileY {
		endTileY = maxTileY
	}

	// Render only visible chunks for better performance
	for cy := 0; cy < maxChunksY; cy++ {
		for cx := 0; cx < maxChunksX; cx++ {
			// Skip chunks that are completely outside the view
			chunkStartX := cx * block.ChunkWidth
			chunkEndX := chunkStartX + block.ChunkWidth
			chunkStartY := cy * block.ChunkHeight
			chunkEndY := chunkStartY + block.ChunkHeight

			if chunkEndX < startTileX || chunkStartX > endTileX ||
				chunkEndY < startTileY || chunkStartY > endTileY {
				continue // Skip this chunk - it's not visible
			}

			chunk := (*g)[cy][cx]
			for y := 0; y < block.ChunkHeight; y++ {
				for x := 0; x < block.ChunkWidth; x++ {
					globalTileX := cx*block.ChunkWidth + x
					globalTileY := cy*block.ChunkHeight + y

					// Skip tiles outside visible area
					if globalTileX < startTileX || globalTileX >= endTileX ||
						globalTileY < startTileY || globalTileY >= endTileY {
						continue
					}

					blockType := chunk[y][x]
					if blockType == block.Air {
						continue // Skip air blocks
					}

					// Calculate screen position with camera offset
					px := float64(globalTileX*block.TileSize) - cameraX
					py := float64(globalTileY*block.TileSize) - cameraY

					// Final bounds check to ensure we're drawing on screen
					if px+float64(block.TileSize) < 0 || px >= float64(screenWidth) ||
						py+float64(block.TileSize) < 0 || py >= float64(screenHeight) {
						continue
					}

					tile := tileImages[blockType]
					if tile == nil {
						continue
					}

					// Reuse the batch renderer to reduce allocations
					batchRenderer.GeoM.Reset()
					batchRenderer.GeoM.Translate(px, py)
					screen.DrawImage(tile, batchRenderer)
				}
			}
		}
	}
}
