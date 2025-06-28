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

	// Performance monitoring
	fpsCounter    int     // Frame counter for FPS calculation
	lastFPSUpdate float64 // Last time FPS was calculated
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
	g.World = world.NewWorld(0, seed)

	// Initialize camera position to follow the player's spawn location with tighter centering
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			// Center camera more tightly on the player for a zoomed-in feel
			g.CameraX = player.X + float64(settings.PlayerWidth)/2 - float64(g.LastScreenW)/2
			g.CameraY = player.Y + float64(settings.PlayerHeight)/2 - float64(g.LastScreenH)/2 - float64(settings.TileSize*2) // Offset upward slightly
		}
	}

	return g
}

func (g *Game) Update() error {
	if g.World == nil {
		return nil
	}

	g.frameCount++

	// Update only entities near the camera/screen - cached grid for better performance
	var grid [][]int
	var gridOffsetX, gridOffsetY int

	// Only regenerate collision grid when necessary (every few frames or when chunks change)
	if g.frameCount%4 == 0 || g.World.IsGridDirty() {
		grid, gridOffsetX, gridOffsetY = g.World.ToIntGrid()
	} else {
		grid, gridOffsetX, gridOffsetY = g.World.GetCachedGrid()
	}

	// Pre-calculate camera bounds once
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	// Update entities (reduce slice allocation by reusing)
	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			// Frustum culling for entities
			if p.X+float64(settings.PlayerWidth) < camLeft || p.X > camRight ||
				p.Y+float64(settings.PlayerHeight) < camTop || p.Y > camBottom {
				continue // Skip entities far from view
			}
			// Set the offset for the player's collision system
			p.AABB.GridOffsetX = gridOffsetX
			p.AABB.GridOffsetY = gridOffsetY

			p.Update()
			p.CollideBlocks(grid)
		}
	}

	// Update camera to follow player more responsively for zoomed-in feel
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			// Tighter camera following with offset for better view ahead
			targetCameraX := player.X + float64(settings.PlayerWidth)/2 - float64(g.LastScreenW)/2
			targetCameraY := player.Y + float64(settings.PlayerHeight)/2 - float64(g.LastScreenH)/2 - float64(settings.TileSize*2)

			// More responsive camera movement for zoomed-in feel
			lerpFactor := 0.12 // Increased from 0.05 for more responsive following
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

	// Pre-calculate camera bounds for entity culling
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	// Reusable draw options to reduce allocations
	var op ebiten.DrawImageOptions

	// Draw only entities near the camera/screen for performance
	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			// Use pre-calculated camera bounds
			if p.X+float64(settings.PlayerWidth) < camLeft || p.X > camRight ||
				p.Y+float64(settings.PlayerHeight) < camTop || p.Y > camBottom {
				continue // Skip entities far from view
			}
			px, py := int(p.X-g.CameraX), int(p.Y-g.CameraY)
			if px > -settings.PlayerWidth && px < g.LastScreenW && py > -settings.PlayerHeight && py < g.LastScreenH {
				// Reuse the DrawImageOptions instead of creating new ones
				op.GeoM.Reset()
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
