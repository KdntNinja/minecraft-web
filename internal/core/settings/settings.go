package settings

// --- Noise/World Generation ---
const (
	PerlinAlpha       = 2.0 // Perlin noise smoothness
	PerlinBeta        = 2.0 // Perlin noise detail
	PerlinOctaves     = 2   // Reduced from 3 for better performance
	PerlinPersistence = 0.5 // Perlin octave contribution

	ChunkWidth   = 16  // Chunk width (blocks)
	ChunkHeight  = 256 // Chunk height (blocks)
	WorldChunksY = 20  // Number of chunks vertically in the world
	WorldChunksX = 25  // Number of chunks horizontally in the world
	TileSize     = 42  // Increased from 30 for more zoomed-in view
	DefaultSeed  = 42  // Default world generation seed
)

// --- Biome/Surface Generation ---
const (
	BiomeCount         = 5    // Number of biome types
	BiomeBlendDistance = 16   // Biome blend width (blocks)
	SurfaceBaseHeight  = 64   // Average surface Y
	SurfaceHeightVar   = 60   // Surface height variation - increased for more dramatic terrain
	TreeChance         = 0.08 // Reduced tree spawn chance for more realistic forests
)

// --- Player/Entity Physics ---
const (
	PlayerWidth        = TileSize     // Player width (px)
	PlayerHeight       = TileSize * 2 // Player height (px)
	PlayerMoveSpeed    = 4.3          // Player move speed
	PlayerJumpSpeed    = -12.0        // Player jump velocity
	PlayerGravity      = 0.7          // Player gravity
	PlayerMaxFallSpeed = 15.0         // Player terminal velocity

	PlayerGroundFriction  = 0.6  // Ground friction
	PlayerAirResistance   = 0.98 // Air resistance
	PlayerGroundThreshold = 0.1  // Jump ground threshold
)

// --- Ore/Cave Generation ---
const (
	OreVeinChance = 0.03 // Increased ore vein chance for more resources
	CaveFrequency = 0.08 // Cave frequency
	CaveThreshold = 0.6  // Cave noise threshold - higher for fewer but larger caves
)

// --- Rendering ---
const (
	AtlasTileSize = 8 // Atlas tile size (px)
)
