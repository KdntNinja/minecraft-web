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
	// Stone variants
	Stone
	// Ore blocks
	CopperOre
	IronOre
	GoldOre
	// Underground blocks
	Mud
	Ash
	// Tree blocks
	Wood
	Leaves
	// Liquids (for future use)
	Water
	// Hell/Underworld blocks (simplified)
	Hellstone
)

type Chunk [settings.ChunkHeight][settings.ChunkWidth]BlockType
