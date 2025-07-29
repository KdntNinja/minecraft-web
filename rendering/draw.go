package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/settings"
)

func Draw(chunks map[coretypes.ChunkCoord]*coretypes.Chunk, screen *ebiten.Image, cameraX, cameraY float64, gridOffsetX, gridOffsetY int) {
	if tileImages == nil {
		initTileImages()
	}
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Sky

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	tileSize := settings.TileSize
	chunkWidth := settings.ChunkWidth
	chunkHeight := settings.ChunkHeight

	// Calculate visible tile bounds to avoid rendering off-screen tiles
	drawOpts := getDrawOptions() // Reuse one instance per frame

	for coord, chunk := range chunks {
		baseTileX := coord.X * chunkWidth
		baseTileY := coord.Y * chunkHeight
		for y := 0; y < chunkHeight; y++ {
			globalTileY := baseTileY + y
			py := float64(globalTileY*tileSize) - cameraY
			if int(py) >= screenHeight || int(py) < -tileSize {
				continue
			}
			for x := 0; x < chunkWidth; x++ {
				globalTileX := baseTileX + x
				px := float64(globalTileX*tileSize) - cameraX
				if int(px) >= screenWidth || int(px) < -tileSize {
					continue
				}
				if y < 0 || y >= len(chunk.Blocks) || x < 0 || x >= len(chunk.Blocks[y]) {
					continue
				}
				blockType := chunk.Blocks[y][x]
				if blockType == coretypes.Air {
					continue // Skip air blocks
				}
				tile := GetBlockTexture(blockType)
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
