package rendering

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/settings"
)

// initTextureImages loads textures using the graphics package
func initTextureImages() {
	if isInitialized {
		return // Already initialized
	}

	// Initialize object pool for draw options
	initObjectPool()

	// Load textures from the graphics package
	if err := LoadTextures(settings.TileSize); err != nil {
		log.Fatalf("Error loading textures: %v", err)
		return
	}

	tileImages = make(map[coretypes.BlockType]*ebiten.Image, coretypes.NumBlockTypes)
	batchRenderer = &ebiten.DrawImageOptions{}

	// Copy textures from graphics package to render package
	for blockType := coretypes.Air; blockType <= coretypes.Hellstone; blockType++ {
		if blockType == coretypes.Air {
			continue // Skip air blocks
		}

		// Try to get the texture from graphics package
		texture := GetBlockTexture(blockType)
		if texture != nil {
			tileImages[blockType] = texture
			continue
		}
	}

	isInitialized = true
	if settings.TextureLogInit {
		log.Printf("Texture system initialized with %d block textures", len(tileImages))
	}
}
