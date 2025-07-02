package rendering

import (
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/graphics"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// initTextureImages loads textures using the graphics package
func initTextureImages() {
	if isInitialized {
		return // Already initialized
	}

	// Initialize object pool for draw options
	initObjectPool()

	// Load textures from the graphics package
	if err := graphics.LoadTextures(settings.TileSize); err != nil {
		if settings.TextureLogFallback {
			log.Printf("Error loading textures: %v", err)
			log.Printf("Falling back to solid colors")
		}
		initFallbackTextures()
		return
	}

	tileImages = make(map[block.BlockType]*ebiten.Image, block.NumBlockTypes)
	batchRenderer = &ebiten.DrawImageOptions{}

	// Copy textures from graphics package to render package
	for blockType := block.Air; blockType <= block.Hellstone; blockType++ {
		if blockType == block.Air {
			continue // Skip air blocks
		}

		// Try to get the texture from graphics package
		texture := graphics.GetBlockTexture(blockType)
		if texture != nil {
			tileImages[blockType] = texture
			continue
		}

		// Fallback to solid color if texture not found
		if settings.TextureLogFallback {
			log.Printf("Using fallback color for block type %v", blockType)
		}
		tile := ebiten.NewImage(settings.TileSize, settings.TileSize)
		tile.Fill(getBlockColorFast(blockType))
		tileImages[blockType] = tile
	}

	isInitialized = true
	if settings.TextureLogInit {
		log.Printf("Texture system initialized with %d block textures", len(tileImages))
	}
}

// initFallbackTextures creates solid color textures as a fallback
func initFallbackTextures() {
	tileImages = make(map[block.BlockType]*ebiten.Image)
	batchRenderer = &ebiten.DrawImageOptions{}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for blockType := block.Air; blockType <= block.Hellstone; blockType++ {
		if blockType == block.Air {
			continue // Skip air blocks
		}
		wg.Add(1)
		go func(bt block.BlockType) {
			defer wg.Done()
			tile := ebiten.NewImage(settings.TileSize, settings.TileSize)
			tile.Fill(getBlockColorFast(bt))
			mu.Lock()
			tileImages[bt] = tile
			mu.Unlock()
		}(blockType)
	}
	wg.Wait()

	isInitialized = true
	log.Printf("Fallback color-based texture system initialized")
}
