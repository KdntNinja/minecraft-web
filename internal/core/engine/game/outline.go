package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// drawBlockOutline draws an outline around a block
func (g *Game) drawBlockOutline(screen *ebiten.Image, x, y int, outlineColor color.RGBA) {
	tileSize := settings.TileSize

	// Use ebitenutil.DrawLine for each edge (much faster than per-pixel Set)
	x0, y0 := float64(x), float64(y)
	x1, y1 := float64(x+tileSize-1), float64(y+tileSize-1)

	// Top edge
	ebitenutil.DrawLine(screen, x0, y0, x1, y0, outlineColor)
	// Bottom edge
	ebitenutil.DrawLine(screen, x0, y1, x1, y1, outlineColor)
	// Left edge
	ebitenutil.DrawLine(screen, x0, y0, x0, y1, outlineColor)
	// Right edge
	ebitenutil.DrawLine(screen, x1, y0, x1, y1, outlineColor)
}
