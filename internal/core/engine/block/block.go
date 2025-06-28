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
	// Ore blocks
	CopperOre
	IronOre
	GoldOre
	// Underground blocks
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
