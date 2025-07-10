package rendering

import (
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
	"github.com/hajimehoshi/ebiten/v2"
)

// DrawEntities draws all entities (e.g. players) near the camera.
func DrawEntities(entities entity.Entities, screen *ebiten.Image, cameraX, cameraY float64, lastScreenW, lastScreenH int, playerImage *ebiten.Image) {
	for _, entity := range entities {
		// Get entity position (this is the collider position)
		entityX, entityY := entity.GetPosition()

		// Calculate screen position relative to camera
		screenX := entityX - cameraX
		screenY := entityY - cameraY

		// Only draw entities that are visible on screen (with some margin)
		margin := 100.0
		if screenX > -margin && screenX < float64(lastScreenW)+margin &&
			screenY > -margin && screenY < float64(lastScreenH)+margin {

			// Adjust sprite position based on entity type
			var spriteX, spriteY float64
			if _, isPlayer := entity.(*player.Player); isPlayer {
				// For players, the collider is centered horizontally and bottom-aligned
				// We need to draw the sprite at its original position
				spriteX = screenX - float64(settings.PlayerSpriteWidth-settings.PlayerColliderWidth)/2
				spriteY = screenY - float64(settings.PlayerSpriteHeight-settings.PlayerColliderHeight)
			} else {
				// For other entities, use the entity position directly
				spriteX = screenX
				spriteY = screenY
			}

			// Draw the entity image at the calculated screen position
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(spriteX, spriteY)
			screen.DrawImage(playerImage, op)
		}
	}
}
