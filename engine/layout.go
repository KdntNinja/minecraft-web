package engine

const (
	ScreenWidth  = 1024
	ScreenHeight = 768
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
