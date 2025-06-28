package terrain

import (
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
		// Use default seed if not set
		terrainSeed = settings.DefaultSeed
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

// ResetWorldGeneration forces regeneration with a new provided seed
func ResetWorldGeneration(seed int64) {
	terrainNoise = noise.NewPerlinNoise(seed)
	terrainSeed = seed
	terrainHeights = make(map[int]int)
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
