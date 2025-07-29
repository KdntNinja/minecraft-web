package rendering

import (
	"image/color"

	"github.com/KdntNinja/webcraft/coretypes"
)

// getBlockColorFast returns fallback colors for blocks when textures aren't available
func getBlockColorFast(blockType coretypes.BlockType) color.RGBA {
	switch blockType {
	case coretypes.Grass:
		return color.RGBA{106, 190, 48, 255} // Green
	case coretypes.Dirt:
		return color.RGBA{151, 105, 79, 255} // Brown
	case coretypes.Clay:
		return color.RGBA{168, 85, 65, 255} // Reddish brown
	case coretypes.Stone:
		return color.RGBA{100, 100, 100, 255} // Gray
	case coretypes.Granite:
		return color.RGBA{145, 110, 95, 255} // Gray-pink
	case coretypes.Andesite:
		return color.RGBA{120, 130, 140, 255} // Bluish-gray
	case coretypes.Diorite:
		return color.RGBA{200, 200, 200, 255} // White-gray
	case coretypes.Slate:
		return color.RGBA{60, 60, 70, 255} // Dark gray
	case coretypes.CopperOre:
		return color.RGBA{184, 115, 51, 255} // Orange-brown
	case coretypes.IronOre:
		return color.RGBA{192, 192, 192, 255} // Silver
	case coretypes.GoldOre:
		return color.RGBA{255, 215, 0, 255} // Gold
	case coretypes.Ash:
		return color.RGBA{128, 128, 128, 255} // Gray
	case coretypes.Hellstone:
		return color.RGBA{139, 0, 0, 255} // Dark red
	case coretypes.Wood:
		return color.RGBA{139, 69, 19, 255} // Saddle brown
	case coretypes.Leaves:
		return color.RGBA{34, 139, 34, 255} // Forest green
	case coretypes.Water:
		return color.RGBA{0, 191, 255, 180} // Semi-transparent blue
	case coretypes.Air:
		return color.RGBA{135, 206, 235, 255} // Sky blue
	default:
		return color.RGBA{255, 0, 255, 255} // Magenta for unknown blocks
	}
}
