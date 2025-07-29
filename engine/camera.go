package engine

import "github.com/KdntNinja/webcraft/settings"

// Camera logic and helpers for the Game struct

// GetTargetCameraPosition returns the target camera X and Y for the player.
func (g *Game) GetTargetCameraPosition(playerX, playerY float64) (float64, float64) {
	// Tighter camera following with offset for better view ahead
	targetCameraX := playerX + float64(g.LastScreenW)/2 - float64(g.LastScreenW)/2
	targetCameraY := playerY + float64(g.LastScreenH)/2 - float64(g.LastScreenH)/2 - float64(settings.TileSize*2)
	return targetCameraX, targetCameraY
}

// LerpCamera moves the camera towards the target position with a given lerp factor.
func (g *Game) LerpCamera(targetX, targetY float64, lerpFactor float64) {
	g.CameraX += (targetX - g.CameraX) * lerpFactor
	g.CameraY += (targetY - g.CameraY) * lerpFactor
}
