package world

import (
	"github.com/KdntNinja/webcraft/settings"
)

// wouldBlockCollideWithEntity checks if placing a block at the given coordinates would collide with any entity
func (w *World) wouldBlockCollideWithEntity(blockX, blockY int) bool {
	// Convert block coordinates to world coordinates
	blockWorldX := float64(blockX * settings.TileSize)
	blockWorldY := float64(blockY * settings.TileSize)
	blockWidth := float64(settings.TileSize)
	blockHeight := float64(settings.TileSize)

	// Check collision with all entities
	for _, entity := range w.Entities {
		entityX, entityY := entity.GetPosition()

		// Get entity dimensions based on interface (assume Player interface has GetAABB method)
		var entityWidth, entityHeight float64
		// Try to get AABB if available, else use default size
		type localAABB struct{ Width, Height int }
		type aabbGetter interface{ GetAABB() localAABB }
		if ag, ok := entity.(aabbGetter); ok {
			aabb := ag.GetAABB()
			entityWidth = float64(aabb.Width)
			entityHeight = float64(aabb.Height)
		} else {
			// Default entity size for other entity types
			entityWidth = float64(settings.TileSize)
			entityHeight = float64(settings.TileSize)
		}

		// Check AABB collision between entity and potential block position
		if entityX < blockWorldX+blockWidth &&
			entityX+entityWidth > blockWorldX &&
			entityY < blockWorldY+blockHeight &&
			entityY+entityHeight > blockWorldY {
			return true // Collision detected
		}
	}
	return false // No collision
}
