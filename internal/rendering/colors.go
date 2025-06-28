package rendering

import (
	"image/color"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
)

// getBlockColorFast returns fallback colors for blocks when textures aren't available
func getBlockColorFast(blockType block.BlockType) color.RGBA {
	switch blockType {
	case block.Grass:
		return color.RGBA{106, 190, 48, 255} // Green
	case block.Dirt:
		return color.RGBA{151, 105, 79, 255} // Brown
	case block.Clay:
		return color.RGBA{168, 85, 65, 255} // Reddish brown
	case block.Stone:
		return color.RGBA{100, 100, 100, 255} // Gray
	case block.CopperOre:
		return color.RGBA{184, 115, 51, 255} // Orange-brown
	case block.IronOre:
		return color.RGBA{192, 192, 192, 255} // Silver
	case block.GoldOre:
		return color.RGBA{255, 215, 0, 255} // Gold
	case block.Ash:
		return color.RGBA{128, 128, 128, 255} // Gray
	case block.Hellstone:
		return color.RGBA{139, 0, 0, 255} // Dark red
	case block.Wood:
		return color.RGBA{139, 69, 19, 255} // Saddle brown
	case block.Leaves:
		return color.RGBA{34, 139, 34, 255} // Forest green
	case block.Water:
		return color.RGBA{0, 191, 255, 180} // Semi-transparent blue
	case block.Air:
		return color.RGBA{135, 206, 235, 255} // Sky blue
	default:
		return color.RGBA{255, 0, 255, 255} // Magenta for unknown blocks
	}
}
