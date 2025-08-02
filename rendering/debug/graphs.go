package debug

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// drawEmbeddedGraphs draws performance graphs embedded within the UI containers
func drawEmbeddedGraphs(screen *ebiten.Image, graphType string) {
	if debugUI == nil {
		return
	}

	screenWidth := screen.Bounds().Dx()

	// Calculate positions based on UI container locations
	// Since containers are positioned at top-right, we need to calculate their screen positions
	baseX := screenWidth - 320 + 20 // Panel width - padding
	baseY := 480                    // Move graphs up (was 480)

	// Draw embedded graphs with smaller size to fit in containers
	graphWidth := 280
	graphHeight := 60
	graphSpacing := 90 // Spacing between embedded graphs

	// Draw only the requested graph type
	switch graphType {
	case "fps":
		if len(debugUI.fpsHistory) > 1 {
			drawEmbeddedGraph(screen, debugUI.fpsHistory, baseX, baseY, graphWidth, graphHeight,
				"", color.NRGBA{100, 255, 100, 255}, color.NRGBA{0, 0, 0, 0})
		}
	case "tick":
		if len(debugUI.tickHistory) > 1 {
			drawEmbeddedGraph(screen, debugUI.tickHistory, baseX, baseY+graphSpacing,
				graphWidth, graphHeight, "",
				color.NRGBA{255, 200, 100, 255}, color.NRGBA{0, 0, 0, 0})
		}
	case "mem":
		if len(debugUI.memHistory) > 1 {
			drawEmbeddedGraph(screen, debugUI.memHistory, baseX, baseY+2*graphSpacing,
				graphWidth, graphHeight, "",
				color.NRGBA{100, 200, 255, 255}, color.NRGBA{0, 0, 0, 0})
		}
	}
}

