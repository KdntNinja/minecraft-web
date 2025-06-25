package engine

import "log"

const (
	ScreenWidth  = 2048
	ScreenHeight = 1536
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	log.Printf("Layout called: outsideWidth=%d, outsideHeight=%d -> %d x %d", outsideWidth, outsideHeight, ScreenWidth, ScreenHeight)
	return ScreenWidth, ScreenHeight
}
