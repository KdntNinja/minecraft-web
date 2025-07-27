package settings

import "runtime"

// --- Biome/Surface Generation ---
const (
	BiomeCount         = 5    // Number of biome types
	BiomeBlendDistance = 16   // Biome blend width in blocks
	SurfaceBaseHeight  = 64   // Average surface Y
	SurfaceHeightVar   = 60   // Surface height variation
	TreeChance         = 0.08 // Probability of tree spawn per block
)

// --- Cave Generation Parameters ---
const (
	CaveSurfaceEntranceMinDepth = -4     // Min depth (relative to surface) for cave entrances
	CaveSurfaceEntranceMaxDepth = 12     // Max depth (relative to surface) for cave entrances
	CaveSurfaceEntranceScale    = 30.0   // Surface cave entrance noise scale
	CaveSurfaceEntranceOffset   = 8000.0 // Surface cave entrance noise offset
	CaveSurfaceEntranceThresh   = 0.55   // Surface cave entrance threshold (higher = fewer)

	CaveLargeScale        = 50.0   // Large cavern noise scale (broad spaces)
	CaveHorizontalScale   = 18.0   // Horizontal tunnel noise scale
	CaveHorizontalYOffset = 1000.0 // Horizontal tunnel noise Y offset
	CaveHorizontalYScale  = 35.0   // Horizontal tunnel Y blending scale
	CaveVerticalScale     = 35.0   // Vertical tunnel noise scale
	CaveVerticalYOffset   = 2000.0 // Vertical tunnel noise Y offset
	CaveVerticalYScale    = 15.0   // Vertical tunnel Y blending scale
	CaveSmallScale        = 10.0   // Small cave noise scale (pockets)
	CaveSmallYOffset      = 3000.0 // Small cave noise Y offset
	CaveAirPocketScale    = 6.0    // Air pocket noise scale (tiny holes)
	CaveAirPocketYOffset  = 4000.0 // Air pocket noise Y offset

	CaveVeryDeepDepth   = 150 // Very deep caves start below this depth
	CaveDeepDepth       = 100 // Deep caves start below this depth
	CaveMediumDepth     = 50  // Medium caves start below this depth
	CaveShallowDepth    = 15  // Shallow caves start below this depth
	CaveMinShallowDepth = 2   // Minimum depth for any cave

	CaveVeryDeepLargeWeight  = 0.4  // Very deep: large cavern weight
	CaveVeryDeepHorizWeight  = 0.3  // Very deep: horizontal tunnel weight
	CaveVeryDeepVertWeight   = 0.2  // Very deep: vertical tunnel weight
	CaveVeryDeepSmallWeight  = 0.1  // Very deep: small cave weight
	CaveVeryDeepThresh       = 0.15 // Very deep: cave generation threshold
	CaveVeryDeepTunnelWeight = 0.6  // Very deep: interconnected tunnel weight
	CaveVeryDeepPocketWeight = 0.4  // Very deep: air pocket weight
	CaveVeryDeepTunnelThresh = 0.35 // Very deep: tunnel generation threshold

	CaveDeepLargeWeight  = 0.3  // Deep: large cavern weight
	CaveDeepHorizWeight  = 0.3  // Deep: horizontal tunnel weight
	CaveDeepSmallWeight  = 0.3  // Deep: small cave weight
	CaveDeepPocketWeight = 0.1  // Deep: air pocket weight
	CaveDeepThresh       = 0.18 // Deep: cave generation threshold
	CaveDeepVertThresh   = 0.45 // Deep: vertical tunnel threshold

	CaveMediumHorizWeight  = 0.4  // Medium: horizontal tunnel weight
	CaveMediumSmallWeight  = 0.4  // Medium: small cave weight
	CaveMediumPocketWeight = 0.2  // Medium: air pocket weight
	CaveMediumThresh       = 0.22 // Medium: cave generation threshold
	CaveMediumVertThresh   = 0.5  // Medium: vertical tunnel threshold

	CaveShallowHorizWeight  = 0.3  // Shallow: horizontal tunnel weight
	CaveShallowSmallWeight  = 0.5  // Shallow: small cave weight
	CaveShallowPocketWeight = 0.2  // Shallow: air pocket weight
	CaveShallowThresh       = 0.18 // Shallow: cave generation threshold
	CaveShallowVertThresh   = 0.45 // Shallow: vertical tunnel threshold

	CaveMinShallowSmallWeight  = 0.5  // Min-depth: small cave weight
	CaveMinShallowPocketWeight = 0.3  // Min-depth: air pocket weight
	CaveMinShallowHorizWeight  = 0.2  // Min-depth: horizontal tunnel weight
	CaveMinShallowThresh       = 0.22 // Min-depth: cave generation threshold
)

// --- Debug Overlay Constants ---
const (
	DebugOverlayWidth  = 340 // Width in pixels of the debug overlay UI
	DebugOverlayHeight = 170 // Height in pixels of the debug overlay UI
	DebugGraphHeight   = 32  // Height in pixels of the debug graph area
	DebugGraphSamples  = 120 // Number of samples shown in the debug graph
)

// --- Ore Generation Constants ---
const (
	OreVeinChance = 0.08 // Base chance for an ore vein to generate at a position
)

