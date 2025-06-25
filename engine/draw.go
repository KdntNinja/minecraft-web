package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileSize = 20
	tilesX   = 16
	tilesY   = 12
)

type BlockType int

const (
	Air BlockType = iota
	Grass
	Dirt
	Stone
)

func generateTerrain() [tilesY][tilesX]BlockType {
	var terrain [tilesY][tilesX]BlockType
	for y := 0; y < tilesY; y++ {
		for x := 0; x < tilesX; x++ {
			if y < 8 {
				terrain[y][x] = Air // Sky
			} else if y == 8 {
				terrain[y][x] = Grass
			} else if y > 8 && y < 11 {
				terrain[y][x] = Dirt
			} else if y >= 11 {
				terrain[y][x] = Stone
			}
		}
	}
	return terrain
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Lighter sky blue background

	terrain := generateTerrain()

	for y := 0; y < tilesY; y++ {
		for x := 0; x < tilesX; x++ {
			var c color.Color
			switch terrain[y][x] {
			case Grass:
				c = color.RGBA{106, 190, 48, 255}
			case Dirt:
				c = color.RGBA{151, 105, 79, 255}
			case Stone:
				c = color.RGBA{100, 100, 100, 255}
			case Air:
				c = color.RGBA{135, 206, 235, 255} // Slightly different blue for sky blocks
			}
			tile := ebiten.NewImage(tileSize, tileSize)
			tile.Fill(c)

			// Draw black border
			borderColor := color.Black
			for i := 0; i < tileSize; i++ {
				tile.Set(i, 0, borderColor)          // Top
				tile.Set(i, tileSize-1, borderColor) // Bottom
			}
			for i := 1; i < tileSize-1; i++ {
				tile.Set(0, i, borderColor)          // Left
				tile.Set(tileSize-1, i, borderColor) // Right
			}

			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			screen.DrawImage(tile, &op)
		}
	}
}
