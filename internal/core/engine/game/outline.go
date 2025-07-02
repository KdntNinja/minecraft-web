package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// drawBlockOutline draws an outline around a block
func (g *Game) drawBlockOutline(screen *ebiten.Image, x, y int, outlineColor color.RGBA) {
	tileSize := settings.TileSize

	// Top edge
	for i := 0; i < tileSize; i++ {
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y >= 0 && y < screen.Bounds().Dy() {
			screen.Set(x+i, y, outlineColor)
		}
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y+1 >= 0 && y+1 < screen.Bounds().Dy() {
			screen.Set(x+i, y+1, outlineColor)
		}
	}

	// Bottom edge
	for i := 0; i < tileSize; i++ {
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y+tileSize-1 >= 0 && y+tileSize-1 < screen.Bounds().Dy() {
			screen.Set(x+i, y+tileSize-1, outlineColor)
		}
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y+tileSize-2 >= 0 && y+tileSize-2 < screen.Bounds().Dy() {
			screen.Set(x+i, y+tileSize-2, outlineColor)
		}
	}

	// Left edge
	for i := 0; i < tileSize; i++ {
		if x >= 0 && x < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x, y+i, outlineColor)
		}
		if x+1 >= 0 && x+1 < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x+1, y+i, outlineColor)
		}
	}

	// Right edge
	for i := 0; i < tileSize; i++ {
		if x+tileSize-1 >= 0 && x+tileSize-1 < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x+tileSize-1, y+i, outlineColor)
		}
		if x+tileSize-2 >= 0 && x+tileSize-2 < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x+tileSize-2, y+i, outlineColor)
		}
	}
}
