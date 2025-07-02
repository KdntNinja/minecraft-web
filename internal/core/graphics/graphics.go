package graphics

import (
	"embed"
	"image"
	_ "image/png"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

//go:embed *.png
var imageFiles embed.FS

var BlockTextures map[block.BlockType]*ebiten.Image

// AtlasCoord represents coordinates within an individual texture file's atlas
type AtlasCoord struct {
	X, Y float64 // Position in the texture's internal atlas grid (supports fractional positions)
}

// BlockTextureConfig maps block types to their file and atlas coordinates
type BlockTextureConfig struct {
	Filename string
	Coord    AtlasCoord // Coordinates within that specific texture file
}

// BlockTextureConfigs maps block types to their texture file and coordinates
var BlockTextureConfigs = map[block.BlockType]BlockTextureConfig{
	block.Dirt:      {Filename: "dirt.png", Coord: AtlasCoord{X: 0, Y: 0}},
	block.Grass:     {Filename: "grass.png", Coord: AtlasCoord{X: 2, Y: 0.3}},
	block.Clay:      {Filename: "clay.png", Coord: AtlasCoord{X: 0, Y: 0}},
	block.Leaves:    {Filename: "leaves.png", Coord: AtlasCoord{X: 0, Y: 0}},
	block.Wood:      {Filename: "wood.png", Coord: AtlasCoord{X: 0, Y: 0}},
	block.GoldOre:   {Filename: "goldore.png", Coord: AtlasCoord{X: 0, Y: 0}},
	block.CopperOre: {Filename: "copperore.png", Coord: AtlasCoord{X: 0, Y: 0}},
	block.IronOre:   {Filename: "ironore.png", Coord: AtlasCoord{X: 0, Y: 0}},
	block.Stone:     {Filename: "stone.png", Coord: AtlasCoord{X: 0, Y: 0}},
	block.Water:     {Filename: "water.png", Coord: AtlasCoord{X: 1, Y: 1}},
	block.Ash:       {Filename: "clay.png", Coord: AtlasCoord{X: 0, Y: 1}},
	block.Hellstone: {Filename: "goldore.png", Coord: AtlasCoord{X: 0, Y: 0}},
	block.Granite:   {Filename: "stone.png", Coord: AtlasCoord{X: 0, Y: 0}}, // Use stone.png, different tint
	block.Andesite:  {Filename: "stone.png", Coord: AtlasCoord{X: 0, Y: 0}}, // Use stone.png, different tint
	block.Diorite:   {Filename: "stone.png", Coord: AtlasCoord{X: 0, Y: 0}}, // Use stone.png, different tint
	block.Slate:     {Filename: "stone.png", Coord: AtlasCoord{X: 0, Y: 0}}, // Use stone.png, different tint
}

// LoadTextures loads all block textures from their individual atlas files
func LoadTextures(tileSize int) error {
	if BlockTextures != nil {
		return nil // Already loaded
	}

	BlockTextures = make(map[block.BlockType]*ebiten.Image)

	var wg sync.WaitGroup
	var mu sync.Mutex
	for blockType, config := range BlockTextureConfigs {
		wg.Add(1)
		go func(blockType block.BlockType, config BlockTextureConfig) {
			defer wg.Done()
			texture, err := loadTextureFromAtlas(config.Filename, config.Coord, tileSize)
			if err != nil {
				log.Printf("Warning: Could not load texture %s for block %v: %v", config.Filename, blockType, err)
				return
			}
			// Tint stone variants and gold ore
			switch blockType {
			case block.Granite:
				texture = tintImage(texture, 1.15, 0.95, 0.85) // Warm pinkish
			case block.Andesite:
				texture = tintImage(texture, 0.85, 0.9, 1.1) // Cool bluish
			case block.Diorite:
				texture = tintImage(texture, 1.2, 1.2, 1.2) // Bright white
			case block.Slate:
				texture = tintImage(texture, 0.6, 0.6, 0.7) // Dark gray
			case block.GoldOre:
				texture = tintImage(texture, 2.0, 2.0, 0.3) // Strong yellow tint
			}
			mu.Lock()
			BlockTextures[blockType] = texture
			mu.Unlock()
		}(blockType, config)
	}
	wg.Wait()
	log.Printf("Graphics textures loaded successfully: %d block textures", len(BlockTextures))
	return nil
}

// loadTextureFromAtlas loads a specific tile from an individual texture atlas file
func loadTextureFromAtlas(filename string, coord AtlasCoord, tileSize int) (*ebiten.Image, error) {
	// Load the texture file
	file, err := imageFiles.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	atlasImg := ebiten.NewImageFromImage(img)

	// Extract the specific tile from the atlas using settings.AtlasTileSize
	x := int(coord.X * float64(settings.AtlasTileSize))
	y := int(coord.Y * float64(settings.AtlasTileSize))

	// Extract the sub-image
	subImg := atlasImg.SubImage(image.Rect(x, y, x+settings.AtlasTileSize, y+settings.AtlasTileSize)).(*ebiten.Image)

	// Scale to the desired tile size
	return scaleTexture(subImg, tileSize), nil
}

// scaleTexture scales a texture to the specified size
func scaleTexture(sourceImg *ebiten.Image, tileSize int) *ebiten.Image {
	if sourceImg == nil {
		return nil
	}

	// If already the right size, return as-is
	if settings.AtlasTileSize == tileSize {
		// Create a copy to avoid reference issues
		newImg := ebiten.NewImage(tileSize, tileSize)
		newImg.DrawImage(sourceImg, &ebiten.DrawImageOptions{})
		return newImg
	}

	// Create a new image with the correct tile size
	scaledImg := ebiten.NewImage(tileSize, tileSize)

	// Calculate scale factor
	scale := float64(tileSize) / float64(settings.AtlasTileSize)

	// Draw the source image scaled to fit the tile size
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scale, scale)
	scaledImg.DrawImage(sourceImg, opts)

	return scaledImg
}

// GetBlockTexture returns the texture for a given block type
func GetBlockTexture(blockType block.BlockType) *ebiten.Image {
	if BlockTextures == nil {
		return nil
	}
	return BlockTextures[blockType]
}

// tintImage applies a simple RGB multiplier to an ebiten.Image
func tintImage(img *ebiten.Image, r, g, b float64) *ebiten.Image {
	if img == nil {
		return nil
	}
	w, h := img.Size()
	tinted := ebiten.NewImage(w, h)
	opts := &ebiten.DrawImageOptions{}
	opts.ColorM.Scale(r, g, b, 1)
	tinted.DrawImage(img, opts)
	return tinted
}
