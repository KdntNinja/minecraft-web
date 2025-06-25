package engine

type Game struct {
	Terrain [][]Block
	Width   int
	Height  int
}

// NewGame initializes the game and terrain
func NewGame(screenWidth, screenHeight int) *Game {
	tilesX := screenWidth / tileSize
	tilesY := screenHeight / tileSize
	terrain := generateTerrainDynamic(tilesY, tilesX)
	return &Game{
		Terrain: terrain,
		Width:   tilesX,
		Height:  tilesY,
	}
}
