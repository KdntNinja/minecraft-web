package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/world"
)

var (
	tileImages    map[world.BlockType]*ebiten.Image
	batchRenderer *ebiten.DrawImageOptions // Reuse draw options to reduce allocations
)

func initTileImages() {
	tileImages = make(map[world.BlockType]*ebiten.Image)
	batchRenderer = &ebiten.DrawImageOptions{}

	for _, t := range []world.BlockType{world.Grass, world.Dirt, world.Stone, world.Air} {
		tile := ebiten.NewImage(block.TileSize, block.TileSize)
		tile.Fill(BlockColor(t))
		tileImages[t] = tile
	}
}

func Draw(g *[][]world.Chunk, screen *ebiten.Image) {
	if tileImages == nil {
		initTileImages()
	}
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Sky

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Calculate visible tile bounds to avoid rendering off-screen tiles
	startTileX := 0
	endTileX := (screenWidth + block.TileSize - 1) / block.TileSize
	startTileY := 0
	endTileY := (screenHeight + block.TileSize - 1) / block.TileSize

	for cy := 0; cy < len(*g); cy++ {
		for cx := 0; cx < len((*g)[cy]); cx++ {
			chunk := (*g)[cy][cx]
			for y := 0; y < block.ChunkHeight; y++ {
				for x := 0; x < block.ChunkWidth; x++ {
					globalTileX := cx*block.ChunkWidth + x
					globalTileY := cy*block.ChunkHeight + y

					// Only render tiles that are within screen bounds
					if globalTileX < startTileX || globalTileX >= endTileX ||
						globalTileY < startTileY || globalTileY >= endTileY {
						continue
					}

					blockType := chunk[y][x]
					if blockType == world.Air {
						continue // Skip air blocks
					}

					px := globalTileX * block.TileSize
					py := globalTileY * block.TileSize

					// Double-check pixel bounds
					if px >= screenWidth || py >= screenHeight {
						continue
					}

					tile := tileImages[blockType]
					if tile == nil {
						continue
					}

					// Reuse the batch renderer to reduce allocations
					batchRenderer.GeoM.Reset()
					batchRenderer.GeoM.Translate(float64(px), float64(py))
					screen.DrawImage(tile, batchRenderer)
				}
			}
		}
	}
}

// DrawWithCamera renders the world with camera offset for following player
func DrawWithCamera(g *[][]world.Chunk, screen *ebiten.Image, cameraX, cameraY float64) {
	if tileImages == nil {
		initTileImages()
	}
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Sky

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Calculate visible tile bounds based on camera position
	startTileX := int(cameraX / float64(block.TileSize))
	endTileX := int((cameraX + float64(screenWidth)) / float64(block.TileSize)) + 1
	startTileY := int(cameraY / float64(block.TileSize))
	endTileY := int((cameraY + float64(screenHeight)) / float64(block.TileSize)) + 1

	// Ensure bounds are not negative
	if startTileX < 0 {
		startTileX = 0
	}
	if startTileY < 0 {
		startTileY = 0
	}

	for cy := 0; cy < len(*g); cy++ {
		for cx := 0; cx < len((*g)[cy]); cx++ {
			chunk := (*g)[cy][cx]
			for y := 0; y < block.ChunkHeight; y++ {
				for x := 0; x < block.ChunkWidth; x++ {
					globalTileX := cx*block.ChunkWidth + x
					globalTileY := cy*block.ChunkHeight + y

					// Only render tiles that are within camera view
					if globalTileX < startTileX || globalTileX >= endTileX ||
						globalTileY < startTileY || globalTileY >= endTileY {
						continue
					}

					blockType := chunk[y][x]
					if blockType == world.Air {
						continue // Skip air blocks
					}

					// Calculate screen position with camera offset
					px := float64(globalTileX*block.TileSize) - cameraX
					py := float64(globalTileY*block.TileSize) - cameraY

					// Skip if outside screen bounds
					if px+float64(block.TileSize) < 0 || px >= float64(screenWidth) ||
						py+float64(block.TileSize) < 0 || py >= float64(screenHeight) {
						continue
					}

					tile := tileImages[blockType]
					if tile == nil {
						continue
					}

					// Reuse the batch renderer to reduce allocations
					batchRenderer.GeoM.Reset()
					batchRenderer.GeoM.Translate(px, py)
					screen.DrawImage(tile, batchRenderer)
				}
			}
		}
	}
}

func BlockColor(b world.BlockType) color.Color {
	switch b {
	case world.Grass:
		return color.RGBA{106, 190, 48, 255}
	case world.Dirt:
		return color.RGBA{151, 105, 79, 255}
	case world.Stone:
		return color.RGBA{100, 100, 100, 255}
	case world.Air:
		return color.RGBA{135, 206, 235, 255}
	}
	return color.Black
}
