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

	// Extreme terrain generation with multiple octaves for dramatic landscapes
	baseHeight := settings.SurfaceBaseHeight
	fx := float64(x)

	// Large scale mountains and valleys - very dramatic
	mountains := terrainNoise.Noise1D(fx*0.005) * float64(settings.SurfaceHeightVar) * 2.0

	// Medium scale hills - adds complexity
	hills := terrainNoise.Noise1D(fx*0.015) * float64(settings.SurfaceHeightVar) * 1.2

	// Rolling terrain for natural variation
	rolling := terrainNoise.Noise1D(fx*0.04) * float64(settings.SurfaceHeightVar) * 0.8

	// Fine detail for realistic texture
	detail := terrainNoise.Noise1D(fx*0.12) * float64(settings.SurfaceHeightVar) * 0.3

	// Sharp ridges for extreme terrain features
	ridges := terrainNoise.RidgedNoise1D(fx*0.02, 3, 0.02, 1.0) * float64(settings.SurfaceHeightVar) * 1.5

	// Combine all layers for extreme terrain
	height := baseHeight + int(mountains+hills+rolling+detail+ridges)

	// Allow more extreme heights
	if height < 5 {
		height = 5
	}
	if height > (settings.WorldChunksY*settings.ChunkHeight)-5 {
		height = (settings.WorldChunksY * settings.ChunkHeight) - 5
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
