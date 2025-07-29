package world

// updateCachedGridBlock efficiently updates a single block in the cached collision grid
import "github.com/KdntNinja/webcraft/coretypes"

func (w *World) updateCachedGridBlock(blockX, blockY int, blockType coretypes.BlockType) {
	// Only update if we have a cached grid
	if w.cachedGrid == nil {
		return
	}

	// Convert world coordinates to grid coordinates
	gridX := blockX - w.cachedGridOffsetX
	gridY := blockY - w.cachedGridOffsetY

	// Bounds check for the cached grid
	if gridY < 0 || gridY >= len(w.cachedGrid) || gridX < 0 || gridX >= len(w.cachedGrid[0]) {
		// Block is outside cached grid bounds, mark as dirty for next full regeneration
		w.gridDirty = true
		return
	}

	// Update just this one block in the cached grid
	w.cachedGrid[gridY][gridX] = int(blockType)

	// Grid is still valid, no need to mark as dirty
}
