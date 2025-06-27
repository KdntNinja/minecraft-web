package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

func Draw(g *[][]block.Chunk, screen *ebiten.Image) {
	if tileImages == nil {
		initTileImages()
	}
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Sky

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Calculate visible tile bounds to avoid rendering off-screen tiles
	startTileX := 0
	endTileX := (screenWidth + block.TileSize - 1) / block.TileSize
	startTileY := 0
	endTileY := (screenHeight + block.TileSize - 1) / block.TileSize

	for cy := 0; cy < len(*g); cy++ {
		for cx := 0; cx < len((*g)[cy]); cx++ {
			chunk := (*g)[cy][cx]
			for y := 0; y < block.ChunkHeight; y++ {
				for x := 0; x < block.ChunkWidth; x++ {
					globalTileX := cx*block.ChunkWidth + x
					globalTileY := cy*block.ChunkHeight + y

					// Only render tiles that are within screen bounds
					if globalTileX < startTileX || globalTileX >= endTileX ||
						globalTileY < startTileY || globalTileY >= endTileY {
						continue
					}

					blockType := chunk[y][x]
					if blockType == block.Air {
						continue // Skip air blocks
					}

					px := globalTileX * block.TileSize
					py := globalTileY * block.TileSize

					// Double-check pixel bounds
					if px >= screenWidth || py >= screenHeight {
						continue
					}

					tile := tileImages[blockType]
					if tile == nil {
						continue
					}

					// Reuse the batch renderer to reduce allocations
					batchRenderer.GeoM.Reset()
					batchRenderer.GeoM.Translate(float64(px), float64(py))
					screen.DrawImage(tile, batchRenderer)
				}
			}
		}
	}
}
