package coretypes

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

func (b BlockType) String() string {
	if b >= 0 && int(b) < len(blockNames) {
		return blockNames[b]
	}
	return "Unknown"
}

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

const NumBlockTypes = int(Hellstone) + 1
