package block

import "github.com/KdntNinja/webcraft/internal/core/settings"

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

type Chunk [settings.ChunkHeight][settings.ChunkWidth]BlockType
