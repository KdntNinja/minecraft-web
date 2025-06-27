package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/player"
	"github.com/KdntNinja/webcraft/engine/render"
	"github.com/KdntNinja/webcraft/engine/world"
)

type Game struct {
	World       *world.World
	LastScreenW int     // Cache last screen width
	LastScreenH int     // Cache last screen height
	CameraX     float64 // Camera X position
	CameraY     float64 // Camera Y position
}

func NewGame() *Game {
	g := &Game{
		LastScreenW: 800, // Default screen width
		LastScreenH: 600, // Default screen height
	}

	// Create a simple world with fixed size
	g.World = world.NewWorld(20, 0) // Create world with 20 chunks vertically

	return g
}

func (g *Game) Update() error {
	if g.World == nil {
		return nil
	}

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

	// Update camera to follow player
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			// Smooth camera following with some easing
			targetCameraX := player.X + float64(player.Width)/2 - float64(g.LastScreenW)/2
			targetCameraY := player.Y + float64(player.Height)/2 - float64(g.LastScreenH)/2

			// Smooth camera movement (lerp)
			lerpFactor := 0.1
			g.CameraX += (targetCameraX - g.CameraX) * lerpFactor
			g.CameraY += (targetCameraY - g.CameraY) * lerpFactor
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.World == nil {
		return
	}

	// Create a temporary image for world rendering
	worldImg := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	// Render world with camera offset
	render.DrawWithCamera(&g.World.Blocks, worldImg, g.CameraX, g.CameraY)

	// Draw all entities with camera offset
	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			playerColor := [4]uint8{255, 255, 0, 255} // Yellow
			px, py := int(p.X-g.CameraX), int(p.Y-g.CameraY)
			img := ebiten.NewImage(player.Width, player.Height)
			img.Fill(color.RGBA{playerColor[0], playerColor[1], playerColor[2], playerColor[3]})
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(px), float64(py))
			worldImg.DrawImage(img, &op)
		}
	}

	// Draw the world image to the screen
	screen.DrawImage(worldImg, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.LastScreenW = outsideWidth
	g.LastScreenH = outsideHeight
	return outsideWidth, outsideHeight
}
