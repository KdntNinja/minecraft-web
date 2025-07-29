package world

import (
	"github.com/KdntNinja/webcraft/coretypes"
)

// GetBlockAt returns the block type at the given world coordinates
func (w *World) GetBlockAt(blockX, blockY int) coretypes.BlockType {
	return w.ChunkManager.GetBlock(blockX, blockY)
}

// SetBlockAt sets the block type at the given world coordinates
func (w *World) SetBlockAt(blockX, blockY int, blockType coretypes.BlockType) bool {
	success := w.ChunkManager.SetBlock(blockX, blockY, blockType)

	if success {
		w.updateCachedGridBlock(blockX, blockY, blockType)
		// Do NOT always mark gridDirty here; updateCachedGridBlock will do so only if needed
	}

	return success
}

// BreakBlock removes a block at the given coordinates
func (w *World) BreakBlock(blockX, blockY int) bool {
	currentBlock := w.GetBlockAt(blockX, blockY)
	if currentBlock == coretypes.Air {
		return false // Cannot break air
	}

	return w.SetBlockAt(blockX, blockY, coretypes.Air)
}

// PlaceBlock places a block at the given coordinates
func (w *World) PlaceBlock(blockX, blockY int, blockType coretypes.BlockType) bool {
	currentBlock := w.GetBlockAt(blockX, blockY)
	if currentBlock != coretypes.Air {
		return false // Cannot place block where one already exists
	}

	// Don't allow placing air blocks
	if blockType == coretypes.Air {
		return false
	}

	// Prevent placing a block inside any entity (including player)
	if w.wouldBlockCollideWithEntity(blockX, blockY) {
		return false
	}

	return w.SetBlockAt(blockX, blockY, blockType)
}
