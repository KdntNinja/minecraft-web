package generation

import (
	"fmt"
	"math/rand"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// GenerateTreeAtPosition generates a tree at the specified position in a chunk
func GenerateTreeAtPosition(chunk *block.Chunk, x, surfaceChunkY int, rng *rand.Rand) {
	// Determine tree height with weighted random selection
	// Heights: 2 (most common), 3, 1, 6 (rarest)
	var treeHeight int
	heightRoll := rng.Float64()
	if heightRoll < 0.5 {
		treeHeight = 2 // 50% chance - most common
	} else if heightRoll < 0.8 {
		treeHeight = 3 // 30% chance
	} else if heightRoll < 0.95 {
		treeHeight = 1 // 15% chance
	} else {
		treeHeight = 6 // 5% chance - rarest
	}

	fmt.Printf("TREE_DEBUG: Placing tree (height %d) at x=%d, surfaceChunkY=%d\n",
		treeHeight, x, surfaceChunkY)

	// Place wood trunk blocks
	for trunkLevel := 0; trunkLevel < treeHeight; trunkLevel++ {
		trunkChunkY := surfaceChunkY - trunkLevel
		if trunkChunkY >= 0 && trunkChunkY < settings.ChunkHeight {
			chunk[trunkChunkY][x] = block.Wood
		}
	}

	// Place leaves above the trunk (going up means decreasing world Y, so decreasing chunk Y)
	// Leaves start from the top of the trunk and go up a few more levels
	leafStartLevel := treeHeight
	leafEndLevel := treeHeight + 2 // 2-3 levels of leaves above trunk

	for leafLevel := leafStartLevel; leafLevel <= leafEndLevel; leafLevel++ {
		leafChunkY := surfaceChunkY - leafLevel
		if leafChunkY >= 0 && leafChunkY < settings.ChunkHeight {
			// Place center leaves
			chunk[leafChunkY][x] = block.Leaves

			// Place side leaves for most leaf levels (except the very top)
			if leafLevel < leafEndLevel {
				if x > 0 {
					chunk[leafChunkY][x-1] = block.Leaves
				}
				if x < settings.ChunkWidth-1 {
					chunk[leafChunkY][x+1] = block.Leaves
				}
			}
		}
	}
}
