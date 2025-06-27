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

	// Pre-allocated images to reduce memory allocation
	playerImage *ebiten.Image
	frameCount  int // For frame rate limiting
}

func NewGame() *Game {
	g := &Game{
		LastScreenW: 800, // Default screen width
		LastScreenH: 600, // Default screen height
	}

	// Pre-allocate player image to avoid recreating it every frame
	g.playerImage = ebiten.NewImage(player.Width, player.Height)
	g.playerImage.Fill(color.RGBA{255, 255, 0, 255}) // Yellow

	// Create a simple world with fixed size
	g.World = world.NewWorld(20, 0) // Create world with 20 chunks vertically

	return g
}

func (g *Game) Update() error {
	if g.World == nil {
		return nil
	}

	g.frameCount++

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

	// Update camera to follow player (only every few frames for smoother performance)
	if g.frameCount%2 == 0 && len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			// Smooth camera following with some easing
			targetCameraX := player.X + float64(player.Width)/2 - float64(g.LastScreenW)/2
			targetCameraY := player.Y + float64(player.Height)/2 - float64(g.LastScreenH)/2

			// Smooth camera movement (lerp) - reduced for better performance
			lerpFactor := 0.05
			g.CameraX += (targetCameraX - g.CameraX) * lerpFactor
			g.CameraY += (targetCameraY - g.CameraY) * lerpFactor
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.World == nil {
		// Fill with a solid color to prevent black flashing
		screen.Fill(color.RGBA{135, 206, 250, 255}) // Sky blue
		return
	}

	// Clear screen with sky color to prevent flashing
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Sky blue

	// Render world directly to screen (avoid intermediate image allocation)
	render.DrawWithCamera(&g.World.Blocks, screen, g.CameraX, g.CameraY)

	// Draw all entities directly to screen with camera offset
	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			px, py := int(p.X-g.CameraX), int(p.Y-g.CameraY)

			// Only draw if player is visible on screen
			if px > -player.Width && px < g.LastScreenW && py > -player.Height && py < g.LastScreenH {
				var op ebiten.DrawImageOptions
				op.GeoM.Translate(float64(px), float64(py))
				screen.DrawImage(g.playerImage, &op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.LastScreenW = outsideWidth
	g.LastScreenH = outsideHeight
	return outsideWidth, outsideHeight
}
