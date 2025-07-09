package rendering

import (
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

// DrawEntities draws all entities (e.g. players) near the camera.
func DrawEntities(entities entity.Entities, screen *ebiten.Image, cameraX, cameraY float64, lastScreenW, lastScreenH int, playerImage *ebiten.Image) {
	// TODO: Move your entity rendering logic here, e.g. iterate and draw each entity.
}
