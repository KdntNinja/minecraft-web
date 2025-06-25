package engine

type Game struct {
	Terrain [][]Block
	Width   int
	Height  int
}

// NewGame initializes the game with zero size; Layout will set the real size
func NewGame() *Game {
	return &Game{
		Terrain: [][]Block{},
		Width:   0,
		Height:  0,
	}
}
