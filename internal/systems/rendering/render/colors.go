package render

import (
	"image/color"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

var (
	// Pre-calculated colors for faster access
	blockColors      [33]color.RGBA // Pre-calculated array for all block types
	colorInitialized bool
)

func initBlockColors() {
	if colorInitialized {
		return
	}

	// Pre-calculate all block colors for faster access
	blockColors[block.Grass] = color.RGBA{106, 190, 48, 255}
	blockColors[block.Dirt] = color.RGBA{151, 105, 79, 255}
	blockColors[block.Sand] = color.RGBA{238, 203, 173, 255}
	blockColors[block.Clay] = color.RGBA{168, 85, 65, 255}
	blockColors[block.Snow] = color.RGBA{255, 255, 255, 255}
	blockColors[block.Ice] = color.RGBA{173, 216, 230, 255}
	blockColors[block.Stone] = color.RGBA{100, 100, 100, 255}
	blockColors[block.Granite] = color.RGBA{120, 120, 120, 255}
	blockColors[block.Marble] = color.RGBA{245, 245, 245, 255}
	blockColors[block.Obsidian] = color.RGBA{50, 50, 50, 255}
	blockColors[block.CopperOre] = color.RGBA{184, 115, 51, 255}
	blockColors[block.IronOre] = color.RGBA{192, 192, 192, 255}
	blockColors[block.SilverOre] = color.RGBA{211, 211, 211, 255}
	blockColors[block.GoldOre] = color.RGBA{255, 215, 0, 255}
	blockColors[block.PlatinumOre] = color.RGBA{229, 228, 226, 255}
	blockColors[block.Mud] = color.RGBA{101, 67, 33, 255}
	blockColors[block.Ash] = color.RGBA{128, 128, 128, 255}
	blockColors[block.Silt] = color.RGBA{139, 119, 101, 255}
	blockColors[block.Cobweb] = color.RGBA{220, 220, 220, 128}
	blockColors[block.Hellstone] = color.RGBA{139, 0, 0, 255}
	blockColors[block.HellstoneOre] = color.RGBA{255, 69, 0, 255}
	blockColors[block.Wood] = color.RGBA{139, 69, 19, 255}
	blockColors[block.Leaves] = color.RGBA{34, 139, 34, 255}
	blockColors[block.Water] = color.RGBA{0, 191, 255, 180}
	blockColors[block.Lava] = color.RGBA{255, 69, 0, 255}
	blockColors[block.Air] = color.RGBA{135, 206, 235, 255}

	colorInitialized = true
}

func getBlockColorFast(blockType block.BlockType) color.RGBA {
	if int(blockType) < len(blockColors) {
		return blockColors[blockType]
	}
	return color.RGBA{0, 0, 0, 255} // Black for unknown blocks
}

// BlockColor returns the color for a given block type (legacy function for compatibility)
func BlockColor(b block.BlockType) color.Color {
	switch b {
	// Surface blocks
	case block.Grass:
		return color.RGBA{106, 190, 48, 255} // Green
	case block.Dirt:
		return color.RGBA{151, 105, 79, 255} // Brown
	case block.Sand:
		return color.RGBA{238, 203, 173, 255} // Sandy yellow
	case block.Clay:
		return color.RGBA{168, 85, 65, 255} // Reddish brown
	case block.Snow:
		return color.RGBA{255, 255, 255, 255} // White
	case block.Ice:
		return color.RGBA{173, 216, 230, 255} // Light blue

	// Stone variants
	case block.Stone:
		return color.RGBA{100, 100, 100, 255} // Gray
	case block.Granite:
		return color.RGBA{120, 120, 120, 255} // Light gray
	case block.Marble:
		return color.RGBA{245, 245, 245, 255} // Off-white
	case block.Obsidian:
		return color.RGBA{50, 50, 50, 255} // Dark gray/black

	// Ore blocks
	case block.CopperOre:
		return color.RGBA{184, 115, 51, 255} // Orange-brown
	case block.IronOre:
		return color.RGBA{192, 192, 192, 255} // Silver
	case block.SilverOre:
		return color.RGBA{211, 211, 211, 255} // Light silver
	case block.GoldOre:
		return color.RGBA{255, 215, 0, 255} // Gold
	case block.PlatinumOre:
		return color.RGBA{229, 228, 226, 255} // Platinum white

	// Underground blocks
	case block.Mud:
		return color.RGBA{101, 67, 33, 255} // Dark brown
	case block.Ash:
		return color.RGBA{128, 128, 128, 255} // Gray
	case block.Silt:
		return color.RGBA{139, 119, 101, 255} // Grayish brown

	// Cave blocks
	case block.Cobweb:
		return color.RGBA{220, 220, 220, 128} // Semi-transparent gray

	// Hell/Underworld blocks
	case block.Hellstone:
		return color.RGBA{139, 0, 0, 255} // Dark red
	case block.HellstoneOre:
		return color.RGBA{255, 69, 0, 255} // Orange-red

	// Tree blocks
	case block.Wood:
		return color.RGBA{139, 69, 19, 255} // Saddle brown
	case block.Leaves:
		return color.RGBA{34, 139, 34, 255} // Forest green

	// Liquids
	case block.Water:
		return color.RGBA{0, 191, 255, 180} // Semi-transparent blue
	case block.Lava:
		return color.RGBA{255, 69, 0, 255} // Orange-red

	case block.Air:
		return color.RGBA{135, 206, 235, 255} // Sky blue
	}
	return color.Black
}
