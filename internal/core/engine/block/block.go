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
	// Additional blocks for biome variety
	Sandstone
	Limestone
	// Liquids (for future use)
	Water
	Lava
)

const (
	TileSize    = 42  // Block size in pixels
	TilesX      = 32  // Horizontal tiles visible on screen
	ChunkWidth  = 16  // Blocks per chunk horizontally
	ChunkHeight = 256 // Blocks per chunk vertically
)

type Chunk [ChunkHeight][ChunkWidth]BlockType
