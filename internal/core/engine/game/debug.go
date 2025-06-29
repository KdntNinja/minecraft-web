package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
)

// drawDebugInfo draws a Minecraft-like F3 debug screen overlay
func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	// Colors
	bgColor := color.RGBA{20, 20, 30, 220}
	borderColor := color.RGBA{80, 180, 255, 255}
	accentColor := color.RGBA{80, 255, 180, 255}
	barBgColor := color.RGBA{40, 40, 60, 200}
	fpsBarColor := color.RGBA{80, 255, 80, 255}
	chunkBarColor := color.RGBA{255, 200, 80, 255}

	// Panel size
	w, h := 360, 200
	bg := ebiten.NewImage(w, h)
	bg.Fill(bgColor)
	// Draw border
	for i := 0; i < 3; i++ {
		borderRect := ebiten.NewImage(w-2*i, h-2*i)
		borderRect.Fill(borderColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i), float64(i))
		bg.DrawImage(borderRect, op)
	}
	// Draw background inside border
	inner := ebiten.NewImage(w-6, h-6)
	inner.Fill(bgColor)
	opInner := &ebiten.DrawImageOptions{}
	opInner.GeoM.Translate(3, 3)
	bg.DrawImage(inner, opInner)
	screen.DrawImage(bg, &ebiten.DrawImageOptions{})

	// Draw FPS graph (last 60 frames)
	graphX, graphY := 20, 30
	graphW, graphH := 120, 40
	fpsHistory := g.GetFPSHistory(60) // implement this method to return []float64
	maxFPS := 120.0
	graph := ebiten.NewImage(graphW, graphH)
	graph.Fill(barBgColor)
	for i := 1; i < len(fpsHistory); i++ {
		x1 := float64(i-1) * float64(graphW) / float64(len(fpsHistory)-1)
		y1 := float64(graphH) - (fpsHistory[i-1]/maxFPS)*float64(graphH)
		x2 := float64(i) * float64(graphW) / float64(len(fpsHistory)-1)
		y2 := float64(graphH) - (fpsHistory[i]/maxFPS)*float64(graphH)
		ebitenutil.DrawLine(graph, x1, y1, x2, y2, fpsBarColor)
	}
	// Draw max line
	ebitenutil.DrawLine(graph, 0, 0, float64(graphW), 0, accentColor)
	// Draw min line
	ebitenutil.DrawLine(graph, 0, float64(graphH-1), float64(graphW), float64(graphH-1), accentColor)
	opGraph := &ebiten.DrawImageOptions{}
	opGraph.GeoM.Translate(float64(graphX), float64(graphY))
	screen.DrawImage(graph, opGraph)
	// FPS label
	fpsLabel := fmt.Sprintf("FPS: %.1f", g.currentFPS)
	ebitenutil.DebugPrintAt(screen, fpsLabel, graphX+2, graphY-16)

	// Draw loaded chunks bar
	barX, barY := 20, 80
	barW, barH := 120, 12
	maxChunks := 64.0
	loadedChunks := 0.0
	if g.World != nil && g.World.ChunkManager != nil {
		loadedChunks = float64(g.World.ChunkManager.GetLoadedChunkCount())
	}
	bar := ebiten.NewImage(barW, barH)
	bar.Fill(barBgColor)
	fillW := int(math.Min(float64(barW), (loadedChunks/maxChunks)*float64(barW)))
	if fillW > 0 {
		fill := ebiten.NewImage(fillW, barH)
		fill.Fill(chunkBarColor)
		bar.DrawImage(fill, &ebiten.DrawImageOptions{})
	}
	// Border
	for i := 0; i < 2; i++ {
		border := ebiten.NewImage(barW-2*i, barH-2*i)
		border.Fill(borderColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i), float64(i))
		bar.DrawImage(border, op)
	}
	opBar := &ebiten.DrawImageOptions{}
	opBar.GeoM.Translate(float64(barX), float64(barY))
	screen.DrawImage(bar, opBar)
	barLabel := fmt.Sprintf("Loaded Chunks: %d", int(loadedChunks))
	ebitenutil.DebugPrintAt(screen, barLabel, barX+2, barY-16)

	// Player info
	if len(g.World.Entities) > 0 {
		if p, ok := g.World.Entities[0].(*player.Player); ok {
			playerInfo := fmt.Sprintf("Player: X=%.2f Y=%.2f", p.X, p.Y)
			ebitenutil.DebugPrintAt(screen, playerInfo, 180, 40)
			chunkX := int(p.X) / (settings.ChunkWidth * settings.TileSize)
			chunkY := int(p.Y) / (settings.ChunkHeight * settings.TileSize)
			chunkInfo := fmt.Sprintf("Chunk: %d, %d", chunkX, chunkY)
			ebitenutil.DebugPrintAt(screen, chunkInfo, 180, 60)
		}
	}
	// Camera info
	camInfo := fmt.Sprintf("Camera: X=%.2f Y=%.2f", g.CameraX, g.CameraY)
	ebitenutil.DebugPrintAt(screen, camInfo, 180, 80)
	// Seed
	seedInfo := fmt.Sprintf("Seed: %d", g.Seed)
	ebitenutil.DebugPrintAt(screen, seedInfo, 180, 100)

	// TODO: Add more graphs (memory, entity count, etc.)
}

// GetFPSHistory returns a slice of the last n FPS values (implement this in your Game struct)
// Example stub:
// func (g *Game) GetFPSHistory(n int) []float64 {
// 	return make([]float64, n) // Replace with actual FPS history
// }
