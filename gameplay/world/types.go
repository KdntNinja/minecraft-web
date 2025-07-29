package world

import (
	"context"
	"sync"

	"github.com/KdntNinja/webcraft/coretypes"
)

// AsyncUpdateTask represents a task for async world updates
type AsyncUpdateTask struct {
	taskType string
	data     interface{}
	callback func(interface{})
}

type World struct {
	ChunkManager coretypes.ChunkManager // Dynamic chunk loading system
	Entities     coretypes.Entities     // All entities in the world

	// Performance optimization caches
	cachedGrid        [][]int // Cached collision grid
	cachedGridOffsetX int     // Cached grid offset X
	cachedGridOffsetY int     // Cached grid offset Y
	gridDirty         bool    // Flag to indicate grid needs regeneration

	// Async update system
	updateTasks      chan AsyncUpdateTask
	updateWorkers    sync.WaitGroup
	updateCtx        context.Context
	updateCancel     context.CancelFunc
	numUpdateWorkers int

	// Grid generation pool
	gridGenerationMutex sync.RWMutex
	gridGenerationPool  sync.Pool
}

// GetEntities returns all entities in the world (implements coretypes.World)
func (w *World) GetEntities() []coretypes.Entity {
	return w.Entities
}
