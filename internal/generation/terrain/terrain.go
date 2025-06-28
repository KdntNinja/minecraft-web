package terrain

import (
	"math/rand"
	"time"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/generation/noise"
)

var (
	terrainHeights = make(map[int]int)
	terrainNoise   *noise.PerlinNoise
	terrainSeed    int64
)

func init() {
	// Set the function variable in chunks.go to avoid import cycle
	GetWorldSeedFunc = GetWorldSeed
}

func initTerrainNoise() {
	if terrainNoise == nil {
		rand.Seed(time.Now().UnixNano())
		terrainSeed = int64(rand.Intn(1000000) + 1) // Random seed in [1, 1000000]
		terrainNoise = noise.NewPerlinNoise(terrainSeed)
		terrainHeights = make(map[int]int)
	}
}

func getSurfaceHeight(x int) int {
	h, ok := terrainHeights[x]
	if ok {
		return h
	}

	// Improved: Use fractal noise for more interesting terrain
	baseHeight := settings.SurfaceBaseHeight
	fx := float64(x)
	// Combine several octaves of noise for hills, valleys, and detail
	hill := terrainNoise.Noise1D(fx*0.01) * float64(settings.SurfaceHeightVar)
	valley := terrainNoise.Noise1D(fx*0.03) * (float64(settings.SurfaceHeightVar) * 0.5)
	detail := terrainNoise.Noise1D(fx*0.09) * (float64(settings.SurfaceHeightVar) * 0.25)

	height := baseHeight + int(hill+valley+detail)

	if height < 3 {
		height = 3
	}
	if height > settings.ChunkHeight-2 {
		height = settings.ChunkHeight - 2
	}

	terrainHeights[x] = height
	return height
}

// GenerateChunksInView generates all chunks in range of the player's view (centered on player)
func GenerateChunksInView(playerX, playerY float64, viewRadius int) {
	playerChunkX := int(playerX) / (settings.ChunkWidth * settings.TileSize)
	playerChunkY := int(playerY) / (settings.ChunkHeight * settings.TileSize)

	for cy := playerChunkY - viewRadius; cy <= playerChunkY+viewRadius; cy++ {
		for cx := playerChunkX - viewRadius; cx <= playerChunkX+viewRadius; cx++ {
			_ = GenerateChunk(cx, cy) // This will generate and cache the chunk if needed
		}
	}
}

// BiomeType represents different surface biomes
// Only a single biome for simplicity

type BiomeType int

const (
	DefaultBiome BiomeType = iota
)

func getBiome(x int) BiomeType {
	return DefaultBiome
}

// getOreType determines what ore should be placed using simple noise
func getOreType(x, y, depthFromSurface int) block.BlockType {
	oreNoise := terrainNoise.Noise2D(float64(x), float64(y))

	// Map ore type numbers to block types (simplified)
	switch {
	case oreNoise < -0.4:
		return block.CopperOre
	case oreNoise < 0:
		return block.IronOre
	case oreNoise < 0.6:
		return block.GoldOre
	default:
		return block.Stone // Default to stone
	}
}

// isInCave determines if a position should be a cave using simple noise
func isInCave(x, y int) bool {
	// Don't generate caves too close to surface
	if y < 8 {
		return false
	}

	// Use simple noise for cave generation
	caveNoise := terrainNoise.Noise2D(float64(x), float64(y))

	// Simple cave threshold
	return caveNoise < -0.2
}

// isUnderworld checks if we're in the underworld layer
func isUnderworld(y, worldHeight int) bool {
	return y > worldHeight-6 // Bottom 6 layers are underworld
}

// ResetWorldGeneration forces regeneration with a new provided seed
func ResetWorldGeneration(seed int64) {
	terrainNoise = noise.NewPerlinNoise(seed)
	terrainSeed = seed
	terrainHeights = make(map[int]int)
	chunkCache = make(map[string]block.Chunk)
}

// GetWorldSeed returns the current world seed for display or saving
func GetWorldSeed() int64 {
	initTerrainNoise()
	return terrainSeed
}

// getBlockType returns the block type for a given world position, generating Minecraft-like layers
func getBlockType(x, y int) block.BlockType {
	surfaceY := getSurfaceHeight(x)
	depth := y - surfaceY

	if y < 0 || y >= settings.ChunkHeight {
		return block.Air
	}

	if y < surfaceY {
		return block.Air
	}

	// Surface block
	if y == surfaceY {
		return block.Grass
	}

	// Dirt layer (3 blocks below surface)
	if y > surfaceY && y <= surfaceY+3 {
		return block.Dirt
	}

	// Stone and ores below dirt
	if y > surfaceY+3 && y < settings.ChunkHeight-6 {
		// Occasionally generate ores
		ore := getOreType(x, y, depth)
		if ore != block.Stone {
			return ore
		}
		return block.Stone
	}

	// Underworld (bottom 6 layers)
	if y >= settings.ChunkHeight-6 {
		return block.Hellstone
	}

	return block.Stone
}
