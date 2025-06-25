package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileSize = 20
)

type BlockType int

const (
	Air BlockType = iota
	Grass
	Dirt
	Stone
)

type Block struct {
	Type BlockType
	// Add more fields here for metadata, etc.
}

func generateTerrainDynamic(tilesY, tilesX int) [][]Block {
	terrain := make([][]Block, tilesY)
	for y := 0; y < tilesY; y++ {
		terrain[y] = make([]Block, tilesX)
		for x := 0; x < tilesX; x++ {
			var t BlockType
			if y < tilesY*2/3 {
				t = Air // Sky
			} else if y == tilesY*2/3 {
				t = Grass
			} else if y > tilesY*2/3 && y < tilesY-1 {
				t = Dirt
			} else if y >= tilesY-1 {
				t = Stone
			}
			terrain[y][x] = Block{Type: t}
		}
	}
	return terrain
}

// SetBlock sets a block at (x, y)
func (g *Game) SetBlock(x, y int, t BlockType) {
	if y >= 0 && y < g.Height && x >= 0 && x < g.Width {
		g.Terrain[y][x].Type = t
	}
}

// GetBlock returns the block at (x, y)
func (g *Game) GetBlock(x, y int) BlockType {
	if y >= 0 && y < g.Height && x >= 0 && x < g.Width {
		return g.Terrain[y][x].Type
	}
	return Air
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Lighter sky blue background

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			var c color.Color
			switch g.Terrain[y][x].Type {
			case Grass:
				c = color.RGBA{106, 190, 48, 255}
			case Dirt:
				c = color.RGBA{151, 105, 79, 255}
			case Stone:
				c = color.RGBA{100, 100, 100, 255}
			case Air:
				c = color.RGBA{135, 206, 235, 255}
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
