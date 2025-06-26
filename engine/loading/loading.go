package loading

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/entity"
	"github.com/KdntNinja/webcraft/engine/player"
	"github.com/KdntNinja/webcraft/engine/world"
)

type LoadingState int

const (
	LoadingTerrain LoadingState = iota
	LoadingPlayer
	LoadingComplete
)

type LoadingScreen struct {
	State    LoadingState
	Progress float64
	Message  string
	World    *world.World
	PlayerX  float64
	PlayerY  float64
	Complete bool
}

func NewLoadingScreen() *LoadingScreen {
	return &LoadingScreen{
		State:    LoadingTerrain,
		Progress: 0.0,
		Message:  "Generating terrain...",
		Complete: false,
	}
}

// Update processes the loading stages
func (ls *LoadingScreen) Update() {
	switch ls.State {
	case LoadingTerrain:
		if ls.World == nil {
			// Start terrain generation
			ls.World = ls.generateWorldWithProgress()
			ls.Progress = 0.5
			ls.State = LoadingPlayer
			ls.Message = "Finding spawn location..."
		}

	case LoadingPlayer:
		if ls.World != nil {
			// Find proper spawn location
			ls.findPlayerSpawnLocation()
			ls.Progress = 1.0
			ls.State = LoadingComplete
			ls.Message = "Ready!"
		}

	case LoadingComplete:
		ls.Complete = true
	}
}

func (ls *LoadingScreen) generateWorldWithProgress() *world.World {
	// Generate world chunks
	centerChunkX := 0
	numChunksY := 10
	width := 5 // 2 chunks left, 1 center, 2 right

	blocks := make([][]world.Chunk, numChunksY)
	for cy := 0; cy < numChunksY; cy++ {
		blocks[cy] = make([]world.Chunk, width)
		for cx := 0; cx < width; cx++ {
			chunkX := centerChunkX + cx - 2
			blocks[cy][cx] = world.GenerateChunk(chunkX, cy)
		}
	}

	w := &world.World{
		Blocks:    blocks,
		Entities:  entity.Entities{},
		MinChunkX: centerChunkX - 2,
		MinChunkY: 0,
	}

	return w
}

func (ls *LoadingScreen) findPlayerSpawnLocation() {
	if ls.World == nil {
		return
	}

	// Calculate center position
	centerChunkCol := len(ls.World.Blocks[0]) / 2
	centerBlockX := centerChunkCol*block.ChunkWidth + block.ChunkWidth/2
	px := float64(centerBlockX * block.TileSize)

	// Convert to grid to find surface
	grid := ls.World.ToIntGrid()

	// Find the surface from top down - look for first solid block
	spawnY := 0
	found := false

	for y := 0; y < len(grid) && !found; y++ {
		if centerBlockX < len(grid[y]) {
			if entity.IsSolid(grid, centerBlockX, y) {
				// Found surface, place player exactly on top
				spawnY = y - 1 // One block above the solid surface
				found = true
			}
		}
	}

	// Fallback if no surface found
	if !found {
		spawnY = 5
	}

	// Ensure spawn position is valid and has clearance for player (2 blocks tall)
	if spawnY < 0 {
		spawnY = 0
	}

	// Make sure there's air space for the player
	for spawnY < len(grid)-2 {
		if spawnY >= 0 && centerBlockX < len(grid[spawnY]) &&
			!entity.IsSolid(grid, centerBlockX, spawnY) &&
			!entity.IsSolid(grid, centerBlockX, spawnY+1) {
			break // Found good spawn location
		}
		spawnY++
	}

	// Store spawn position
	ls.PlayerX = px
	ls.PlayerY = float64(spawnY * block.TileSize)

	// Create player at the calculated position
	ls.World.Entities = append(ls.World.Entities, player.NewPlayer(ls.PlayerX, ls.PlayerY))
}

// Draw renders the loading screen
func (ls *LoadingScreen) Draw(screen *ebiten.Image) {
	// Fill background with dark color
	screen.Fill(color.RGBA{32, 32, 64, 255})

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Draw loading bar background
	barWidth := 400
	barHeight := 20
	barX := (screenWidth - barWidth) / 2
	barY := screenHeight/2 - 10

	// Loading bar background
	barBg := ebiten.NewImage(barWidth, barHeight)
	barBg.Fill(color.RGBA{64, 64, 64, 255})

	barOp := &ebiten.DrawImageOptions{}
	barOp.GeoM.Translate(float64(barX), float64(barY))
	screen.DrawImage(barBg, barOp)

	// Loading bar progress
	progressWidth := int(float64(barWidth) * ls.Progress)
	if progressWidth > 0 {
		barFg := ebiten.NewImage(progressWidth, barHeight)
		barFg.Fill(color.RGBA{64, 255, 64, 255})

		progressOp := &ebiten.DrawImageOptions{}
		progressOp.GeoM.Translate(float64(barX), float64(barY))
		screen.DrawImage(barFg, progressOp)
	}

	// Draw loading message
	messageX := screenWidth/2 - len(ls.Message)*4 // Rough text centering
	messageY := barY - 40

	// Create a simple text rendering (since we don't have a font loaded)
	// We'll create a simple pixel-based text renderer
	ls.drawText(screen, ls.Message, messageX, messageY, color.White)

	// Draw percentage
	percentText := fmt.Sprintf("%.0f%%", ls.Progress*100)
	percentX := screenWidth/2 - len(percentText)*4
	percentY := barY + 40
	ls.drawText(screen, percentText, percentX, percentY, color.White)
}

// Simple text renderer using rectangles (since we don't have font loaded)
func (ls *LoadingScreen) drawText(screen *ebiten.Image, text string, x, y int, clr color.Color) {
	// This is a very basic text renderer - in a real game you'd use a proper font
	// For now, we'll just draw the text as simple rectangles for each character
	for i, char := range text {
		if char == ' ' {
			continue
		}

		// Draw a simple rectangle for each character
		charImg := ebiten.NewImage(8, 12)
		charImg.Fill(clr)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+i*10), float64(y))
		screen.DrawImage(charImg, op)
	}
}

func (ls *LoadingScreen) IsComplete() bool {
	return ls.Complete
}

func (ls *LoadingScreen) GetWorld() *world.World {
	return ls.World
}
