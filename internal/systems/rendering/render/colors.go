package render

import (
	"image/color"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

var (
	// Pre-calculated colors for faster access
	blockColors      [20]color.RGBA // Reduced array size for simplified block types
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
	blockColors[block.Stone] = color.RGBA{100, 100, 100, 255}
	blockColors[block.CopperOre] = color.RGBA{184, 115, 51, 255}
	blockColors[block.IronOre] = color.RGBA{192, 192, 192, 255}
	blockColors[block.GoldOre] = color.RGBA{255, 215, 0, 255}
	blockColors[block.Mud] = color.RGBA{101, 67, 33, 255}
	blockColors[block.Ash] = color.RGBA{128, 128, 128, 255}
	blockColors[block.Hellstone] = color.RGBA{139, 0, 0, 255}
	blockColors[block.Wood] = color.RGBA{139, 69, 19, 255}
	blockColors[block.Leaves] = color.RGBA{34, 139, 34, 255}
	blockColors[block.Water] = color.RGBA{0, 191, 255, 180}
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

	// Stone variants
	case block.Stone:
		return color.RGBA{100, 100, 100, 255} // Gray

	// Ore blocks
	case block.CopperOre:
		return color.RGBA{184, 115, 51, 255} // Orange-brown
	case block.IronOre:
		return color.RGBA{192, 192, 192, 255} // Silver
	case block.GoldOre:
		return color.RGBA{255, 215, 0, 255} // Gold

	// Underground blocks
	case block.Mud:
		return color.RGBA{101, 67, 33, 255} // Dark brown
	case block.Ash:
		return color.RGBA{128, 128, 128, 255} // Gray

	// Hell/Underworld blocks
	case block.Hellstone:
		return color.RGBA{139, 0, 0, 255} // Dark red

	// Tree blocks
	case block.Wood:
		return color.RGBA{139, 69, 19, 255} // Saddle brown
	case block.Leaves:
		return color.RGBA{34, 139, 34, 255} // Forest green

	// Liquids
	case block.Water:
		return color.RGBA{0, 191, 255, 180} // Semi-transparent blue

	case block.Air:
		return color.RGBA{135, 206, 235, 255} // Sky blue
	}
	return color.Black
}
