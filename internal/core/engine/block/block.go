package block

import "github.com/KdntNinja/webcraft/internal/core/settings"

type BlockType int

const (
	Air BlockType = iota
	// Surface blocks
	Grass
	Dirt
	Clay
	// Stone variants
	Stone
	Granite  // Underground stone variant: gray-pink
	Andesite // Underground stone variant: bluish-gray
	Diorite  // Underground stone variant: white-gray
	Slate    // Underground stone variant: dark gray
	// Ore blocks
	CopperOre
	IronOre
	GoldOre
	// Underground blocks
	Ash
	// Tree blocks
	Wood
	Leaves
	// Liquids
	Water
	// Hell/Underworld blocks
	Hellstone
)

// NumBlockTypes is the total number of defined block types (for fast array access)
const NumBlockTypes = int(Hellstone) + 1

type Chunk [settings.ChunkHeight][settings.ChunkWidth]BlockType
