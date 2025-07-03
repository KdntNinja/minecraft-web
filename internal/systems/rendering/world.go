package rendering

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/hajimehoshi/ebiten/v2"
)

// DrawWorld draws the world chunks with camera offset.
// Accepts a map of chunk coordinates to chunk pointers, as in the game engine.
func DrawWorld(chunks map[[2]int]*block.Chunk, screen *ebiten.Image, cameraX, cameraY float64) {
	// TODO: Move your chunk rendering logic here, e.g. iterate and draw each chunk.
}
