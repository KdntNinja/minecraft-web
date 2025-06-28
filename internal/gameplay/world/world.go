package world

import (
	"fmt"

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
func NewWorld(seed int64) *World {
	terrain.ResetWorldGeneration(seed)
	w := &World{
		Chunks:    make(map[ChunkCoord]block.Chunk),
		Entities:  entity.Entities{},
		gridDirty: true, // Grid needs initial generation
	}
	// Generate a large fixed world area - no dynamic loading
	worldWidth := settings.WorldChunksX  // Total chunks horizontally
	worldHeight := settings.WorldChunksY // Total chunks vertically

	fmt.Printf("DEBUG: WorldChunksX=%d, WorldChunksY=%d\n", worldWidth, worldHeight)

	// Calculate chunk range to center the world around (0,0)
	halfWidth := worldWidth / 2
	var startX, endX int
	if worldWidth%2 == 0 {
		// Even number of chunks: generate equal chunks on both sides
		// For 24 chunks: -12 to 11 (24 total)
		startX = -halfWidth
		endX = halfWidth - 1
	} else {
		// Odd number of chunks: center chunk at 0
		// For 25 chunks: -12 to 12 (25 total)
		startX = -halfWidth
		endX = halfWidth
	}

	fmt.Printf("DEBUG: Generating chunks from X=%d to X=%d (total: %d chunks)\n", startX, endX, endX-startX+1)

	for cy := 0; cy < worldHeight; cy++ {
		for cx := startX; cx <= endX; cx++ {
			coord := ChunkCoord{X: cx, Y: cy}
			w.Chunks[coord] = GenerateChunk(coord.X, coord.Y)
		}
	}

	fmt.Printf("DEBUG: Generated %d chunks total\n", len(w.Chunks))

	// Print first few and last few chunk coordinates for verification
	chunkCount := 0
	fmt.Printf("DEBUG: First 5 chunks generated: ")
	for coord := range w.Chunks {
		if chunkCount < 5 {
			fmt.Printf("(%d,%d) ", coord.X, coord.Y)
		}
		chunkCount++
	}
	fmt.Printf("\n")

	// Add player entity at center of the world (chunk 0)
	playerChunkX := 0 // Always spawn at the center chunk

	fmt.Printf("DEBUG: Player spawning in chunk X=%d\n", playerChunkX)

	// Place player at the center of their spawn chunk
	centerBlockX := playerChunkX*settings.ChunkWidth + settings.ChunkWidth/2
	px := float64(centerBlockX * settings.TileSize)

	fmt.Printf("DEBUG: Initial spawn block X=%d, pixel X=%f\n", centerBlockX, px)

	// Find the surface height at the center position - try multiple points for best spawn
	bestSpawnX := centerBlockX
	bestSpawnY := 1000

	// Sample multiple positions around center to find the best spawn point
	for testX := centerBlockX - 4; testX <= centerBlockX+4; testX++ {
		testSurfaceY := FindSurfaceHeight(testX, w)
		fmt.Printf("DEBUG: Surface height at X=%d is Y=%d\n", testX, testSurfaceY)
		if testSurfaceY < bestSpawnY && testSurfaceY > 10 { // Avoid spawning too high or underground
			bestSpawnY = testSurfaceY
			bestSpawnX = testX
		}
	}

	// Use the best spawn position found
	px = float64(bestSpawnX * settings.TileSize)
	spawnY := bestSpawnY - 3 // Spawn 3 blocks above surface for safety

	fmt.Printf("DEBUG: Best spawn found at block X=%d, Y=%d\n", bestSpawnX, spawnY)

	// Ensure spawn position is reasonable
	if spawnY < 5 {
		spawnY = 5
	}
	if spawnY > 200 {
		spawnY = 200
	}

	py := float64(spawnY * settings.TileSize)
	fmt.Printf("DEBUG: Final player spawn at pixel position (%f, %f)\n", px, py)

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
