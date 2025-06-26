package game

import (
	"fmt"

	"github.com/KdntNinja/webcraft/engine/world"
)

// loadChunksAroundPlayer loads chunks in a spiral pattern around the player position
func (g *Game) loadChunksAroundPlayer() {
	chunksGenerated := 0

	// Load chunks in a spiral pattern starting from center (closest first)
	for distance := 0; distance <= g.LoadDistance && chunksGenerated < g.ChunksPerFrame; distance++ {
		// Load center chunk first (distance 0)
		if distance == 0 {
			chunkX := g.CenterChunkX
			chunkY := g.CenterChunkY
			chunkKey := chunkKey(chunkX, chunkY)
			if !g.LoadedChunks[chunkKey] {
				g.ensureChunkExists(chunkX, chunkY)
				g.LoadedChunks[chunkKey] = true
				chunksGenerated++
			}
			continue
		}

		// Load chunks at current distance in a square pattern
		for dx := -distance; dx <= distance && chunksGenerated < g.ChunksPerFrame; dx++ {
			for dy := -distance; dy <= distance && chunksGenerated < g.ChunksPerFrame; dy++ {
				// Skip chunks that are not at the current distance (not on the perimeter)
				if abs(dx) != distance && abs(dy) != distance {
					continue
				}

				chunkX := g.CenterChunkX + dx
				chunkY := g.CenterChunkY + dy

				chunkKey := chunkKey(chunkX, chunkY)
				if !g.LoadedChunks[chunkKey] {
					g.ensureChunkExists(chunkX, chunkY)
					g.LoadedChunks[chunkKey] = true
					chunksGenerated++
				}
			}
		}
	}

	// Only unload chunks occasionally to avoid constant management overhead
	if len(g.LoadedChunks)%10 == 0 { // Every 10 loaded chunks, check for unloading
		g.unloadDistantChunks()
	}
}

// ensureChunkExists makes sure a chunk exists at the given coordinates
func (g *Game) ensureChunkExists(chunkX, chunkY int) {
	// Ensure the world blocks array can accommodate this chunk
	g.expandWorldGrid(chunkX, chunkY)

	// Convert world coordinates to array indices
	arrayX, arrayY := g.worldToArrayCoords(chunkX, chunkY)

	// Check bounds
	if arrayY >= len(g.World.Blocks) || arrayX >= len(g.World.Blocks[arrayY]) {
		return
	}

	// Generate the chunk (always generate, don't check if empty since we track loaded chunks)
	g.World.Blocks[arrayY][arrayX] = world.GenerateChunk(chunkX, chunkY)
}

// expandWorldGrid expands the world blocks array to accommodate the given chunk coordinates
func (g *Game) expandWorldGrid(chunkX, chunkY int) {
	// Calculate required array dimensions
	minChunkX := chunkX
	maxChunkX := chunkX
	minChunkY := chunkY
	maxChunkY := chunkY

	// Check existing chunks to find current bounds
	if len(g.World.Blocks) > 0 && len(g.World.Blocks[0]) > 0 {
		currentMinX, currentMinY := g.arrayToWorldCoords(0, 0)
		currentMaxX := currentMinX + len(g.World.Blocks[0]) - 1
		currentMaxY := currentMinY + len(g.World.Blocks) - 1

		minChunkX = min(minChunkX, currentMinX)
		maxChunkX = max(maxChunkX, currentMaxX)
		minChunkY = min(minChunkY, currentMinY)
		maxChunkY = max(maxChunkY, currentMaxY)
	}

	newWidth := maxChunkX - minChunkX + 1
	newHeight := maxChunkY - minChunkY + 1

	// Create new blocks array if needed
	if len(g.World.Blocks) == 0 || len(g.World.Blocks[0]) == 0 {
		g.World.Blocks = make([][]world.Chunk, newHeight)
		for y := 0; y < newHeight; y++ {
			g.World.Blocks[y] = make([]world.Chunk, newWidth)
		}
		g.World.MinChunkX = minChunkX
		g.World.MinChunkY = minChunkY
		return
	}

	// Expand existing array if necessary
	if newWidth > len(g.World.Blocks[0]) || newHeight > len(g.World.Blocks) ||
		minChunkX < g.World.MinChunkX || minChunkY < g.World.MinChunkY {

		// Create new larger array
		newBlocks := make([][]world.Chunk, newHeight)
		for y := 0; y < newHeight; y++ {
			newBlocks[y] = make([]world.Chunk, newWidth)
		}

		// Copy existing chunks to new positions
		for y := 0; y < len(g.World.Blocks); y++ {
			for x := 0; x < len(g.World.Blocks[y]); x++ {
				oldChunkX, oldChunkY := g.arrayToWorldCoords(x, y)
				newArrayX := oldChunkX - minChunkX
				newArrayY := oldChunkY - minChunkY
				newBlocks[newArrayY][newArrayX] = g.World.Blocks[y][x]
			}
		}

		g.World.Blocks = newBlocks
		g.World.MinChunkX = minChunkX
		g.World.MinChunkY = minChunkY
	}
}

// worldToArrayCoords converts world chunk coordinates to array indices
func (g *Game) worldToArrayCoords(chunkX, chunkY int) (int, int) {
	return chunkX - g.World.MinChunkX, chunkY - g.World.MinChunkY
}

// arrayToWorldCoords converts array indices to world chunk coordinates
func (g *Game) arrayToWorldCoords(arrayX, arrayY int) (int, int) {
	return arrayX + g.World.MinChunkX, arrayY + g.World.MinChunkY
}

// unloadDistantChunks removes chunks that are too far from the player
func (g *Game) unloadDistantChunks() {
	for chunkKey := range g.LoadedChunks {
		chunkX, chunkY := parseChunkKey(chunkKey)

		// Calculate distance from player
		dx := abs(chunkX - g.CenterChunkX)
		dy := abs(chunkY - g.CenterChunkY)
		distance := max(dx, dy)

		// Unload if too far
		if distance > g.UnloadDistance {
			arrayX, arrayY := g.worldToArrayCoords(chunkX, chunkY)
			if arrayY < len(g.World.Blocks) && arrayX < len(g.World.Blocks[arrayY]) {
				// Clear the chunk (set to all air)
				g.World.Blocks[arrayY][arrayX] = world.Chunk{}
			}
			delete(g.LoadedChunks, chunkKey)
		}
	}
}

// chunkKey creates a unique key for a chunk coordinate
func chunkKey(chunkX, chunkY int) string {
	return fmt.Sprintf("%d,%d", chunkX, chunkY)
}

// parseChunkKey parses a chunk key back to coordinates
func parseChunkKey(key string) (int, int) {
	var x, y int
	fmt.Sscanf(key, "%d,%d", &x, &y)
	return x, y
}
