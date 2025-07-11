package rendering

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// DrawDebugOverlay draws a Minecraft-like F3 debug screen overlay
func DrawDebugOverlay(
	screen *ebiten.Image,
	fpsHistory []float64,
	minFPS, maxFPS float64,
	currentFPS float64,
	loadedChunks int,
	entityCount int,
	memUsage, maxMem float64,
	playerInfo, chunkInfo, playerStats, camInfo, seedInfo, worldInfo string,
	tickTimes []float64, minTick, maxTick float64,
	gcPercent float64,
	renderedBlocksHistory []int,
	generatedBlocksHistory []int,
) {
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
	w, h := 280, 550
	// Use a package-level variable to cache the background image
	var debugOverlayBg *ebiten.Image
	if debugOverlayBg == nil || debugOverlayBg.Bounds().Dx() != w || debugOverlayBg.Bounds().Dy() != h {
		bg := ebiten.NewImage(w, h)
		bg.Fill(bgColor)
		for i := 0; i < 3; i++ {
			borderRect := ebiten.NewImage(w-2*i, h-2*i)
			borderRect.Fill(borderColor)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i), float64(i))
			bg.DrawImage(borderRect, op)
		}
		inner := ebiten.NewImage(w-6, h-6)
		inner.Fill(bgColor)
		opInner := &ebiten.DrawImageOptions{}
		opInner.GeoM.Translate(3, 3)
		bg.DrawImage(inner, opInner)
		debugOverlayBg = bg
	}
	screen.DrawImage(debugOverlayBg, &ebiten.DrawImageOptions{})

	// Draw FPS graph (last 120 frames, larger, relative scaling to runtime min/max)
	graphX, graphY := 30, 40
	graphW, graphH := 220, 60
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
	// FPS label
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.1f", currentFPS), 30, 20)

	// Draw loaded chunks bar (larger)
	barX, barY := 30, 120
	barW, barH := 220, 18
	maxChunks := 128.0
	bar := ebiten.NewImage(barW, barH)
	bar.Fill(barBgColor)
	fillW := int(math.Min(float64(barW), (float64(loadedChunks)/maxChunks)*float64(barW)))
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
	barLabel := fmt.Sprintf("Loaded Chunks: %d", loadedChunks)
	ebitenutil.DebugPrintAt(screen, barLabel, barX+2, barY-18)

	// Draw entity count bar
	entityX, entityY := 30, 160
	entityW, entityH := 220, 18
	maxEntities := 64.0
	entityBar := ebiten.NewImage(entityW, entityH)
	entityBar.Fill(barBgColor)
	entityFillW := int(math.Min(float64(entityW), (float64(entityCount)/maxEntities)*float64(entityW)))
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
	entityLabel := fmt.Sprintf("Entities: %d", entityCount)
	ebitenutil.DebugPrintAt(screen, entityLabel, entityX+2, entityY-18)

	// Draw memory usage bar (real stats)
	memX, memY := 30, 200
	memW, memH := 220, 18
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
	playerY := 338
	lineSpacing := 28
	ebitenutil.DebugPrintAt(screen, playerInfo, 30, playerY)
	ebitenutil.DebugPrintAt(screen, chunkInfo, 30, playerY+lineSpacing*2)
	ebitenutil.DebugPrintAt(screen, playerStats, 30, playerY+lineSpacing*3)
	// Camera info
	camY := 338 + 28*4
	ebitenutil.DebugPrintAt(screen, camInfo, 30, camY)
	// Seed
	seedY := camY + 28
	ebitenutil.DebugPrintAt(screen, seedInfo, 30, seedY)
	// World info
	worldY := seedY + 28
	ebitenutil.DebugPrintAt(screen, worldInfo, 30, worldY)

	// Draw tick time graph (real stats)
	tickGraphX, tickGraphY := 30, 260
	tickGraphW, tickGraphH := 220, 40
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

	// Draw rendered blocks chart
	chartW, chartH := 220, 40
	chartX, chartY := 30, 370
	chart := ebiten.NewImage(chartW, chartH)
	chart.Fill(color.RGBA{30, 30, 30, 200})
	if len(renderedBlocksHistory) > 1 {
		maxBlocks := 1
		for _, v := range renderedBlocksHistory {
			if v > maxBlocks {
				maxBlocks = v
			}
		}
		for i := 1; i < len(renderedBlocksHistory); i++ {
			prev := renderedBlocksHistory[i-1]
			curr := renderedBlocksHistory[i]
			x0 := (i - 1) * chartW / len(renderedBlocksHistory)
			y0 := chartH - (prev * chartH / maxBlocks)
			x1 := i * chartW / len(renderedBlocksHistory)
			y1 := chartH - (curr * chartH / maxBlocks)
			col := color.RGBA{80, 180, 255, 255}
			for dx := 0; dx < 2; dx++ {
				if x0+dx < chartW && x1+dx < chartW {
					for dy := y0; dy <= y1 && dy < chartH; dy++ {
						chart.Set(x0+dx, dy, col)
					}
				}
			}
		}
	}
	chartOp := &ebiten.DrawImageOptions{}
	chartOp.GeoM.Translate(float64(chartX), float64(chartY))
	screen.DrawImage(chart, chartOp)
	// Label
	ebitenutil.DebugPrintAt(screen, "Rendered Blocks", chartX, chartY-16)

	// Draw generated blocks chart
	chart2 := ebiten.NewImage(chartW, chartH)
	chart2.Fill(color.RGBA{30, 30, 30, 200})
	if len(generatedBlocksHistory) > 1 {
		maxBlocks := 1
		for _, v := range generatedBlocksHistory {
			if v > maxBlocks {
				maxBlocks = v
			}
		}
		for i := 1; i < len(generatedBlocksHistory); i++ {
			prev := generatedBlocksHistory[i-1]
			curr := generatedBlocksHistory[i]
			x0 := (i - 1) * chartW / len(generatedBlocksHistory)
			y0 := chartH - (prev * chartH / maxBlocks)
			x1 := i * chartW / len(generatedBlocksHistory)
			y1 := chartH - (curr * chartH / maxBlocks)
			col := color.RGBA{255, 180, 80, 255}
			for dx := 0; dx < 2; dx++ {
				if x0+dx < chartW && x1+dx < chartW {
					for dy := y0; dy <= y1 && dy < chartH; dy++ {
						chart2.Set(x0+dx, dy, col)
					}
				}
			}
		}
	}
	chart2Op := &ebiten.DrawImageOptions{}
	chart2Op.GeoM.Translate(float64(chartX), float64(chartY+chartH+10))
	screen.DrawImage(chart2, chart2Op)
	ebitenutil.DebugPrintAt(screen, "Generated Blocks", chartX, chartY+chartH-6)
}
