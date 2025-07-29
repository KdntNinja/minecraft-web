package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/KdntNinja/webcraft/settings"
)

// DrawBlockOutline draws an outline around a block
func DrawBlockOutline(screen *ebiten.Image, x, y int, outlineColor color.RGBA) {
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

// DrawUITextOutline draws text with a colored outline for readability
func DrawUITextOutline(screen *ebiten.Image, text string, x, y int, outline, fill color.Color) {
	// Draw outline (4 directions)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx != 0 || dy != 0 {
				ebitenutil.DebugPrintAt(screen, text, x+dx, y+dy)
			}
		}
	}
	// Draw fill
	ebitenutil.DebugPrintAt(screen, text, x, y)
}
