package game

import (
	"fmt"
	"image/color"
	"math"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/KdntNinja/webcraft/internal/core/engine/util"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
)

// drawDebugInfo draws a Minecraft-like F3 debug screen overlay
func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	// Colors
	bgColor := color.RGBA{20, 20, 30, 230}
	borderColor := color.RGBA{80, 180, 255, 255}
	accentColor := color.RGBA{80, 255, 180, 255}
	barBgColor := color.RGBA{40, 40, 60, 200}
	fpsBarColor := color.RGBA{80, 255, 80, 255}
	chunkBarColor := color.RGBA{255, 200, 80, 255}
	entityBarColor := color.RGBA{255, 80, 180, 255}
	memBarColor := color.RGBA{120, 120, 255, 255}

	// Panel size (taller and thinner for more vertical info)
	w, h := 280, 490 // was 520, 490
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

	// Draw FPS graph (last 120 frames, larger, relative scaling to runtime min/max)
	graphX, graphY := 30, 40
	graphW, graphH := 220, 60
	fpsHistory := g.GetFPSHistory(120)
	if len(g.fpsHistory) > 0 {
		// Track min/max FPS over all time
		if g.fpsHistoryMin == 0 || g.fpsHistoryMin > g.fpsHistory[0] {
			g.fpsHistoryMin = g.fpsHistory[0]
		}
		if g.fpsHistoryMax < g.fpsHistory[0] {
			g.fpsHistoryMax = g.fpsHistory[0]
		}
		for _, v := range g.fpsHistory {
			if v < g.fpsHistoryMin {
				g.fpsHistoryMin = v
			}
			if v > g.fpsHistoryMax {
				g.fpsHistoryMax = v
			}
		}
	}
	minFPS := g.fpsHistoryMin
	maxFPS := g.fpsHistoryMax
	if minFPS == maxFPS {
		minFPS = 0
		maxFPS = 120
	}
	graph := ebiten.NewImage(graphW, graphH)
	graph.Fill(barBgColor)
	for i := 1; i < len(fpsHistory); i++ {
		x1 := float64(i-1) * float64(graphW) / float64(len(fpsHistory)-1)
		y1 := float64(graphH) - ((fpsHistory[i-1]-minFPS)/(maxFPS-minFPS))*float64(graphH)
		x2 := float64(i) * float64(graphW) / float64(len(fpsHistory)-1)
		y2 := float64(graphH) - ((fpsHistory[i]-minFPS)/(maxFPS-minFPS))*float64(graphH)
		ebitenutil.DrawLine(graph, x1, y1, x2, y2, fpsBarColor)
	}
	// Draw max/min lines
	ebitenutil.DrawLine(graph, 0, 0, float64(graphW), 0, accentColor)
	ebitenutil.DrawLine(graph, 0, float64(graphH-1), float64(graphW), float64(graphH-1), accentColor)
	if len(fpsHistory) > 0 {
		ebitenutil.DebugPrintAt(graph, fmt.Sprintf("max: %.0f", maxFPS), 2, 2)
		ebitenutil.DebugPrintAt(graph, fmt.Sprintf("min: %.0f", minFPS), 2, graphH-14)
	}
	opGraph := &ebiten.DrawImageOptions{}
	opGraph.GeoM.Translate(float64(graphX), float64(graphY))
	screen.DrawImage(graph, opGraph)
	fpsLabel := fmt.Sprintf("FPS: %.1f", g.currentFPS)
	ebitenutil.DebugPrintAt(screen, fpsLabel, graphX+2, graphY-18)

	// Draw loaded chunks bar (larger)
	barX, barY := 30, 120
	barW, barH := 220, 18
	maxChunks := 128.0
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
	ebitenutil.DebugPrintAt(screen, barLabel, barX+2, barY-18)

	// Draw entity count bar
	entityX, entityY := 30, 160
	entityW, entityH := 220, 18
	maxEntities := 64.0
	entityCount := 0.0
	if g.World != nil {
		entityCount = float64(len(g.World.Entities))
	}
	entityBar := ebiten.NewImage(entityW, entityH)
	entityBar.Fill(barBgColor)
	entityFillW := int(math.Min(float64(entityW), (entityCount/maxEntities)*float64(entityW)))
	if entityFillW > 0 {
		fill := ebiten.NewImage(entityFillW, entityH)
		fill.Fill(entityBarColor)
		entityBar.DrawImage(fill, &ebiten.DrawImageOptions{})
	}
	for i := 0; i < 2; i++ {
		border := ebiten.NewImage(entityW-2*i, entityH-2*i)
		border.Fill(borderColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i), float64(i))
		entityBar.DrawImage(border, op)
	}
	opEntity := &ebiten.DrawImageOptions{}
	opEntity.GeoM.Translate(float64(entityX), float64(entityY))
	screen.DrawImage(entityBar, opEntity)
	entityLabel := fmt.Sprintf("Entities: %d", int(entityCount))
	ebitenutil.DebugPrintAt(screen, entityLabel, entityX+2, entityY-18)

	// Draw memory usage bar (real stats)
	memX, memY := 30, 200
	memW, memH := 220, 18
	maxMem := float64(runtime.MemStats{}.Sys) / 1024.0 / 1024.0 // System memory in MB
	memUsage := util.GetMemoryUsageMB()
	if memUsage > maxMem {
		maxMem = memUsage * 1.2 // auto-scale if needed
	}
	memBar := ebiten.NewImage(memW, memH)
	memBar.Fill(barBgColor)
	memFillW := int(math.Min(float64(memW), (memUsage/maxMem)*float64(memW)))
	if memFillW > 0 {
		fill := ebiten.NewImage(memFillW, memH)
		fill.Fill(memBarColor)
		memBar.DrawImage(fill, &ebiten.DrawImageOptions{})
	}
	for i := 0; i < 2; i++ {
		border := ebiten.NewImage(memW-2*i, memH-2*i)
		border.Fill(borderColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i), float64(i))
		memBar.DrawImage(border, op)
	}
	opMem := &ebiten.DrawImageOptions{}
	opMem.GeoM.Translate(float64(memX), float64(memY))
	screen.DrawImage(memBar, opMem)
	memLabel := fmt.Sprintf("Memory: %.0fMB / %.0fMB", memUsage, maxMem)
	// Draw memory label inside the bar, clipped
	labelX := int(memX) + 8
	labelY := int(memY) + (memH / 2) - 8
	if labelX+120 > int(memX)+memW {
		labelX = int(memX) + memW - 120
	}
	ebitenutil.DebugPrintAt(screen, memLabel, labelX, labelY)

	// Player info (more detail)
	if len(g.World.Entities) > 0 {
		if p, ok := g.World.Entities[0].(*player.Player); ok {
			// Move player info to the left side, below the last stat bar/graph
			playerInfo := fmt.Sprintf("Player: X=%.2f Y=%.2f\nVX=%.2f VY=%.2f", p.X, p.Y, p.VX, p.VY)
			// Place at (30, 338) (below GC bar, with padding)
			ebitenutil.DebugPrintAt(screen, playerInfo, 30, 338)
			chunkX := int(p.X) / (settings.ChunkWidth * settings.TileSize)
			chunkY := int(p.Y) / (settings.ChunkHeight * settings.TileSize)
			chunkInfo := fmt.Sprintf("Chunk: %d, %d", chunkX, chunkY)
			// Place chunk info below player info
			ebitenutil.DebugPrintAt(screen, chunkInfo, 30, 370)
			playerStats := fmt.Sprintf("OnGround: %v", p.OnGround)
			// Place player stats below chunk info
			ebitenutil.DebugPrintAt(screen, playerStats, 30, 390)
		}
	}
	// Camera info
	camInfo := fmt.Sprintf("Camera: X=%.2f Y=%.2f", g.CameraX, g.CameraY)
	// Move camera info below player stats
	ebitenutil.DebugPrintAt(screen, camInfo, 30, 410)
	// Seed
	seedInfo := fmt.Sprintf("Seed: %d", g.Seed)
	// Move seed info below camera info
	ebitenutil.DebugPrintAt(screen, seedInfo, 30, 430)
	// World info (show loaded chunk count only)
	if g.World != nil && g.World.ChunkManager != nil {
		worldInfo := fmt.Sprintf("Loaded Chunks: %d", g.World.ChunkManager.GetLoadedChunkCount())
		// Move world info below seed info
		ebitenutil.DebugPrintAt(screen, worldInfo, 30, 450)
	}

	// TODO: Add more graphs (tick time, GC, etc.)

	// Draw tick time graph (real stats)
	tickGraphX, tickGraphY := 30, 260
	tickGraphW, tickGraphH := 220, 40
	tickTimes := util.GetTickTimes(120)
	minTick, maxTick := 999.0, 0.0
	for _, v := range tickTimes {
		if v < minTick {
			minTick = v
		}
		if v > maxTick {
			maxTick = v
		}
	}
	tickGraph := ebiten.NewImage(tickGraphW, tickGraphH)
	tickGraph.Fill(barBgColor)
	for i := 1; i < len(tickTimes); i++ {
		x1 := float64(i-1) * float64(tickGraphW) / float64(len(tickTimes)-1)
		y1 := float64(tickGraphH) - ((tickTimes[i-1]-minTick)/(maxTick-minTick))*float64(tickGraphH)
		x2 := float64(i) * float64(tickGraphW) / float64(len(tickTimes)-1)
		y2 := float64(tickGraphH) - ((tickTimes[i]-minTick)/(maxTick-minTick))*float64(tickGraphH)
		ebitenutil.DrawLine(tickGraph, x1, y1, x2, y2, memBarColor)
	}
	ebitenutil.DrawLine(tickGraph, 0, 0, float64(tickGraphW), 0, accentColor)
	ebitenutil.DrawLine(tickGraph, 0, float64(tickGraphH-1), float64(tickGraphW), float64(tickGraphH-1), accentColor)
	ebitenutil.DebugPrintAt(tickGraph, fmt.Sprintf("max: %.1fms", maxTick), 2, 2)
	ebitenutil.DebugPrintAt(tickGraph, fmt.Sprintf("min: %.1fms", minTick), 2, tickGraphH-14)
	tickOp := &ebiten.DrawImageOptions{}
	tickOp.GeoM.Translate(float64(tickGraphX), float64(tickGraphY))
	screen.DrawImage(tickGraph, tickOp)
	ebitenutil.DebugPrintAt(screen, "Tick Time (ms)", tickGraphX+2, tickGraphY-16)

	// Draw GC bar (real stats)
	gcX, gcY := 30, 310
	gcW, gcH := 220, 18
	gcPercent := util.GetGCPercent()
	gcBar := ebiten.NewImage(gcW, gcH)
	gcBar.Fill(barBgColor)
	gcFillW := int(math.Min(float64(gcW), (gcPercent/100.0)*float64(gcW)))
	if gcFillW > 0 {
		fill := ebiten.NewImage(gcFillW, gcH)
		fill.Fill(accentColor)
		gcBar.DrawImage(fill, &ebiten.DrawImageOptions{})
	}
	for i := 0; i < 2; i++ {
		border := ebiten.NewImage(gcW-2*i, gcH-2*i)
		border.Fill(borderColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i), float64(i))
		gcBar.DrawImage(border, op)
	}
	gcOp := &ebiten.DrawImageOptions{}
	gcOp.GeoM.Translate(float64(gcX), float64(gcY))
	screen.DrawImage(gcBar, gcOp)
	gcLabel := fmt.Sprintf("GC: %.1f%%", gcPercent)
	labelX = int(gcX) + 8
	labelY = int(gcY) + (gcH / 2) - 8
	if labelX+80 > int(gcX)+gcW {
		labelX = int(gcX) + gcW - 80
	}
	ebitenutil.DebugPrintAt(screen, gcLabel, labelX, labelY)
}

// GetFPSHistory returns a slice of the last n FPS values (implement this in your Game struct)
// Example stub:
// func (g *Game) GetFPSHistory(n int) []float64 {
// 	return make([]float64, n) // Replace with actual FPS history
// }
