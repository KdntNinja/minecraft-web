package game

import (
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
)

// UpdateEntitiesNearCamera updates only entities near the camera/screen for performance.
func (g *Game) UpdateEntitiesNearCamera() {
	// Pre-calculate camera bounds once
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			// Frustum culling for entities
			if p.X+float64(settings.PlayerColliderWidth) < camLeft || p.X > camRight ||
				p.Y+float64(settings.PlayerColliderHeight) < camTop || p.Y > camBottom {
				continue // Skip entities far from view
			}
			// Set the offset for the player's collision system
			p.AABB.GridOffsetX = g.physicsOffsetX
			p.AABB.GridOffsetY = g.physicsOffsetY

			// Update player movement
			p.Update()

			// Handle block interactions separately
			blockInteraction := p.HandleBlockInteractions(g.CameraX, g.CameraY)
			if blockInteraction != nil {
				switch blockInteraction.Type {
				case player.BreakBlock:
					g.World.BreakBlock(blockInteraction.BlockX, blockInteraction.BlockY)
				case player.PlaceBlock:
					g.World.PlaceBlock(blockInteraction.BlockX, blockInteraction.BlockY, p.SelectedBlock)
				}
			}

			// Use cached physics world
			if g.physicsWorld != nil {
				p.CollideBlocksAdvanced(g.physicsWorld)
			}
		}
	}
}
