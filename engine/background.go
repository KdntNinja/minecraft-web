package engine

import (
	"image/color"

	"github.com/KdntNinja/webcraft/settings"
)

// GetBackgroundColor returns the background color based on the player's Y position.
func GetBackgroundColor(playerY float64) color.RGBA {
	skyColor := color.RGBA{135, 206, 235, 255}      // Terraria-like sky blue
	undergroundColor := color.RGBA{10, 10, 30, 255} // Deep blue/black

	if playerY <= settings.SkyTransitionStartY {
		return skyColor
	}
	if playerY >= settings.SkyTransitionEndY {
		return undergroundColor
	}
	t := (playerY - settings.SkyTransitionStartY) / (settings.SkyTransitionEndY - settings.SkyTransitionStartY)
	bgR := uint8(float64(skyColor.R)*(1-t) + float64(undergroundColor.R)*t)
	bgG := uint8(float64(skyColor.G)*(1-t) + float64(undergroundColor.G)*t)
	bgB := uint8(float64(skyColor.B)*(1-t) + float64(undergroundColor.B)*t)
	return color.RGBA{bgR, bgG, bgB, 255}
}
