package render

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
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

	// Initialize object pool for draw options
	initObjectPool()

	// Pre-calculate colors for faster access
	initBlockColors()

	tileImages = make(map[block.BlockType]*ebiten.Image)
	batchRenderer = &ebiten.DrawImageOptions{}

	// Create tile images for all block types
	for blockType := block.Air; blockType <= block.Hellstone; blockType++ {
		if blockType == block.Air {
			continue // Skip air blocks
		}
		tile := ebiten.NewImage(settings.TileSize, settings.TileSize)
		tile.Fill(getBlockColorFast(blockType))
		tileImages[blockType] = tile
	}

	isInitialized = true
}
