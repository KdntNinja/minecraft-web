package physics

import (
	"context"
	"runtime"
	"sync"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/settings"
)

// PhysicsUpdateJob represents a physics update job
type PhysicsUpdateJob struct {
	entity       coretypes.Entity
	physicsWorld *PhysicsWorld
	callback     func(coretypes.Entity)
}

// AsyncPhysicsSystem handles multithreaded physics updates
type AsyncPhysicsSystem struct {
	jobs       chan PhysicsUpdateJob
	workers    sync.WaitGroup
	numWorkers int
	ctx        context.Context
	cancel     context.CancelFunc

	// Spatial partitioning for better performance
	spatialGrid map[int]map[int][]coretypes.Entity
	gridMutex   sync.RWMutex
	cellSize    int
}

var (
	asyncPhysics *AsyncPhysicsSystem
	physicsOnce  sync.Once
)

// GetAsyncPhysicsSystem returns the singleton async physics system
func GetAsyncPhysicsSystem() *AsyncPhysicsSystem {
	physicsOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		asyncPhysics = &AsyncPhysicsSystem{
			jobs:        make(chan PhysicsUpdateJob, 256),
			numWorkers:  runtime.NumCPU(),
			ctx:         ctx,
			cancel:      cancel,
			spatialGrid: make(map[int]map[int][]coretypes.Entity),
			cellSize:    settings.TileSize * 4, // Each cell is 4x4 tiles
		}
		asyncPhysics.startWorkers()
	})
	return asyncPhysics
}

// startWorkers starts the physics worker pool
func (aps *AsyncPhysicsSystem) startWorkers() {
	for i := 0; i < aps.numWorkers; i++ {
		aps.workers.Add(1)
		go aps.physicsWorker(i)
	}
}

// physicsWorker processes physics update jobs
func (aps *AsyncPhysicsSystem) physicsWorker(workerID int) {
	defer aps.workers.Done()

	for {
		select {
		case <-aps.ctx.Done():
			return
		case job := <-aps.jobs:
			aps.processPhysicsJob(job)
		}
	}
}

// processPhysicsJob processes a single physics job
func (aps *AsyncPhysicsSystem) processPhysicsJob(job PhysicsUpdateJob) {
	// Update the entity's physics
	job.entity.Update()

	// If there's a physics world, handle collisions
	if job.physicsWorld != nil {
		// Apply physics collision detection
		switch e := job.entity.(type) {
		case interface{ CollideBlocks(*PhysicsWorld) }:
			e.CollideBlocks(job.physicsWorld)
		}
	}

	// Call callback if provided
	if job.callback != nil {
		job.callback(job.entity)
	}
}

// UpdateSpatialGrid updates the spatial partitioning grid
func (aps *AsyncPhysicsSystem) UpdateSpatialGrid(entities []coretypes.Entity) {
	aps.gridMutex.Lock()
	defer aps.gridMutex.Unlock()

	// Clear the grid
	aps.spatialGrid = make(map[int]map[int][]coretypes.Entity)

	// Add entities to spatial grid
	for _, e := range entities {
		x, y := e.GetPosition()
		cellX := int(x) / aps.cellSize
		cellY := int(y) / aps.cellSize

		if aps.spatialGrid[cellX] == nil {
			aps.spatialGrid[cellX] = make(map[int][]coretypes.Entity)
		}
		aps.spatialGrid[cellX][cellY] = append(aps.spatialGrid[cellX][cellY], e)
	}
}

// GetEntitiesInRadius returns entities within a radius of a position
func (aps *AsyncPhysicsSystem) GetEntitiesInRadius(x, y, radius float64) []coretypes.Entity {
	aps.gridMutex.RLock()
	defer aps.gridMutex.RUnlock()

	var results []coretypes.Entity
	cellRadius := int(radius)/aps.cellSize + 1
	centerCellX := int(x) / aps.cellSize
	centerCellY := int(y) / aps.cellSize

	for dx := -cellRadius; dx <= cellRadius; dx++ {
		for dy := -cellRadius; dy <= cellRadius; dy++ {
			cellX := centerCellX + dx
			cellY := centerCellY + dy

			if cells, exists := aps.spatialGrid[cellX]; exists {
				if entities, exists := cells[cellY]; exists {
					for _, entity := range entities {
						ex, ey := entity.GetPosition()
						distance := (ex-x)*(ex-x) + (ey-y)*(ey-y)
						if distance <= radius*radius {
							results = append(results, entity)
						}
					}
				}
			}
		}
	}

	return results
}

// SubmitPhysicsJob submits a physics job to the worker pool
func (aps *AsyncPhysicsSystem) SubmitPhysicsJob(job PhysicsUpdateJob) {
	select {
	case aps.jobs <- job:
	default:
		// If queue is full, process synchronously
		aps.processPhysicsJob(job)
	}
}

// ProcessEntitiesAsync processes multiple entities in parallel
func (aps *AsyncPhysicsSystem) ProcessEntitiesAsync(entities []coretypes.Entity, physicsWorld *PhysicsWorld, callback func(coretypes.Entity)) {
	var wg sync.WaitGroup

	for _, e := range entities {
		wg.Add(1)
		go func(entity coretypes.Entity) {
			defer wg.Done()
			job := PhysicsUpdateJob{
				entity:       entity,
				physicsWorld: physicsWorld,
				callback:     callback,
			}
			aps.processPhysicsJob(job)
		}(e)
	}

	wg.Wait()
}

// Shutdown stops all physics workers
func (aps *AsyncPhysicsSystem) Shutdown() {
	aps.cancel()
	aps.workers.Wait()
	close(aps.jobs)
}
