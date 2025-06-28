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
	Chunks     map[ChunkCoord]block.Chunk // Infinite world: map of chunk coordinates to chunks
	ChunkCache map[ChunkCoord]block.Chunk // Cache for previously generated chunks
	Entities   entity.Entities
	MinChunkX  int // Minimum chunk X coordinate in the world grid (optional, can remove)
	MinChunkY  int // Minimum chunk Y coordinate in the world grid (optional, can remove)
}

// NewWorld constructs a new World instance with generated chunks, using the provided seed
func NewWorld(numChunksY int, centerChunkX int, seed int64) *World {
	terrain.ResetWorldGeneration(seed)
	w := &World{
		Chunks:     make(map[ChunkCoord]block.Chunk),
		ChunkCache: make(map[ChunkCoord]block.Chunk),
		Entities:   entity.Entities{},
	}
	// Generate initial window of chunks around center
	for cy := 0; cy < numChunksY; cy++ {
		for cx := -3; cx <= 3; cx++ {
			coord := ChunkCoord{X: centerChunkX + cx, Y: cy}
			w.Chunks[coord] = GenerateChunk(coord.X, coord.Y)
		}
	}
	// Add player entity at center
	centerChunkCol := 0 // center is always X=0
	centerBlockX := centerChunkCol*settings.ChunkWidth + settings.ChunkWidth/2
	px := float64(centerBlockX * settings.TileSize)

	// Find the surface height at the center position
	surfaceY := FindSurfaceHeight(centerBlockX, w)

	// Spawn player 2 blocks above the surface for safety
	spawnY := surfaceY - 2

	// Ensure spawn position is reasonable
	if spawnY < 0 {
		spawnY = 0
	}

	py := float64(spawnY * settings.TileSize)
	w.Entities = append(w.Entities, player.NewPlayer(px, py))
	return w
}

// ToIntGrid flattens the world's blocks into a [][]int grid for entity collision, and returns the offset (minX, minY)
func (w *World) ToIntGrid() ([][]int, int, int) {
	if len(w.Chunks) == 0 {
		return [][]int{}, 0, 0
	}
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
	grid := make([][]int, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
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
	return grid, minX * settings.ChunkWidth, minY * settings.ChunkHeight
}

// UpdateChunksWindow generates and manages chunks based only on distance to the player (no window logic)
func (w *World) UpdateChunksWindow(playerX, playerY float64) {
	playerChunkX := int(playerX) / (settings.ChunkWidth * settings.TileSize)
	playerChunkY := int(playerY) / (settings.ChunkHeight * settings.TileSize)

	radiusLeft := settings.ChunkGenRadiusLeft
	radiusRight := settings.ChunkGenRadiusRight
	buffer := settings.ChunkGenBuffer

	// Use a set to track which chunks should be present after this update
	needed := make(map[ChunkCoord]struct{})

	// 1. Generate all chunks within the asymmetric radius+buffer of the player
	for cy := playerChunkY - radiusLeft - buffer; cy <= playerChunkY+radiusLeft+buffer; cy++ {
		for cx := playerChunkX - radiusLeft - buffer; cx <= playerChunkX+radiusRight+buffer; cx++ {
			coord := ChunkCoord{X: cx, Y: cy}
			needed[coord] = struct{}{}
			if _, ok := w.Chunks[coord]; !ok || len(w.Chunks[coord]) == 0 {
				// Try to restore from cache first
				if cached, found := w.ChunkCache[coord]; found {
					w.Chunks[coord] = cached
				} else {
					w.Chunks[coord] = GenerateChunk(cx, cy)
				}
			}
		}
	}

	// 2. Remove only chunks outside the (asymmetric) area, and cache them
	// Avoid modifying the map while iterating: collect toRemove first
	toRemove := make([]ChunkCoord, 0)
	for coord := range w.Chunks {
		if _, keep := needed[coord]; !keep {
			toRemove = append(toRemove, coord)
		}
	}
	for _, coord := range toRemove {
		w.ChunkCache[coord] = w.Chunks[coord]
		delete(w.Chunks, coord)
	}
}
