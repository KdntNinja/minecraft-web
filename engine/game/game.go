package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/player"
	"github.com/KdntNinja/webcraft/engine/render"
	"github.com/KdntNinja/webcraft/engine/world"
)

type Game struct {
	World      *world.World
	LoadedY    int // Number of vertical chunk rows loaded
	MaxChunksY int // Total vertical chunk rows to eventually load
	CenterX    int // Center chunk X
}

func NewGame() *Game {
	centerChunkX := 0
	initialChunksY := 10 // Start with enough chunks to fill most screens
	w := world.NewWorld(initialChunksY, centerChunkX)
	return &Game{
		World:      w,
		LoadedY:    initialChunksY,
		MaxChunksY: 0, // Will be set dynamically based on window height
		CenterX:    centerChunkX,
	}
}

func (g *Game) Update() error {
	// Update all entities (including player)
	grid := g.World.ToIntGrid()
	for _, e := range g.World.Entities {
		e.Update()
		if p, ok := e.(interface {
			CollideBlocks([][]int)
		}); ok {
			p.CollideBlocks(grid)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	render.Draw(&g.World.Blocks, screen)
	// Draw all entities (player, etc.)
	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			playerColor := [4]uint8{255, 255, 0, 255} // Yellow
			px, py := int(p.X), int(p.Y)
			img := ebiten.NewImage(player.Width, player.Height)
			img.Fill(color.RGBA{playerColor[0], playerColor[1], playerColor[2], playerColor[3]})
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(px), float64(py))
			screen.DrawImage(img, &op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Use exact browser window dimensions - no chunks, just pixel-perfect sizing
	// This ensures the game world matches exactly the browser window size
	return outsideWidth, outsideHeight
}
