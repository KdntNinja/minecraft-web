package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
	"github.com/KdntNinja/webcraft/internal/gameplay/world"
	"github.com/KdntNinja/webcraft/internal/generation"
	"github.com/KdntNinja/webcraft/internal/rendering"
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
	fpsCounter    int       // Frame counter for FPS calculation
	lastFPSUpdate time.Time // Last time FPS was calculated
	currentFPS    float64   // Current FPS value to display
}

func NewGame() *Game {
	seed := time.Now().UnixNano() % 1000000
	g := &Game{
		LastScreenW:   800, // Default screen width
		LastScreenH:   600, // Default screen height
		Seed:          seed,
		lastFPSUpdate: time.Now(), // Initialize FPS tracking
		currentFPS:    60.0,       // Default FPS value
	}

	// Hide the cursor for better gameplay experience
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	// Always reset world generation with the new seed
	generation.ResetWorldGeneration(seed)

	// Pre-allocate player image to avoid recreating it every frame
	g.playerImage = ebiten.NewImage(settings.PlayerWidth, settings.PlayerHeight)
	g.playerImage.Fill(color.RGBA{255, 255, 0, 255}) // Yellow

	// Create a simple world with fixed size, passing the seed
	g.World = world.NewWorld(seed)

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
	g.fpsCounter++

	// Update FPS calculation every second
	now := time.Now()
	if now.Sub(g.lastFPSUpdate) >= time.Second {
		g.currentFPS = float64(g.fpsCounter) / now.Sub(g.lastFPSUpdate).Seconds()
		g.fpsCounter = 0
		g.lastFPSUpdate = now
	}

	// Update only entities near the camera/screen - cached grid for better performance
	var grid [][]int
	var gridOffsetX, gridOffsetY int

	// Only regenerate collision grid when absolutely necessary (much less frequent)
	if g.frameCount%60 == 0 || g.World.IsGridDirty() {
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

			// Update player movement
			p.Update()

			// Handle block interactions separately
			blockInteraction := p.HandleBlockInteractions(g.CameraX, g.CameraY)
			if blockInteraction != nil {
				g.handleBlockInteraction(p, blockInteraction)
			}

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
	rendering.DrawWithCamera(g.World.Chunks, screen, g.CameraX, g.CameraY)

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

	// Draw crosshair/target indicator
	g.drawCrosshair(screen)

	// Display FPS in the top left corner
	fpsText := fmt.Sprintf("FPS: %.1f", g.currentFPS)
	ebitenutil.DebugPrint(screen, fpsText)

	// Display currently selected block and controls
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			selectedBlockText := fmt.Sprintf("\nSelected Block: %v", player.SelectedBlock)
			controlsText := "\nControls: Left Click = Break, Right Click = Place"
			numbersText := "\nBlocks: 1=Grass 2=Dirt 3=Clay 4=Stone 5=Copper 6=Iron 7=Gold 8=Ash 9=Wood 0=Leaves"
			uiText := fpsText + selectedBlockText + controlsText + numbersText
			ebitenutil.DebugPrint(screen, uiText)
		}
	} else {
		ebitenutil.DebugPrint(screen, fpsText)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.LastScreenW = outsideWidth
	g.LastScreenH = outsideHeight
	return outsideWidth, outsideHeight
}

// handleBlockInteraction processes block interaction events from the player
func (g *Game) handleBlockInteraction(p *player.Player, interaction *player.BlockInteraction) {
	switch interaction.Type {
	case player.BreakBlock:
		g.World.BreakBlock(interaction.BlockX, interaction.BlockY)
	case player.PlaceBlock:
		g.World.PlaceBlock(interaction.BlockX, interaction.BlockY, p.SelectedBlock)
	}
}