// drawEmbeddedGraph draws a single performance graph within a UI container
func drawEmbeddedGraph(screen *ebiten.Image, data []float64, x, y, width, height int,
	title string, lineColor, bgColor color.Color) {

	if len(data) < 2 {
		return
	}

	// Calculate data range
	minVal, maxVal := data[0], data[0]
	for _, val := range data {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	// Ensure reasonable range
	if maxVal-minVal < 1 {
		maxVal = minVal + 1
	}

	// Draw data line and info inside the graph outline
	// Use minimal padding (2px) to match the outline
	graphX := x + 2
	graphY := y + 2
	graphW := width - 4
	graphH := height - 4
	stepX := float64(graphW) / float64(len(data)-1)
	lineImg := ebiten.NewImage(2, 2)
	lineImg.Fill(lineColor)

	for i := 1; i < len(data); i++ {
		x1 := float64(graphX) + float64(i-1)*stepX
		y1 := float64(graphY+graphH) - (data[i-1]-minVal)/(maxVal-minVal)*float64(graphH)
		x2 := float64(graphX) + float64(i)*stepX
		y2 := float64(graphY+graphH) - (data[i]-minVal)/(maxVal-minVal)*float64(graphH)

		// Draw line segments with interpolation
		steps := math.Max(math.Abs(x2-x1), math.Abs(y2-y1))
		if steps < 1 {
			steps = 1
		}
		for step := 0.0; step <= steps; step += 1.0 {
			t := step / steps
			px := x1 + t*(x2-x1)
			py := y1 + t*(y2-y1)

			pointOp := &ebiten.DrawImageOptions{}
			pointOp.GeoM.Translate(px, py)
			screen.DrawImage(lineImg, pointOp)
		}
	}

	// Draw current value in bottom right of graph area
	if debugUI.smallFont != nil {
		currentVal := data[len(data)-1]
		valueText := fmt.Sprintf("%.1f", currentVal)
		valueOp := &text.DrawOptions{}
		valueOp.GeoM.Translate(float64(x+width-55), float64(y+height-22))
		valueOp.ColorScale.ScaleWithColor(lineColor)
		text.Draw(screen, valueText, debugUI.smallFont, valueOp)
	}
}

// drawGraph draws a single performance graph with labels and grid
func drawGraph(screen *ebiten.Image, data []float64, x, y, width, height int,
	title string, lineColor, bgColor color.Color) {

	if len(data) < 2 {
		return
	}

	// Create graph background with border
	graphBg := ebiten.NewImage(width, height)
	graphBg.Fill(bgColor)

	// Draw border
	borderColor := color.NRGBA{80, 140, 200, 150}
	for i := 0; i < 2; i++ {
		// Top and bottom borders
		for px := i; px < width-i; px++ {
			graphBg.Set(px, i, borderColor)
			graphBg.Set(px, height-1-i, borderColor)
		}
		// Left and right borders
		for py := i; py < height-i; py++ {
			graphBg.Set(i, py, borderColor)
			graphBg.Set(width-1-i, py, borderColor)
		}
	}

	// Draw background to screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(graphBg, op)

	// Calculate data range
	minVal, maxVal := data[0], data[0]
	for _, val := range data {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	// Ensure reasonable range
	if maxVal-minVal < 1 {
		maxVal = minVal + 1
	}

	// Draw grid lines
	gridColor := color.NRGBA{60, 120, 180, 80}
	gridImg := ebiten.NewImage(1, 1)
	gridImg.Fill(gridColor)

	// Horizontal grid lines
	for i := 1; i < 4; i++ {
		gridY := float64(y) + float64(height*i/4)
		gridOp := &ebiten.DrawImageOptions{}
		gridOp.GeoM.Scale(float64(width), 1)
		gridOp.GeoM.Translate(float64(x), gridY)
		screen.DrawImage(gridImg, gridOp)
	}

	// Draw data line with anti-aliasing effect
	stepX := float64(width-4) / float64(len(data)-1)
	lineImg := ebiten.NewImage(2, 2)
	lineImg.Fill(lineColor)

	for i := 1; i < len(data); i++ {
		x1 := float64(x+2) + float64(i-1)*stepX
		y1 := float64(y+height-2) - (data[i-1]-minVal)/(maxVal-minVal)*float64(height-4)
		x2 := float64(x+2) + float64(i)*stepX
		y2 := float64(y+height-2) - (data[i]-minVal)/(maxVal-minVal)*float64(height-4)

		// Draw line segments
		steps := math.Max(math.Abs(x2-x1), math.Abs(y2-y1))
		for step := 0.0; step <= steps; step += 1.0 {
			t := step / steps
			px := x1 + t*(x2-x1)
			py := y1 + t*(y2-y1)

			pointOp := &ebiten.DrawImageOptions{}
			pointOp.GeoM.Translate(px, py)
			screen.DrawImage(lineImg, pointOp)
		}
	}

	// Draw title and value labels
	if debugUI.smallFont != nil {
		// Title
		titleOp := &text.DrawOptions{}
		titleOp.GeoM.Translate(float64(x+4), float64(y+2))
		titleOp.ColorScale.ScaleWithColor(color.NRGBA{220, 220, 220, 255})
		text.Draw(screen, title, debugUI.smallFont, titleOp)

		// Current value
		currentVal := data[len(data)-1]
		valueText := fmt.Sprintf("%.1f", currentVal)
		valueOp := &text.DrawOptions{}
		valueOp.GeoM.Translate(float64(x+width-50), float64(y+2))
		valueOp.ColorScale.ScaleWithColor(lineColor)
		text.Draw(screen, valueText, debugUI.smallFont, valueOp)

		// Min/Max labels
		minText := fmt.Sprintf("%.1f", minVal)
		minOp := &text.DrawOptions{}
		minOp.GeoM.Translate(float64(x+4), float64(y+height-12))
		minOp.ColorScale.ScaleWithColor(color.NRGBA{160, 160, 160, 255})
		text.Draw(screen, minText, debugUI.smallFont, minOp)

		maxText := fmt.Sprintf("%.1f", maxVal)
		maxOp := &text.DrawOptions{}
		maxOp.GeoM.Translate(float64(x+4), float64(y+15))
		maxOp.ColorScale.ScaleWithColor(color.NRGBA{160, 160, 160, 255})
		text.Draw(screen, maxText, debugUI.smallFont, maxOp)
	}
}

// DrawPerformanceGraphs draws FPS, tick, and memory graphs in the debug panel
func DrawPerformanceGraphs(screen *ebiten.Image) {
	if debugUI == nil {
		return
	}

	// Panel position and graph layout
	screenWidth := screen.Bounds().Dx()
	baseX := screenWidth - 320 + 20 // Panel width - padding
	baseY := 340                    // Move graphs up (was 480)
	graphWidth := 280
	graphHeight := 60
	graphSpacing := 90

	// FPS Graph
	if len(debugUI.fpsHistory) > 1 {
		drawEmbeddedGraph(screen, debugUI.fpsHistory, baseX, baseY, graphWidth, graphHeight,
			"FPS", color.NRGBA{100, 255, 100, 255}, color.NRGBA{30, 30, 30, 200})
	}
	// Tick Time Graph
	if len(debugUI.tickHistory) > 1 {
		drawEmbeddedGraph(screen, debugUI.tickHistory, baseX, baseY+graphSpacing, graphWidth, graphHeight,
			"Tick Time", color.NRGBA{255, 200, 100, 255}, color.NRGBA{30, 30, 30, 200})
	}
	// Memory Graph
	if len(debugUI.memHistory) > 1 {
		drawEmbeddedGraph(screen, debugUI.memHistory, baseX, baseY+2*graphSpacing, graphWidth, graphHeight,
			"Memory", color.NRGBA{100, 200, 255, 255}, color.NRGBA{30, 30, 30, 200})
	}
}

// DrawDebugOverlay maintains compatibility with existing engine code
func DrawDebugOverlay(screen *ebiten.Image, fpsHistory []float64, fpsMin, fpsMax, currentFPS float64,
	loadedChunks, entityCount int, memUsage, maxMem float64,
	playerInfo, chunkInfo, playerStats, camInfo, seedInfo, worldInfo string,
	tickTimes []float64, tickTimeMin, tickTimeMax, gcPercent float64,
	renderedBlocksHistory, generatedBlocksHistory []int) {

	if debugUI == nil {
		return
	}

	// Update all metrics
	UpdateDebugMetrics(fpsHistory, fpsMin, fpsMax, currentFPS, loadedChunks, entityCount,
		memUsage, maxMem, playerInfo, chunkInfo, playerStats, camInfo, seedInfo, worldInfo,
		tickTimes, tickTimeMin, tickTimeMax, gcPercent, renderedBlocksHistory, generatedBlocksHistory)

	// Draw the UI
	DrawDebugUI(screen)
}