// --- Multithreading/Performance/Rendering ---
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
	MaxChunksPerFrame            = 1   // Limit chunk loading per frame to prevent stutter (reduced from 2)
	PhysicsUpdateInterval        = 60  // Frames between physics grid regeneration

	// --- Anti-Stutter Configuration ---
	ChunkGenerationTimeSlice  = 8    // Max milliseconds per frame for chunk generation
	MaxConcurrentChunkJobs    = 4    // Max concurrent chunk generation jobs
	ChunkPriorityRadius       = 2    // Priority radius around player for chunk loading
	BackgroundChunkGeneration = true // Generate chunks in background when not needed immediately

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

// --- Player Physics ---
const (
	PlayerSpriteWidth    = TileSize             // Visual width of the player sprite.
	PlayerSpriteHeight   = TileSize * 2         // Visual height of the player sprite.
	PlayerColliderWidth  = (TileSize * 9) / 10  // Physics bounding box width for collision.
	PlayerColliderHeight = (TileSize * 18) / 10 // Physics bounding box height for collision.
	PlayerMoveSpeed      = 4.3                  // Maximum horizontal walking speed (blocks/sec).
	PlayerJumpSpeed      = -9.0                 // Initial vertical jump velocity (upwards).
	PlayerGravity        = 0.45                 // Gravity force applied each frame.
	PlayerMaxFallSpeed   = 16.0                 // Maximum downward velocity (terminal velocity).

	// --- Movement Tuning ---
	PlayerWalkAccel      = 0.25                  // Acceleration when walking on the ground.
	PlayerAirAccel       = 0.04                  // Acceleration when in the air.
	PlayerGroundFriction = 0.55                  // Friction applied when on the ground and not moving.
	PlayerAirFriction    = 0.985                 // Air resistance applied when airborne.
	PlayerSneakSpeed     = PlayerMoveSpeed * 0.3 // Movement speed when sneaking.
	PlayerSneakAccel     = PlayerWalkAccel * 0.5 // Acceleration when sneaking.
	PlayerSprintSpeed    = PlayerMoveSpeed * 1.3 // Movement speed when sprinting.
	PlayerSprintAccel    = PlayerWalkAccel * 1.2 // Acceleration when sprinting.

	// --- Jump Mechanics ---
	PlayerCoyoteFrames     = 8    // Grace period (frames) to jump after leaving a ledge.
	PlayerJumpBufferFrames = 8    // Grace period (frames) to buffer a jump before landing.
	PlayerJumpHoldMax      = 12   // Max duration (frames) to hold jump for variable height.
	PlayerJumpHoldForce    = 0.32 // Upward force applied each frame when holding jump.
)

// --- Rendering/Texture System Constants ---
const (
	AtlasTileSize      = 8    // Tile size in the texture atlas (pixels)
	TextureLogFallback = true // Log when fallback color-based textures are used
	TextureLogInit     = true // Log when texture system is initialized
)

// --- Terrain Height Generation Constants ---
const (
	TerrainBaseScale       = 120.0  // Base scale for low-frequency terrain
	TerrainHillScale       = 40.0   // Scale for hills (medium-frequency noise)
	TerrainHillOffset      = 500.0  // Offset for hill noise
	TerrainDetailScale     = 10.0   // Scale for small terrain details (high-frequency noise)
	TerrainDetailOffset    = 2000.0 // Offset for detail noise
	TerrainBaseWeight      = 0.6    // Weight for base terrain layer
	TerrainHillWeight      = 0.3    // Weight for hill layer
	TerrainDetailWeight    = 0.1    // Weight for detail layer
	TerrainMinHeight       = 8      // Minimum allowed terrain height (blocks)
	TerrainMaxHeightBuffer = 8      // Buffer from world bottom for max height
)

// --- Tree/Biome Noise Constants (used in surface.go) ---
const (
	TreeBiomeNoiseScale = 80.0 // Scale for biome noise (surface variation)
	TreeNoiseScale      = 12.0 // Scale for tree placement noise
	TreeClayBiomeThresh = 0.45 // Threshold for clay biome
	TreeClayNoiseThresh = 0.5  // Threshold for clay noise
)

// --- Terraria-like Sky Transition Constants ---
const (
	SkyTopY             = -100.0 // Y above this is always sky blue
	SkyTransitionStartY = -50.0  // Start fading here
	SkyTransitionEndY   = 150.0  // Fully dark here
)

// --- Noise/World Generation ---
const (
	PerlinAlpha       = 2.0 // Perlin noise smoothness (higher = smoother terrain)
	PerlinBeta        = 2.0 // Perlin noise detail (higher = more detail)
	PerlinOctaves     = 2   // Number of Perlin noise octaves (layers)
	PerlinPersistence = 0.5 // Perlin octave contribution (blending factor)

	ChunkWidth        = 16  // Chunk width in blocks (reduced from 32 for better performance)
	ChunkHeight       = 128 // Chunk height in blocks (reduced from 256 for faster generation)
	WorldChunksY      = 15  // Number of chunks vertically in the world (reduced)
	WorldChunksX      = 32  // Number of chunks horizontally in the world (reduced)
	TileSize          = 32  // Tile size in pixels (reduced from 30 for better rendering performance)
	DefaultSeed       = 0   // Default world generation seed
	ChunkViewDistance = 2   // Chunks to keep loaded around the player (reduced from 3)
)

// --- Performance optimization flags ---
var (
	// These can be modified at runtime to tune performance
	DynamicChunkLoading   = true // Enable/disable dynamic chunk loading
	ParallelEntityUpdates = true // Enable/disable parallel entity updates
	AsyncGridGeneration   = true // Enable/disable async collision grid generation
	FramerateLimiting     = true // Enable/disable framerate limiting for consistent performance
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