// drawCrosshair draws a targeting reticle and highlights the block under the cursor
func (g *Game) drawCrosshair(screen *ebiten.Image) {
	if len(g.World.Entities) == 0 {
		return
	}

	player, ok := g.World.Entities[0].(*player.Player)
	if !ok {
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()

	// Convert screen coordinates to world coordinates
	worldX := float64(mouseX) + g.CameraX
	worldY := float64(mouseY) + g.CameraY

	// Convert to block coordinates
	blockX := int(worldX / float64(settings.TileSize))
	blockY := int(worldY / float64(settings.TileSize))

	// Handle negative coordinates properly
	if worldX < 0 {
		blockX = int(worldX/float64(settings.TileSize)) - 1
	}
	if worldY < 0 {
		blockY = int(worldY/float64(settings.TileSize)) - 1
	}

	// Check if block is in range
	playerCenterX := player.X + float64(settings.PlayerWidth)/2
	playerCenterY := player.Y + float64(settings.PlayerHeight)/2
	blockCenterX := float64(blockX)*float64(settings.TileSize) + float64(settings.TileSize)/2
	blockCenterY := float64(blockY)*float64(settings.TileSize) + float64(settings.TileSize)/2

	dx := blockCenterX - playerCenterX
	dy := blockCenterY - playerCenterY
	distance := dx*dx + dy*dy
	inRange := distance <= player.InteractionRange*player.InteractionRange

	// Calculate screen position of the target block
	blockScreenX := float64(blockX*settings.TileSize) - g.CameraX
	blockScreenY := float64(blockY*settings.TileSize) - g.CameraY

	// Only draw if block is on screen
	if blockScreenX >= -float64(settings.TileSize) && blockScreenX < float64(screen.Bounds().Dx()) &&
		blockScreenY >= -float64(settings.TileSize) && blockScreenY < float64(screen.Bounds().Dy()) {

		// Create highlight color based on whether block is in range
		var highlightColor color.RGBA
		if inRange {
			highlightColor = color.RGBA{255, 255, 255, 128} // White semi-transparent
		} else {
			highlightColor = color.RGBA{255, 0, 0, 128} // Red semi-transparent (out of range)
		}

		// Draw block outline
		g.drawBlockOutline(screen, int(blockScreenX), int(blockScreenY), highlightColor)
	}

	// Draw simple crosshair at cursor
	crosshairSize := 8
	crosshairColor := color.RGBA{255, 255, 255, 200}

	// Horizontal line
	for i := -crosshairSize; i <= crosshairSize; i++ {
		px, py := mouseX+i, mouseY
		if px >= 0 && px < screen.Bounds().Dx() && py >= 0 && py < screen.Bounds().Dy() {
			screen.Set(px, py, crosshairColor)
		}
	}

	// Vertical line
	for i := -crosshairSize; i <= crosshairSize; i++ {
		px, py := mouseX, mouseY+i
		if px >= 0 && px < screen.Bounds().Dx() && py >= 0 && py < screen.Bounds().Dy() {
			screen.Set(px, py, crosshairColor)
		}
	}
}

// drawBlockOutline draws an outline around a block
func (g *Game) drawBlockOutline(screen *ebiten.Image, x, y int, outlineColor color.RGBA) {
	tileSize := settings.TileSize

	// Top edge
	for i := 0; i < tileSize; i++ {
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y >= 0 && y < screen.Bounds().Dy() {
			screen.Set(x+i, y, outlineColor)
		}
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y+1 >= 0 && y+1 < screen.Bounds().Dy() {
			screen.Set(x+i, y+1, outlineColor)
		}
	}

	// Bottom edge
	for i := 0; i < tileSize; i++ {
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y+tileSize-1 >= 0 && y+tileSize-1 < screen.Bounds().Dy() {
			screen.Set(x+i, y+tileSize-1, outlineColor)
		}
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y+tileSize-2 >= 0 && y+tileSize-2 < screen.Bounds().Dy() {
			screen.Set(x+i, y+tileSize-2, outlineColor)
		}
	}

	// Left edge
	for i := 0; i < tileSize; i++ {
		if x >= 0 && x < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x, y+i, outlineColor)
		}
		if x+1 >= 0 && x+1 < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x+1, y+i, outlineColor)
		}
	}

	// Right edge
	for i := 0; i < tileSize; i++ {
		if x+tileSize-1 >= 0 && x+tileSize-1 < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x+tileSize-1, y+i, outlineColor)
		}
		if x+tileSize-2 >= 0 && x+tileSize-2 < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x+tileSize-2, y+i, outlineColor)
		}
	}
}
