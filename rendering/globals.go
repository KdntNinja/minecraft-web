package rendering

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/coretypes"
)

var (
	tileImages    map[coretypes.BlockType]*ebiten.Image
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
