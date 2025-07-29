package world

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/KdntNinja/webcraft/worldgen"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/gameplay"
	"github.com/KdntNinja/webcraft/progress"
)

// NewWorld constructs a new World instance with dynamic chunk loading and a given spawn point
func NewWorld(seed int64, chunkManager coretypes.ChunkManager, spawnPoint worldgen.SpawnPoint) *World {
	// Step: World Setup
	progress.UpdateCurrentStepProgress(1, "Setting up world structure...")
	progress.UpdateCurrentStepProgress(2, "Reset world generation")

	w := &World{
		Entities:         coretypes.Entities{},
		gridDirty:        true,
		numUpdateWorkers: runtime.NumCPU(),                // Use all available CPUs for updates
		updateTasks:      make(chan AsyncUpdateTask, 100), // Buffered channel for update tasks
		ChunkManager:     chunkManager,
	}

	// Initialize async update system
	w.updateCtx, w.updateCancel = context.WithCancel(context.Background())
	w.startUpdateWorkers()

	// Initialize grid generation pool
	w.gridGenerationPool = sync.Pool{
		New: func() interface{} {
			return make([][]int, 0, 256) // Pre-allocate slice capacity
		},
	}

	progress.UpdateCurrentStepProgress(3, "Created world structure")
	progress.CompleteCurrentStep()

	// Step: Finding Player Spawn
	progress.UpdateCurrentStepProgress(1, "Finding spawn location...")
	// Use provided spawnPoint
	progress.UpdateCurrentStepProgress(2, fmt.Sprintf("Found spawn at (%.1f, %.1f)", spawnPoint.X, spawnPoint.Y))
	progress.CompleteCurrentStep()

	// Step: Spawning Player (before chunk loading)
	progress.UpdateCurrentStepProgress(1, "Creating player entity...")
	playerEntity := gameplay.NewPlayer(spawnPoint.X, spawnPoint.Y, w)
	w.Entities = append(w.Entities, playerEntity)
	progress.UpdateCurrentStepProgress(2, "Created player entity")
	progress.CompleteCurrentStep()

	// Step: Generating Terrain around player
	progress.UpdateCurrentStepProgress(1, "Loading initial chunks around player...")
	w.ChunkManager.InitialLoadWithProgress(playerEntity.X, playerEntity.Y)
	progress.CompleteCurrentStep()

	// Step: Finalizing
	progress.UpdateCurrentStepProgress(1, "World generation finished!")
	progress.CompleteCurrentStep()

	// Initialize the grid generation pool
	w.gridGenerationPool = sync.Pool{
		New: func() interface{} {
			// Allocate a new grid slice
			return make([][]int, 0)
		},
	}

	// Start async update workers
	w.updateCtx, w.updateCancel = context.WithCancel(context.Background())
	for i := 0; i < w.numUpdateWorkers; i++ {
		w.updateWorkers.Add(1)
		go w.updateWorker(i)
	}

	return w
}
