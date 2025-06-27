package block

type BlockType int

const (
	Air BlockType = iota
	// Surface blocks
	Grass
	Dirt
	Sand
	Clay
	Snow
	Ice
	// Stone variants
	Stone
	Granite
	Marble
	Obsidian
	// Ore blocks
	CopperOre
	IronOre
	SilverOre
	GoldOre
	PlatinumOre
	// Underground blocks
	Mud
	Ash
	Silt
	// Cave blocks
	Cobweb
	// Hell/Underworld blocks
	Hellstone
	HellstoneOre
	// Tree blocks
	Wood
	Leaves
	// Liquids (for future use)
	Water
	Lava
)

const (
	TileSize    = 42 // Larger blocks for better visibility
	TilesX      = 32 // More tiles horizontally
	ChunkWidth  = 6  // Each chunk is 6 blocks wide
	ChunkHeight = 24 // Each chunk is 24 blocks tall (or any value for vertical size)
)

type Chunk [ChunkHeight][ChunkWidth]BlockType
