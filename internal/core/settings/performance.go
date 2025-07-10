package settings

import "runtime"

// Performance and multithreading settings

const (
	// --- Multithreading Configuration ---
	// MaxWorkers sets the maximum number of worker goroutines for various systems
	// Set to 0 to use runtime.NumCPU()
	MaxChunkWorkers   = 0 // Chunk generation workers
	MaxPhysicsWorkers = 0 // Physics update workers
	MaxRenderWorkers  = 0 // Rendering workers
	MaxUpdateWorkers  = 0 // General update workers

	// --- Queue Sizes ---
	ChunkJobQueueSize   = 128  // Buffer size for chunk generation jobs
	PhysicsJobQueueSize = 256  // Buffer size for physics update jobs
	RenderJobQueueSize  = 1024 // Buffer size for render jobs
	UpdateTaskQueueSize = 100  // Buffer size for general update tasks

	// --- Performance Thresholds ---
	SlowChunkGenerationThreshold = 100 // Milliseconds - log chunks that take longer
	MaxChunksPerFrame            = 2   // Limit chunk loading per frame to prevent stutter
	PhysicsUpdateInterval        = 60  // Frames between physics grid regeneration

	// --- Spatial Optimization ---
	SpatialGridCellSize = TileSize * 4 // Size of each spatial grid cell for physics

	// --- Memory Management ---
	GridPoolInitialCapacity = 256 // Initial capacity for grid object pool
	TextureAtlasSize        = 512 // Size of texture atlas for batching
	EntityBatchSize         = 32  // Number of entities to process in each batch

	// --- Culling and LOD ---
	EntityCullingMargin = TileSize * 2 // Extra margin for entity culling
	ChunkCullingMargin  = 1            // Extra chunks to render beyond visible area

	// --- Async Processing ---
	EnableAsyncChunkGeneration = true // Enable multithreaded chunk generation
	EnableAsyncPhysics         = true // Enable multithreaded physics updates
	EnableAsyncRendering       = true // Enable multithreaded rendering preparation
	EnableSpatialOptimization  = true // Enable spatial partitioning for entities

	// --- Debug and Profiling ---
	EnablePerformanceMetrics = true // Track and display performance metrics
	EnableGCOptimization     = true // Force garbage collection at optimal times
	MetricsHistorySize       = 120  // Number of frames to keep for performance history
)

// GetOptimalWorkerCount returns the optimal number of workers for a given task type
func GetOptimalWorkerCount(taskType string) int {
	switch taskType {
	case "chunk":
		if MaxChunkWorkers > 0 {
			return MaxChunkWorkers
		}
	case "physics":
		if MaxPhysicsWorkers > 0 {
			return MaxPhysicsWorkers
		}
	case "render":
		if MaxRenderWorkers > 0 {
			return MaxRenderWorkers
		}
	case "update":
		if MaxUpdateWorkers > 0 {
			return MaxUpdateWorkers
		}
	}

	// Default to number of CPU cores, but cap at 8 to avoid excessive goroutines
	cores := runtime.NumCPU()
	if cores > 8 {
		return 8
	}
	return cores
}

// Performance optimization flags
var (
	// These can be modified at runtime to tune performance
	DynamicChunkLoading   = true // Enable/disable dynamic chunk loading
	ParallelEntityUpdates = true // Enable/disable parallel entity updates
	AsyncGridGeneration   = true // Enable/disable async collision grid generation
	FramerateLimiting     = true // Enable/disable framerate limiting for consistent performance
)
