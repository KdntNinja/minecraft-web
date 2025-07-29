package world

import (
	"sync"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/settings"
)

func (w *World) generateGridDataAsync(allChunks map[coretypes.ChunkCoord]*coretypes.Chunk) interface{} {
	if len(allChunks) == 0 {
		return map[string]interface{}{
			"grid":    [][]int{},
			"offsetX": 0,
			"offsetY": 0,
		}
	}

	// Calculate bounds
	minX, maxX, minY, maxY := 0, 0, 0, 0
	first := true
	for coord := range allChunks {
		if first {
			minX, maxX, minY, maxY = coord.X, coord.X, coord.Y, coord.Y
			first = false
		} else {
			if coord.X < minX {
				minX = coord.X
			}
			if coord.X > maxX {
				maxX = coord.X
			}
			if coord.Y < minY {
				minY = coord.Y
			}
			if coord.Y > maxY {
				maxY = coord.Y
			}
		}
	}

	width := (maxX - minX + 1) * settings.ChunkWidth
	height := (maxY - minY + 1) * settings.ChunkHeight

	// Get grid from pool or create new
	var grid [][]int
	if poolGrid := w.gridGenerationPool.Get(); poolGrid != nil {
		if pGrid, ok := poolGrid.([][]int); ok {
			grid = pGrid[:0] // Reset slice but keep capacity
		}
	}
	if grid == nil {
		grid = make([][]int, height)
	}

	// Ensure grid has correct dimensions
	for len(grid) < height {
		grid = append(grid, make([]int, width))
	}
	for i := 0; i < height; i++ {
		if len(grid[i]) < width {
			grid[i] = make([]int, width)
		}
	}

	// Fill grid in parallel
	var wg sync.WaitGroup
	for coord, chunk := range allChunks {
		wg.Add(1)
		go func(coord coretypes.ChunkCoord, chunk *coretypes.Chunk) {
			defer wg.Done()
			for y := 0; y < settings.ChunkHeight; y++ {
				for x := 0; x < settings.ChunkWidth; x++ {
					globalX := (coord.X-minX)*settings.ChunkWidth + x
					globalY := (coord.Y-minY)*settings.ChunkHeight + y
					if y < len(chunk.Blocks) && x < len(chunk.Blocks[y]) {
						grid[globalY][globalX] = int(chunk.Blocks[y][x])
					} else {
						grid[globalY][globalX] = int(coretypes.Air)
					}
				}
			}
		}(coord, chunk)
	}
	wg.Wait()

	return map[string]interface{}{
		"grid":    grid,
		"offsetX": minX * settings.ChunkWidth,
		"offsetY": minY * settings.ChunkHeight,
	}
}
