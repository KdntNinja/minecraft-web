package render

import (
	"image/color"

	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/world"
	"github.com/hajimehoshi/ebiten/v2"
)

var tileImages map[world.BlockType]*ebiten.Image

func initTileImages() {
	tileImages = make(map[world.BlockType]*ebiten.Image)
	for _, t := range []world.BlockType{world.Grass, world.Dirt, world.Stone, world.Air} {
		tile := ebiten.NewImage(block.TileSize, block.TileSize)
		tile.Fill(BlockColor(t))
		tileImages[t] = tile
	}
}

func Draw(g *[][]world.Chunk, screen *ebiten.Image) {
	if tileImages == nil {
		initTileImages()
	}
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Sky

	screenWidth, screenHeight := screen.Size()

	for cy := 0; cy < len(*g); cy++ {
		for cx := 0; cx < len((*g)[cy]); cx++ {
			chunk := (*g)[cy][cx]
			for y := 0; y < block.ChunkHeight; y++ {
				for x := 0; x < block.ChunkWidth; x++ {
					px := (cx*block.ChunkWidth + x) * block.TileSize
					py := (cy*block.ChunkHeight + y) * block.TileSize
					if px+block.TileSize < 0 || px >= screenWidth || py+block.TileSize < 0 || py >= screenHeight {
						continue
					}
					tile := tileImages[chunk[y][x]]
					if tile == nil {
						continue
					}
					var op ebiten.DrawImageOptions
					op.GeoM.Translate(float64(px), float64(py))
					screen.DrawImage(tile, &op)
				}
			}
		}
	}
}

func BlockColor(b world.BlockType) color.Color {
	switch b {
	case world.Grass:
		return color.RGBA{106, 190, 48, 255}
	case world.Dirt:
		return color.RGBA{151, 105, 79, 255}
	case world.Stone:
		return color.RGBA{100, 100, 100, 255}
	case world.Air:
		return color.RGBA{135, 206, 235, 255}
	}
	return color.Black
}
