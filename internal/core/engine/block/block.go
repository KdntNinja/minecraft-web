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

// blockNames provides a string representation for each block type.
var blockNames = [...]string{
	"Air",
	"Grass",
	"Dirt",
	"Clay",
	"Stone",
	"Granite",
	"Andesite",
	"Diorite",
	"Slate",
	"Copper Ore",
	"Iron Ore",
	"Gold Ore",
	"Ash",
	"Wood",
	"Leaves",
	"Water",
	"Hellstone",
}

// String returns the human-readable name of the block type.
func (b BlockType) String() string {
	if b >= 0 && int(b) < len(blockNames) {
		return blockNames[b]
	}
	return "Unknown"
}

type Chunk [settings.ChunkHeight][settings.ChunkWidth]BlockType
