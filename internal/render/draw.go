package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/block"
)

var (
	tileImages    map[block.BlockType]*ebiten.Image
	batchRenderer *ebiten.DrawImageOptions // Reuse draw options to reduce allocations
	isInitialized bool                     // Track initialization state

	// Object pooling for performance
	drawOptionsPool []*ebiten.DrawImageOptions
	poolIndex       int

	// Pre-calculated colors for faster access
	blockColors      [33]color.RGBA // Pre-calculated array for all block types
	colorInitialized bool
)

func initTileImages() {
	if isInitialized {
		return // Already initialized
	}

	// Initialize object pool for draw options
	drawOptionsPool = make([]*ebiten.DrawImageOptions, 50)
	for i := range drawOptionsPool {
		drawOptionsPool[i] = &ebiten.DrawImageOptions{}
	}
	poolIndex = 0

	// Pre-calculate colors for faster access
	initBlockColors()

	tileImages = make(map[block.BlockType]*ebiten.Image)
	batchRenderer = &ebiten.DrawImageOptions{}

	// Create tile images for all block types
	for blockType := block.Air; blockType <= block.Lava; blockType++ {
		if blockType == block.Air {
			continue // Skip air blocks
		}
		tile := ebiten.NewImage(block.TileSize, block.TileSize)
		tile.Fill(getBlockColorFast(blockType))
		tileImages[blockType] = tile
	}

	isInitialized = true
}

func initBlockColors() {
	if colorInitialized {
		return
	}

	// Pre-calculate all block colors for faster access
	blockColors[block.Grass] = color.RGBA{106, 190, 48, 255}
	blockColors[block.Dirt] = color.RGBA{151, 105, 79, 255}
	blockColors[block.Sand] = color.RGBA{238, 203, 173, 255}
	blockColors[block.Clay] = color.RGBA{168, 85, 65, 255}
	blockColors[block.Snow] = color.RGBA{255, 255, 255, 255}
	blockColors[block.Ice] = color.RGBA{173, 216, 230, 255}
	blockColors[block.Stone] = color.RGBA{100, 100, 100, 255}
	blockColors[block.Granite] = color.RGBA{120, 120, 120, 255}
	blockColors[block.Marble] = color.RGBA{245, 245, 245, 255}
	blockColors[block.Obsidian] = color.RGBA{50, 50, 50, 255}
	blockColors[block.CopperOre] = color.RGBA{184, 115, 51, 255}
	blockColors[block.IronOre] = color.RGBA{192, 192, 192, 255}
	blockColors[block.SilverOre] = color.RGBA{211, 211, 211, 255}
	blockColors[block.GoldOre] = color.RGBA{255, 215, 0, 255}
	blockColors[block.PlatinumOre] = color.RGBA{229, 228, 226, 255}
	blockColors[block.Mud] = color.RGBA{101, 67, 33, 255}
	blockColors[block.Ash] = color.RGBA{128, 128, 128, 255}
	blockColors[block.Silt] = color.RGBA{139, 119, 101, 255}
	blockColors[block.Cobweb] = color.RGBA{220, 220, 220, 128}
	blockColors[block.Hellstone] = color.RGBA{139, 0, 0, 255}
	blockColors[block.HellstoneOre] = color.RGBA{255, 69, 0, 255}
	blockColors[block.Wood] = color.RGBA{139, 69, 19, 255}
	blockColors[block.Leaves] = color.RGBA{34, 139, 34, 255}
	blockColors[block.Water] = color.RGBA{0, 191, 255, 180}
	blockColors[block.Lava] = color.RGBA{255, 69, 0, 255}
	blockColors[block.Air] = color.RGBA{135, 206, 235, 255}

	colorInitialized = true
}

func getBlockColorFast(blockType block.BlockType) color.RGBA {
	if int(blockType) < len(blockColors) {
		return blockColors[blockType]
	}
	return color.RGBA{0, 0, 0, 255} // Black for unknown blocks
}

func getDrawOptions() *ebiten.DrawImageOptions {
	// Simple object pooling
	if poolIndex >= len(drawOptionsPool) {
		poolIndex = 0
	}
	options := drawOptionsPool[poolIndex]
	options.GeoM.Reset() // Reset transform
	poolIndex++
	return options
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
	if !isInitialized {
		initTileImages()
	}

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Calculate visible tile bounds based on camera position
	startTileX := int(cameraX / float64(block.TileSize))
	endTileX := int((cameraX+float64(screenWidth))/float64(block.TileSize)) + 2 // +2 for safety margin
	startTileY := int(cameraY / float64(block.TileSize))
	endTileY := int((cameraY+float64(screenHeight))/float64(block.TileSize)) + 2 // +2 for safety margin

	// Ensure bounds are not negative
	if startTileX < 0 {
		startTileX = 0
	}
	if startTileY < 0 {
		startTileY = 0
	}

	// Pre-calculate maximum bounds to avoid recalculating
	maxChunksY := len(*g)
	if maxChunksY == 0 {
		return
	}
	maxChunksX := len((*g)[0])
	if maxChunksX == 0 {
		return
	}

	maxTileX := maxChunksX * block.ChunkWidth
	maxTileY := maxChunksY * block.ChunkHeight

	// Clamp end bounds
	if endTileX > maxTileX {
		endTileX = maxTileX
	}
	if endTileY > maxTileY {
		endTileY = maxTileY
	}

	// Render only visible chunks for better performance
	for cy := 0; cy < maxChunksY; cy++ {
		for cx := 0; cx < maxChunksX; cx++ {
			// Skip chunks that are completely outside the view
			chunkStartX := cx * block.ChunkWidth
			chunkEndX := chunkStartX + block.ChunkWidth
			chunkStartY := cy * block.ChunkHeight
			chunkEndY := chunkStartY + block.ChunkHeight

			if chunkEndX < startTileX || chunkStartX > endTileX ||
				chunkEndY < startTileY || chunkStartY > endTileY {
				continue // Skip this chunk - it's not visible
			}

			chunk := (*g)[cy][cx]
			for y := 0; y < block.ChunkHeight; y++ {
				for x := 0; x < block.ChunkWidth; x++ {
					globalTileX := cx*block.ChunkWidth + x
					globalTileY := cy*block.ChunkHeight + y

					// Skip tiles outside visible area
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

					// Final bounds check to ensure we're drawing on screen
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
