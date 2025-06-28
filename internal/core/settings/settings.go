package settings

// --- Noise/World Generation ---
const (
	PerlinAlpha       = 2.0 // Perlin noise smoothness
	PerlinBeta        = 2.0 // Perlin noise detail
	PerlinOctaves     = 3   // Perlin noise layers
	PerlinPersistence = 0.5 // Perlin octave contribution

	ChunkWidth  = 8  // Chunk width (blocks)
	ChunkHeight = 24 // Chunk height (blocks)
	WorldHeight = 32 // World vertical limit
	TileSize    = 42 // Block size (pixels)
)

// --- Biome/Surface Generation ---
const (
	BiomeCount         = 5    // Number of biome types
	BiomeBlendDistance = 16   // Biome blend width (blocks)
	SurfaceBaseHeight  = 64   // Average surface Y
	SurfaceHeightVar   = 16   // Surface height variation
	TreeChance         = 0.15 // Tree spawn chance
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
	OreVeinChance = 0.02 // Ore vein chance
	CaveFrequency = 0.08 // Cave frequency
	CaveThreshold = 0.5  // Cave noise threshold
)

// --- Chunk Generation ---
const (
	ChunkGenRadiusLeft  = 7                             // Chunks to the left of the player
	ChunkGenRadiusRight = 7                             // Chunks to the right of the player
	ChunkGenBuffer      = (2 * ChunkWidth) / ChunkWidth // 2 chunks (in blocks) => 2 chunk coords
)
