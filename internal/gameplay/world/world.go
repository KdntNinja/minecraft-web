package world

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
	"github.com/KdntNinja/webcraft/internal/generation/terrain"
)

type ChunkCoord struct {
	X int
	Y int
}

type World struct {
	Chunks   map[ChunkCoord]block.Chunk // Fixed world: map of chunk coordinates to chunks
	Entities entity.Entities

	// Performance optimization caches
	cachedGrid        [][]int // Cached collision grid
	cachedGridOffsetX int     // Cached grid offset X
	cachedGridOffsetY int     // Cached grid offset Y
	gridDirty         bool    // Flag to indicate grid needs regeneration
}

// NewWorld constructs a new World instance with a fixed set of pre-generated chunks
func NewWorld(numChunksY int, centerChunkX int, seed int64) *World {
	terrain.ResetWorldGeneration(seed)
	w := &World{
		Chunks:    make(map[ChunkCoord]block.Chunk),
		Entities:  entity.Entities{},
		gridDirty: true, // Grid needs initial generation
	}
	// Generate a large fixed world area - no dynamic loading
	worldWidth := settings.WorldChunksX  // Total chunks horizontally
	worldHeight := settings.WorldChunksY // Total chunks vertically

	for cy := 0; cy < worldHeight; cy++ {
		for cx := -worldWidth / 2; cx <= worldWidth/2; cx++ {
			coord := ChunkCoord{X: centerChunkX + cx, Y: cy}
			w.Chunks[coord] = GenerateChunk(coord.X, coord.Y)
		}
	}

	// Add player entity at center - improved spawning with rightward shift
	centerChunkCol := 2 // Shift spawn point 2 chunks to the right for better world exploration
	centerBlockX := centerChunkCol*settings.ChunkWidth + settings.ChunkWidth/2
	px := float64(centerBlockX * settings.TileSize)

	// Find the surface height at the center position - try multiple points for best spawn
	bestSpawnX := centerBlockX
	bestSpawnY := 1000

	// Sample multiple positions around center to find the best spawn point
	for testX := centerBlockX - 4; testX <= centerBlockX+4; testX++ {
		testSurfaceY := FindSurfaceHeight(testX, w)
		if testSurfaceY < bestSpawnY && testSurfaceY > 10 { // Avoid spawning too high or underground
			bestSpawnY = testSurfaceY
			bestSpawnX = testX
		}
	}

	// Use the best spawn position found
	px = float64(bestSpawnX * settings.TileSize)
	spawnY := bestSpawnY - 3 // Spawn 3 blocks above surface for safety

	// Ensure spawn position is reasonable
	if spawnY < 5 {
		spawnY = 5
	}
	if spawnY > 200 {
		spawnY = 200
	}

	py := float64(spawnY * settings.TileSize)
	w.Entities = append(w.Entities, player.NewPlayer(px, py))
	return w
}

// ToIntGrid flattens the world's blocks into a [][]int grid for entity collision, and returns the offset (minX, minY)
// Uses caching to avoid regenerating the grid every frame
func (w *World) ToIntGrid() ([][]int, int, int) {
	if len(w.Chunks) == 0 {
		return [][]int{}, 0, 0
	}

	// Return cached grid if it's still valid
	if !w.gridDirty && w.cachedGrid != nil {
		return w.cachedGrid, w.cachedGridOffsetX, w.cachedGridOffsetY
	}

	// Regenerate grid only when necessary
	minX, maxX, minY, maxY := 0, 0, 0, 0
	first := true
	for coord := range w.Chunks {
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

	// Reuse cached grid if dimensions match, otherwise allocate new
	if w.cachedGrid == nil || len(w.cachedGrid) != height || (height > 0 && len(w.cachedGrid[0]) != width) {
		w.cachedGrid = make([][]int, height)
		for y := 0; y < height; y++ {
			w.cachedGrid[y] = make([]int, width)
		}
	}

	grid := w.cachedGrid
	for y := 0; y < height; y++ {
		cy := minY + y/settings.ChunkHeight
		inChunkY := y % settings.ChunkHeight
		for x := 0; x < width; x++ {
			cx := minX + x/settings.ChunkWidth
			inChunkX := x % settings.ChunkWidth
			coord := ChunkCoord{X: cx, Y: cy}
			chunk, ok := w.Chunks[coord]
			if !ok || len(chunk) == 0 {
				grid[y][x] = int(block.Air)
				continue
			}
			grid[y][x] = int(chunk[inChunkY][inChunkX])
		}
	}

	// Cache the results
	w.cachedGridOffsetX = minX * settings.ChunkWidth
	w.cachedGridOffsetY = minY * settings.ChunkHeight
	w.gridDirty = false

	return grid, w.cachedGridOffsetX, w.cachedGridOffsetY
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
