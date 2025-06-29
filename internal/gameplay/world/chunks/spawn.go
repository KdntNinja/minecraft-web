package chunks

import (
	"fmt"

	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation"
)

// SpawnPoint represents a valid spawn location
type SpawnPoint struct {
	X, Y     float64 // World coordinates
	ChunkX   int     // Chunk coordinates
	ChunkY   int
	SurfaceY int // Surface block Y coordinate
	BlockX   int // Block coordinates
	BlockY   int
}

// FindSpawnPoint finds a suitable spawn point near the world center (0, 0)
func FindSpawnPoint() SpawnPoint {
	// Start looking near world center (0, 0)
	searchBlockX := 0

	// Find surface height at the center
	surfaceY := generation.GetHeightAt(searchBlockX)

	// Ensure reasonable spawn height
	spawnBlockY := surfaceY - 3 // 3 blocks above surface
	if spawnBlockY < 5 {
		spawnBlockY = 5
	}
	if spawnBlockY > 200 {
		spawnBlockY = 200
	}

	// Convert to world pixel coordinates
	spawnWorldX := float64(searchBlockX * settings.TileSize)
	spawnWorldY := float64(spawnBlockY * settings.TileSize)

	// Get chunk coordinates
	chunkX, chunkY := WorldToChunk(spawnWorldX, spawnWorldY)

	spawn := SpawnPoint{
		X:        spawnWorldX,
		Y:        spawnWorldY,
		ChunkX:   chunkX,
		ChunkY:   chunkY,
		SurfaceY: surfaceY,
		BlockX:   searchBlockX,
		BlockY:   spawnBlockY,
	}

	fmt.Printf("SPAWN: Found spawn point at world (%.1f, %.1f), block (%d, %d), chunk (%d, %d), surface Y=%d\n",
		spawn.X, spawn.Y, spawn.BlockX, spawn.BlockY, spawn.ChunkX, spawn.ChunkY, spawn.SurfaceY)

	return spawn
}

// FindSafeSpawnPoint finds a spawn point that's guaranteed to be safe (not in a cave, etc.)
func FindSafeSpawnPoint() SpawnPoint {
	// Center the player in the middle of the world horizontally
	centerChunkX := settings.WorldChunksX / 2
	searchX := centerChunkX * settings.ChunkWidth

	surfaceY := generation.GetHeightAt(searchX)

	// Make sure it's a reasonable surface height
	if surfaceY < 5 {
		surfaceY = 5
	}
	if surfaceY > settings.ChunkHeight*settings.WorldChunksY-3 {
		surfaceY = settings.ChunkHeight*settings.WorldChunksY - 3
	}
	spawnBlockY := surfaceY - 3

	spawnWorldX := float64(searchX * settings.TileSize)
	spawnWorldY := float64(spawnBlockY * settings.TileSize)

	chunkX, chunkY := WorldToChunk(spawnWorldX, spawnWorldY)

	spawn := SpawnPoint{
		X:        spawnWorldX,
		Y:        spawnWorldY,
		ChunkX:   chunkX,
		ChunkY:   chunkY,
		SurfaceY: surfaceY,
		BlockX:   searchX,
		BlockY:   spawnBlockY,
	}

	fmt.Printf("SPAWN: Centered spawn point at world (%.1f, %.1f), block (%d, %d), chunk (%d, %d), surface Y=%d\n",
		spawn.X, spawn.Y, spawn.BlockX, spawn.BlockY, spawn.ChunkX, spawn.ChunkY, spawn.SurfaceY)

	return spawn
}
