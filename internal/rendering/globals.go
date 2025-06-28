package rendering

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

var (
	tileImages    map[block.BlockType]*ebiten.Image
	batchRenderer *ebiten.DrawImageOptions // Reuse draw options to reduce allocations
	isInitialized bool                     // Track initialization state
)

func initTileImages() {
	if isInitialized {
		return // Already initialized
	}

	// Use the new atlas-based texture initialization
	initTextureImages()
}
