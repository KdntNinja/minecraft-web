package generation

import (
	"fmt"
	"math/rand"

	"github.com/KdntNinja/webcraft/coretypes"
)

// GenerateTreeAtPosition generates a tree at the specified position in a chunk
func GenerateTreeAtPosition(chunk *coretypes.Chunk, x, surfaceChunkY int, rng *rand.Rand) {
	treeType, shape := GetTreeTypeAndShape(rng)
	shape.ValidateShape()

	fmt.Printf("TREE_DEBUG: Placing %v tree (height %d) at x=%d, surfaceChunkY=%d\n",
		treeType, shape.TrunkHeight, x, surfaceChunkY)

	// Use custom pattern if available
	if shape.CustomPattern != nil {
		shape.CustomPattern(chunk, x, surfaceChunkY, rng, shape)
		return
	}

	// Generate standard tree pattern
	generateStandardTree(chunk, x, surfaceChunkY, rng, treeType, shape)
}
