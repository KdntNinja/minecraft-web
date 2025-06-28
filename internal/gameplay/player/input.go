package player

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// HandleInput processes keyboard input and returns movement intentions
func (p *Player) HandleInput() (isMoving bool, targetVX float64, jumpKeyPressed bool) {
	isMoving = false
	targetVX = 0.0

	// Check horizontal movement keys (WASD + arrows)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		targetVX = -settings.PlayerMoveSpeed
		isMoving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		targetVX = settings.PlayerMoveSpeed
		isMoving = true
	}

	// Check jump keys (multiple options for accessibility)
	jumpKeyPressed = ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeySpace)

	return isMoving, targetVX, jumpKeyPressed
}
