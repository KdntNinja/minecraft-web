package game

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
	"github.com/KdntNinja/webcraft/internal/gameplay/world"
	"github.com/KdntNinja/webcraft/internal/generation/terrain"
	"github.com/KdntNinja/webcraft/internal/systems/rendering/render"
)

type Game struct {
	World       *world.World
	LastScreenW int     // Cache last screen width
	LastScreenH int     // Cache last screen height
	CameraX     float64 // Camera X position
	CameraY     float64 // Camera Y position
	Seed        int64   // World seed for deterministic generation

	// Pre-allocated images to reduce memory allocation
	playerImage *ebiten.Image
	frameCount  int // For frame rate limiting
}

func NewGame() *Game {
	seed := time.Now().UnixNano() % 1000000
	g := &Game{
		LastScreenW: 800, // Default screen width
		LastScreenH: 600, // Default screen height
		Seed:        seed,
	}

	// Always reset world generation with the new seed
	terrain.ResetWorldGeneration(seed)

	// Pre-allocate player image to avoid recreating it every frame
	g.playerImage = ebiten.NewImage(settings.PlayerWidth, settings.PlayerHeight)
	g.playerImage.Fill(color.RGBA{255, 255, 0, 255}) // Yellow

	// Create a simple world with fixed size, passing the seed
	g.World = world.NewWorld(settings.ChunkHeight/2, 0, seed)

	return g
}

func (g *Game) Update() error {
	if g.World == nil {
		return nil
	}

	g.frameCount++

	// --- Dynamic chunk window management ---
	// Find player position (assume first entity is player)
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			g.World.UpdateChunksWindow(player.X, player.Y)
		}
	}

	// Update only entities near the camera/screen
	grid, gridOffsetX, gridOffsetY := g.World.ToIntGrid()
	visibleEntities := make([]interface {
		Update()
		CollideBlocks([][]int)
	}, 0, len(g.World.Entities))
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			if p.X+float64(settings.PlayerWidth) < camLeft || p.X > camRight ||
				p.Y+float64(settings.PlayerHeight) < camTop || p.Y > camBottom {
				continue // Skip entities far from view
			}
			// Set the offset for the player's collision system
			p.AABB.GridOffsetX = gridOffsetX
			p.AABB.GridOffsetY = gridOffsetY
			visibleEntities = append(visibleEntities, p)
		}
	}

	for _, e := range visibleEntities {
		e.Update()
		e.CollideBlocks(grid)
	}

	// Update camera to follow player (only every few frames for smoother performance)
	if g.frameCount%2 == 0 && len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			// Smooth camera following with some easing
			targetCameraX := player.X + float64(settings.PlayerWidth)/2 - float64(g.LastScreenW)/2
			targetCameraY := player.Y + float64(settings.PlayerHeight)/2 - float64(g.LastScreenH)/2

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
	render.DrawWithCamera(g.World.Chunks, screen, g.CameraX, g.CameraY)

	// Draw only entities near the camera/screen for performance
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			if p.X+float64(settings.PlayerWidth) < camLeft || p.X > camRight ||
				p.Y+float64(settings.PlayerHeight) < camTop || p.Y > camBottom {
				continue // Skip entities far from view
			}
			px, py := int(p.X-g.CameraX), int(p.Y-g.CameraY)
			if px > -settings.PlayerWidth && px < g.LastScreenW && py > -settings.PlayerHeight && py < g.LastScreenH {
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
