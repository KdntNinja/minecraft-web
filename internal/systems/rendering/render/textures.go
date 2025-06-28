package render

import (
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

//go:embed ../../../core/graphics/atlas.png
var atlasData embed.FS

var atlasImage *ebiten.Image

// Atlas tile coordinates (x, y positions in the atlas)
// Based on the atlas layout, each tile appears to be 16x16 pixels
type AtlasCoord struct {
	X, Y int // Position in the atlas grid
}

// blockAtlasMap maps block types to their atlas coordinates
var blockAtlasMap = map[block.BlockType]AtlasCoord{
	block.Grass:     {X: 0, Y: 0},   // Top-left grass tile
	block.Dirt:      {X: 1, Y: 0},   // Brown dirt tile
	block.Sand:      {X: 2, Y: 0},   // Sandy tile
	block.Stone:     {X: 3, Y: 0},   // Gray stone tile
	block.Clay:      {X: 4, Y: 0},   // Clay tile
	block.Snow:      {X: 5, Y: 0},   // Snow/white tile
	block.Wood:      {X: 0, Y: 1},   // Wood log tile
	block.Leaves:    {X: 1, Y: 1},   // Green leaves tile
	block.CopperOre: {X: 2, Y: 1},   // Copper-colored ore
	block.IronOre:   {X: 3, Y: 1},   // Iron-colored ore
	block.GoldOre:   {X: 4, Y: 1},   // Gold-colored ore
	block.Mud:       {X: 5, Y: 1},   // Mud tile
	block.Water:     {X: 0, Y: 2},   // Water tile
	block.Ash:       {X: 1, Y: 2},   // Ash tile
	block.Hellstone: {X: 2, Y: 2},   // Hellstone tile
}

const (
	// Atlas tile size (appears to be 16x16 based on the atlas image)
	atlasTileSize = 16
)

// loadTextureFromAssets loads a texture from embedded assets
func loadTextureFromAssets(filename string) *ebiten.Image {
	filePath := "internal/core/graphics/Tiles/" + filename
	
	file, err := tileAssets.Open(filePath)
	if err != nil {
		log.Printf("Warning: Could not load texture %s: %v", filename, err)
		return nil
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Warning: Could not decode texture %s: %v", filename, err)
		return nil
	}

	return ebiten.NewImageFromImage(img)
}

// createScaledTexture creates a texture scaled to the tile size
func createScaledTexture(sourceImg *ebiten.Image) *ebiten.Image {
	if sourceImg == nil {
		return nil
	}

	// Create a new image with the correct tile size
	scaledImg := ebiten.NewImage(settings.TileSize, settings.TileSize)
	
	// Calculate scale factors
	srcW, srcH := sourceImg.Bounds().Dx(), sourceImg.Bounds().Dy()
	scaleX := float64(settings.TileSize) / float64(srcW)
	scaleY := float64(settings.TileSize) / float64(srcH)

	// Draw the source image scaled to fit the tile size
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scaleX, scaleY)
	scaledImg.DrawImage(sourceImg, opts)

	return scaledImg
}

// initTextureImages loads actual texture assets for blocks
func initTextureImages() {
	if isInitialized {
		return // Already initialized
	}

	// Initialize object pool for draw options
	initObjectPool()

	tileImages = make(map[block.BlockType]*ebiten.Image)
	batchRenderer = &ebiten.DrawImageOptions{}

	// Load texture images for all block types
	for blockType := block.Air; blockType <= block.Hellstone; blockType++ {
		if blockType == block.Air {
			continue // Skip air blocks
		}

		// Try to load texture from assets
		if textureName, exists := blockTextureMap[blockType]; exists {
			sourceTexture := loadTextureFromAssets(textureName)
			if sourceTexture != nil {
				scaledTexture := createScaledTexture(sourceTexture)
				if scaledTexture != nil {
					tileImages[blockType] = scaledTexture
					continue
				}
			}
		}

		// Fallback to solid color if texture loading fails
		log.Printf("Using fallback color for block type %v", blockType)
		tile := ebiten.NewImage(settings.TileSize, settings.TileSize)
		tile.Fill(getBlockColorFast(blockType))
		tileImages[blockType] = tile
	}

	isInitialized = true
}
