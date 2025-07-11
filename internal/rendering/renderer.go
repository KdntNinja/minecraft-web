//go:build js && wasm
// +build js,wasm

package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/world/chunks"
)

// TileRenderJob represents a single tile rendering job
type TileRenderJob struct {
	blockType block.BlockType
	px, py    float64
	tile      *ebiten.Image
}

// DrawWithCameraAsync renders the world with camera offset (single-threaded, no concurrency)
func DrawWithCameraAsync(worldChunks map[chunks.ChunkCoord]*block.Chunk, screen *ebiten.Image, cameraX, cameraY float64) {
	if !isInitialized {
		initTileImages()
	}

	screen.Fill(color.RGBA{135, 206, 250, 255})

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	tileSize := settings.TileSize
	chunkWidth := settings.ChunkWidth
	chunkHeight := settings.ChunkHeight

	startChunkX := int(cameraX / float64(chunkWidth*tileSize))
	endChunkX := int((cameraX+float64(screenWidth))/float64(chunkWidth*tileSize)) + 1
	startChunkY := int(cameraY / float64(chunkHeight*tileSize))
	endChunkY := int((cameraY+float64(screenHeight))/float64(chunkHeight*tileSize)) + 1

	fScreenWidth := float64(screenWidth)
	fScreenHeight := float64(screenHeight)
	fTileSize := float64(tileSize)

	var renderJobs []TileRenderJob

	for coord, chunk := range worldChunks {
		if coord.X < startChunkX || coord.X > endChunkX || coord.Y < startChunkY || coord.Y > endChunkY {
			continue
		}
		baseTileX := coord.X * chunkWidth
		baseTileY := coord.Y * chunkHeight
		for y := 0; y < chunkHeight; y++ {
			globalTileY := baseTileY + y
			py := float64(globalTileY*tileSize) - cameraY
			if py+fTileSize < 0 || py >= fScreenHeight {
				continue
			}
			for x := 0; x < chunkWidth; x++ {
				globalTileX := baseTileX + x
				px := float64(globalTileX*tileSize) - cameraX
				if px+fTileSize < 0 || px >= fScreenWidth {
					continue
				}
				blockType := (*chunk)[y][x]
				if blockType == block.Air {
					continue
				}
				tile := tileImages[blockType]
				if tile == nil {
					continue
				}
				renderJobs = append(renderJobs, TileRenderJob{
					blockType: blockType,
					px:        px,
					py:        py,
					tile:      tile,
				})
			}
		}
	}

	drawOpts := getDrawOptions()
	for _, job := range renderJobs {
		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Translate(job.px, job.py)
		screen.DrawImage(job.tile, drawOpts)
	}
}

// DrawWithCameraCountBlocks draws the world and counts rendered blocks
func DrawWithCameraCountBlocks(worldChunks map[chunks.ChunkCoord]*block.Chunk, screen *ebiten.Image, cameraX, cameraY float64, blockCounter *int) {
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	tileSize := settings.TileSize
	chunkWidth := settings.ChunkWidth
	chunkHeight := settings.ChunkHeight

	startChunkX := int(cameraX / float64(chunkWidth*tileSize))
	endChunkX := int((cameraX+float64(screenWidth))/float64(chunkWidth*tileSize)) + 1
	startChunkY := int(cameraY / float64(chunkHeight*tileSize))
	endChunkY := int((cameraY+float64(screenHeight))/float64(chunkHeight*tileSize)) + 1

	fScreenWidth := float64(screenWidth)
	fScreenHeight := float64(screenHeight)
	fTileSize := float64(tileSize)

	drawOpts := getDrawOptions()
	for coord, chunk := range worldChunks {
		if coord.X < startChunkX || coord.X > endChunkX || coord.Y < startChunkY || coord.Y > endChunkY {
			continue
		}
		baseTileX := coord.X * chunkWidth
		baseTileY := coord.Y * chunkHeight
		for y := 0; y < chunkHeight; y++ {
			globalTileY := baseTileY + y
			py := float64(globalTileY*tileSize) - cameraY
			if py+fTileSize < 0 || py >= fScreenHeight {
				continue
			}
			for x := 0; x < chunkWidth; x++ {
				globalTileX := baseTileX + x
				px := float64(globalTileX*tileSize) - cameraX
				if px+fTileSize < 0 || px >= fScreenWidth {
					continue
				}
				blockType := (*chunk)[y][x]
				if blockType == block.Air {
					continue
				}
				tile := tileImages[blockType]
				if tile == nil {
					continue
				}
				drawOpts.GeoM.Reset()
				drawOpts.GeoM.Translate(px, py)
				screen.DrawImage(tile, drawOpts)
				(*blockCounter)++
			}
		}
	}
}
