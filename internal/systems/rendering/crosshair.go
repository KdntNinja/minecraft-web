package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawCrosshair draws a simple crosshair at the center of the screen.
func DrawCrosshair(screen *ebiten.Image, x, y float64) {
	// Example: simple white crosshair at the center
	w, h := screen.Size()
	cx, cy := w/2, h/2
	crosshairSize := 8
	crosshairColor := color.RGBA{255, 255, 255, 200}
	for i := -crosshairSize; i <= crosshairSize; i++ {
		px, py := cx+i, cy
		if px >= 0 && px < w && py >= 0 && py < h {
			screen.Set(px, py, crosshairColor)
		}
	}
	for i := -crosshairSize; i <= crosshairSize; i++ {
		px, py := cx, cy+i
		if px >= 0 && px < w && py >= 0 && py < h {
			screen.Set(px, py, crosshairColor)
		}
	}
}
