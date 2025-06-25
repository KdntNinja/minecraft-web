package engine

import "log"

const (
	ScreenWidth  = 1024
	ScreenHeight = 768
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	log.Printf("Layout called: outsideWidth=%d, outsideHeight=%d -> %d x %d", outsideWidth, outsideHeight, ScreenWidth, ScreenHeight)
	return ScreenWidth, ScreenHeight
}
