package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

func Draw(g *[][]block.Chunk, screen *ebiten.Image) {
	if tileImages == nil {
		initTileImages()
	}
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Sky

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	tileSize := settings.TileSize
	chunkWidth := settings.ChunkWidth
	chunkHeight := settings.ChunkHeight

	// Calculate visible tile bounds to avoid rendering off-screen tiles
	startTileX := 0
	endTileX := (screenWidth + tileSize - 1) / tileSize
	startTileY := 0
	endTileY := (screenHeight + tileSize - 1) / tileSize

	for cy := 0; cy < len(*g); cy++ {
		for cx := 0; cx < len((*g)[cy]); cx++ {
			chunk := (*g)[cy][cx]
			baseTileX := cx * chunkWidth
			baseTileY := cy * chunkHeight
			for y := 0; y < chunkHeight; y++ {
				globalTileY := baseTileY + y
				py := globalTileY * tileSize
				if globalTileY < startTileY || globalTileY >= endTileY || py >= screenHeight {
					continue
				}
				for x := 0; x < chunkWidth; x++ {
					globalTileX := baseTileX + x
					px := globalTileX * tileSize
					if globalTileX < startTileX || globalTileX >= endTileX || px >= screenWidth {
						continue
					}
					blockType := chunk[y][x]
					if blockType == block.Air {
						continue // Skip air blocks
					}
					tile := tileImages[blockType]
					if tile == nil {
						continue
					}
					drawOpts := getDrawOptions()
					drawOpts.GeoM.Reset()
					drawOpts.GeoM.Translate(float64(px), float64(py))
					screen.DrawImage(tile, drawOpts)
				}
			}
		}
	}
}
