package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/settings"
)

func Draw(g *[][]coretypes.Chunk, screen *ebiten.Image, cameraX, cameraY float64) {
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

	drawOpts := getDrawOptions() // Reuse one instance per frame

	for cy := 0; cy < len(*g); cy++ {
		for cx := 0; cx < len((*g)[cy]); cx++ {
			chunk := (*g)[cy][cx]
			baseTileX := cx * chunkWidth
			baseTileY := cy * chunkHeight
			for y := 0; y < chunkHeight; y++ {
				globalTileY := baseTileY + y
				py := float64(globalTileY*tileSize) - cameraY
				if globalTileY < startTileY || globalTileY >= endTileY || int(py) >= screenHeight {
					continue
				}
				for x := 0; x < chunkWidth; x++ {
					globalTileX := baseTileX + x
					px := float64(globalTileX*tileSize) - cameraX
					if globalTileX < startTileX || globalTileX >= endTileX || int(px) > screenWidth {
						continue
					}
					if y < 0 || y >= len(chunk.Blocks) || x < 0 || x >= len(chunk.Blocks[y]) {
						continue
					}
					blockType := chunk.Blocks[y][x]
					if blockType == coretypes.Air {
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
}
