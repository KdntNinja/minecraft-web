package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/block"
)

var (
	tileImages    map[block.BlockType]*ebiten.Image
	batchRenderer *ebiten.DrawImageOptions // Reuse draw options to reduce allocations
)

func initTileImages() {
	tileImages = make(map[block.BlockType]*ebiten.Image)
	batchRenderer = &ebiten.DrawImageOptions{}

	// Create tile images for all block types
	for blockType := block.Air; blockType <= block.Lava; blockType++ {
		if blockType == block.Air {
			continue // Skip air blocks
		}
		tile := ebiten.NewImage(block.TileSize, block.TileSize)
		tile.Fill(BlockColor(blockType))
		tileImages[blockType] = tile
	}
}

func Draw(g *[][]block.Chunk, screen *ebiten.Image) {
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
					if blockType == block.Air {
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
func DrawWithCamera(g *[][]block.Chunk, screen *ebiten.Image, cameraX, cameraY float64) {
	if tileImages == nil {
		initTileImages()
	}
	screen.Fill(color.RGBA{135, 206, 250, 255}) // Sky

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Calculate visible tile bounds based on camera position
	startTileX := int(cameraX / float64(block.TileSize))
	endTileX := int((cameraX+float64(screenWidth))/float64(block.TileSize)) + 1
	startTileY := int(cameraY / float64(block.TileSize))
	endTileY := int((cameraY+float64(screenHeight))/float64(block.TileSize)) + 1

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
					if blockType == block.Air {
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

func BlockColor(b block.BlockType) color.Color {
	switch b {
	// Surface blocks
	case block.Grass:
		return color.RGBA{106, 190, 48, 255} // Green
	case block.Dirt:
		return color.RGBA{151, 105, 79, 255} // Brown
	case block.Sand:
		return color.RGBA{238, 203, 173, 255} // Sandy yellow
	case block.Clay:
		return color.RGBA{168, 85, 65, 255} // Reddish brown
	case block.Snow:
		return color.RGBA{255, 255, 255, 255} // White
	case block.Ice:
		return color.RGBA{173, 216, 230, 255} // Light blue

	// Stone variants
	case block.Stone:
		return color.RGBA{100, 100, 100, 255} // Gray
	case block.Granite:
		return color.RGBA{120, 120, 120, 255} // Light gray
	case block.Marble:
		return color.RGBA{245, 245, 245, 255} // Off-white
	case block.Obsidian:
		return color.RGBA{50, 50, 50, 255} // Dark gray/black

	// Ore blocks
	case block.CopperOre:
		return color.RGBA{184, 115, 51, 255} // Orange-brown
	case block.IronOre:
		return color.RGBA{192, 192, 192, 255} // Silver
	case block.SilverOre:
		return color.RGBA{211, 211, 211, 255} // Light silver
	case block.GoldOre:
		return color.RGBA{255, 215, 0, 255} // Gold
	case block.PlatinumOre:
		return color.RGBA{229, 228, 226, 255} // Platinum white

	// Underground blocks
	case block.Mud:
		return color.RGBA{101, 67, 33, 255} // Dark brown
	case block.Ash:
		return color.RGBA{128, 128, 128, 255} // Gray
	case block.Silt:
		return color.RGBA{139, 119, 101, 255} // Grayish brown

	// Cave blocks
	case block.Cobweb:
		return color.RGBA{220, 220, 220, 128} // Semi-transparent gray

	// Hell/Underworld blocks
	case block.Hellstone:
		return color.RGBA{139, 0, 0, 255} // Dark red
	case block.HellstoneOre:
		return color.RGBA{255, 69, 0, 255} // Orange-red

	// Tree blocks
	case block.Wood:
		return color.RGBA{139, 69, 19, 255} // Saddle brown
	case block.Leaves:
		return color.RGBA{34, 139, 34, 255} // Forest green

	// Liquids
	case block.Water:
		return color.RGBA{0, 191, 255, 180} // Semi-transparent blue
	case block.Lava:
		return color.RGBA{255, 69, 0, 255} // Orange-red

	case block.Air:
		return color.RGBA{135, 206, 235, 255} // Sky blue
	}
	return color.Black
}
