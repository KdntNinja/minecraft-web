package engine

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	tilesX := outsideWidth / tileSize
	tilesY := outsideHeight / tileSize
	if tilesX != g.Width || tilesY != g.Height {
		g.Terrain = generateTerrainDynamic(tilesY, tilesX)
		g.Width = tilesX
		g.Height = tilesY
	}
	return outsideWidth, outsideHeight
}
