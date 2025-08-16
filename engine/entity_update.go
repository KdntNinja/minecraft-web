package engine

import (
	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/physics"
	"github.com/KdntNinja/webcraft/settings"
)

// UpdateEntitiesNearCamera updates only entities near the camera/screen for performance.
func (g *Game) UpdateEntitiesNearCamera() {
	// Pre-calculate camera bounds once
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	for _, e := range g.World.GetEntities() {
		entity, ok := e.(interface {
			GetX() float64
			GetY() float64
			GetColliderWidth() float64
			GetColliderHeight() float64
			SetGridOffset(x, y int)
			Update()
			GetSelectedBlock() int
			CollideBlocks(*physics.PhysicsWorld)
			HandleBlockInteractions(cameraX, cameraY float64) interface {
				GetType() int
				GetBlockX() int
				GetBlockY() int
				GetSelectedBlock() int
			}
		})
		if !ok {
			continue
		}
		// Check if entity is within camera bounds
		if entity.GetX()+entity.GetColliderWidth() < camLeft || entity.GetX() > camRight ||
			entity.GetY()+entity.GetColliderHeight() < camTop || entity.GetY() > camBottom {
			continue
		}
		entity.SetGridOffset(g.physicsOffsetX, g.physicsOffsetY)
		entity.Update()
		blockInteraction := entity.HandleBlockInteractions(g.CameraX, g.CameraY)
		if blockInteraction != nil {
			blockType := blockInteraction.GetType()
			blockX := blockInteraction.GetBlockX()
			blockY := blockInteraction.GetBlockY()
			switch blockType {
			case 0: // BreakBlock
				g.World.BreakBlock(blockX, blockY)
			case 1: // PlaceBlock
				g.World.PlaceBlock(blockX, blockY, coretypes.BlockType(entity.GetSelectedBlock()))
			}
		}
		// Apply collision using cached physics world
		if g.physicsWorld != nil {
			entity.CollideBlocks(g.physicsWorld)
		}
	}
}
