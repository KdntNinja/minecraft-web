package engine

// Layout implements ebiten.Game's Layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.LastScreenW = outsideWidth
	g.LastScreenH = outsideHeight
	return outsideWidth, outsideHeight
}
