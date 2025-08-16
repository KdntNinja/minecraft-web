package debug

import (
	"fmt"
	"image/color"
	"time"
)

// UpdateDebugUI updates all debug metrics with comprehensive information
func UpdateDebugUI(currentFPS float64) {
	if debugUI == nil {
		return
	}

	// Throttle updates to avoid UI spam (update every 100ms)
	now := time.Now()
	if now.Sub(debugUI.lastUpdateTime) < 100*time.Millisecond {
		return
	}
	debugUI.lastUpdateTime = now

	// Update FPS with color coding
	var fpsColor color.Color = color.NRGBA{100, 255, 100, 255} // Green for good FPS
	if currentFPS < 30 {
		fpsColor = color.NRGBA{255, 100, 100, 255} // Red for poor FPS
	} else if currentFPS < 50 {
		fpsColor = color.NRGBA{255, 200, 100, 255} // Orange for mediocre FPS
	}

	debugUI.fpsLabel.Label = fmt.Sprintf("FPS: %.1f", currentFPS)
	debugUI.fpsLabel.Color = fpsColor

	// Update FPS history
	debugUI.fpsHistory = append(debugUI.fpsHistory, currentFPS)
	if len(debugUI.fpsHistory) > 120 {
		debugUI.fpsHistory = debugUI.fpsHistory[1:]
	}

	// Calculate FPS min/max
	if len(debugUI.fpsHistory) > 0 {
		minFPS, maxFPS := debugUI.fpsHistory[0], debugUI.fpsHistory[0]
		for _, fps := range debugUI.fpsHistory {
			if fps < minFPS {
				minFPS = fps
			}
			if fps > maxFPS {
				maxFPS = fps
			}
		}
		debugUI.fpsMinMaxLabel.Label = fmt.Sprintf("Min/Max: %.1f/%.1f", minFPS, maxFPS)
	}
}

// UpdateDebugMetrics updates all the debug information with detailed metrics
func UpdateDebugMetrics(fpsHistory []float64, fpsMin, fpsMax, currentFPS float64,
	loadedChunks, entityCount int, memUsage, maxMem float64,
	playerInfo, chunkInfo, playerStats, camInfo, seedInfo, worldInfo string,
	tickTimes []float64, tickTimeMin, tickTimeMax, gcPercent float64,
	renderedBlocksHistory, generatedBlocksHistory []int) {

	if debugUI == nil {
		return
	}

	// Update FPS
	UpdateDebugUI(currentFPS)

	// Update tick time metrics
	if len(tickTimes) > 0 {
		avgTickTime := 0.0
		for _, t := range tickTimes {
			avgTickTime += t
		}
		avgTickTime /= float64(len(tickTimes))

		var tickColor color.Color = color.NRGBA{100, 255, 100, 255}
		if avgTickTime > 16.0 { // More than 16ms is concerning for 60fps
			tickColor = color.NRGBA{255, 100, 100, 255}
		} else if avgTickTime > 10.0 {
			tickColor = color.NRGBA{255, 200, 100, 255}
		}

		debugUI.tickTimeLabel.Label = fmt.Sprintf("Tick: %.1fms (%.1f-%.1f)", avgTickTime, tickTimeMin, tickTimeMax)
		debugUI.tickTimeLabel.Color = tickColor

		// Update tick history
		debugUI.tickHistory = append(debugUI.tickHistory, avgTickTime)
		if len(debugUI.tickHistory) > 120 {
			debugUI.tickHistory = debugUI.tickHistory[1:]
		}
	}

	// Update memory metrics with color coding
	memPercent := (memUsage / maxMem) * 100
	var memColor color.Color = color.NRGBA{100, 255, 100, 255}
	if memPercent > 80 {
		memColor = color.NRGBA{255, 100, 100, 255}
	} else if memPercent > 60 {
		memColor = color.NRGBA{255, 200, 100, 255}
	}

	debugUI.memoryLabel.Label = fmt.Sprintf("Memory: %.1fMB / %.1fMB (%.1f%%)",
		memUsage, maxMem, memPercent)
	debugUI.memoryLabel.Color = memColor

	// Update memory history
	debugUI.memHistory = append(debugUI.memHistory, memUsage)
	if len(debugUI.memHistory) > 120 {
		debugUI.memHistory = debugUI.memHistory[1:]
	}

	// Update GC metrics
	var gcColor color.Color = color.NRGBA{100, 255, 100, 255}
	if gcPercent > 5.0 {
		gcColor = color.NRGBA{255, 100, 100, 255}
	} else if gcPercent > 2.0 {
		gcColor = color.NRGBA{255, 200, 100, 255}
	}
	debugUI.gcLabel.Label = fmt.Sprintf("GC: %.2f%%", gcPercent)
	debugUI.gcLabel.Color = gcColor

	// Update world info
	debugUI.playerLabel.Label = playerInfo
	debugUI.chunkLabel.Label = fmt.Sprintf("Chunks: %d loaded", loadedChunks)
	debugUI.cameraLabel.Label = camInfo
	debugUI.seedLabel.Label = seedInfo
	debugUI.entityLabel.Label = fmt.Sprintf("Entities: %d", entityCount)
}
