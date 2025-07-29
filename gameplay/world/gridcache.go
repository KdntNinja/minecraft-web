package world

import (
	"sync"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/settings"
)

// ToIntGrid flattens the world's blocks into a [][]int grid for entity collision, and returns the offset (minX, minY)
// Uses caching to avoid regenerating the grid every frame
func (w *World) ToIntGrid() ([][]int, int, int) {
	allChunks := w.ChunkManager.GetAllChunks()
	if len(allChunks) == 0 {
		return [][]int{}, 0, 0
	}

	// Regenerate grid only when necessary
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

	// Always allocate a new grid to avoid stale data when moving into negative coords
	grid := make([][]int, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
	}

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

	w.cachedGrid = grid
	w.cachedGridOffsetX = minX * settings.ChunkWidth
	w.cachedGridOffsetY = minY * settings.ChunkHeight
	w.gridDirty = false

	return grid, w.cachedGridOffsetX, w.cachedGridOffsetY
}

// IsGridDirty returns whether the collision grid needs to be regenerated
func (w *World) IsGridDirty() bool {
	return w.gridDirty
}

// GetCachedGrid returns the cached collision grid without regenerating it
func (w *World) GetCachedGrid() ([][]int, int, int) {
	if w.cachedGrid == nil {
		return [][]int{}, 0, 0
	}
	return w.cachedGrid, w.cachedGridOffsetX, w.cachedGridOffsetY
}
