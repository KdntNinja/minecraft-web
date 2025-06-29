package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
)

// drawCrosshair draws a targeting reticle and highlights the block under the cursor
func (g *Game) drawCrosshair(screen *ebiten.Image) {
	if len(g.World.Entities) == 0 {
		return
	}

	p, ok := g.World.Entities[0].(*player.Player)
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
	playerCenterX := p.X + float64(settings.PlayerColliderWidth)/2
	playerCenterY := p.Y + float64(settings.PlayerColliderHeight)/2
	blockCenterX := float64(blockX)*float64(settings.TileSize) + float64(settings.TileSize)/2
	blockCenterY := float64(blockY)*float64(settings.TileSize) + float64(settings.TileSize)/2
	dx := blockCenterX - playerCenterX
	dy := blockCenterY - playerCenterY
	distance := dx*dx + dy*dy
	inRange := distance <= p.InteractionRange*p.InteractionRange

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
	for i := -crosshairSize; i <= crosshairSize; i++ {
		px, py := mouseX+i, mouseY
		if px >= 0 && px < screen.Bounds().Dx() && py >= 0 && py < screen.Bounds().Dy() {
			screen.Set(px, py, crosshairColor)
		}
	}
	for i := -crosshairSize; i <= crosshairSize; i++ {
		px, py := mouseX, mouseY+i
		if px >= 0 && px < screen.Bounds().Dx() && py >= 0 && py < screen.Bounds().Dy() {
			screen.Set(px, py, crosshairColor)
		}
	}
}
