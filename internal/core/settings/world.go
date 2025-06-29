package settings

// --- Noise/World Generation ---
const (
	PerlinAlpha       = 2.0 // Perlin noise smoothness (higher = smoother terrain)
	PerlinBeta        = 2.0 // Perlin noise detail (higher = more detail)
	PerlinOctaves     = 2   // Number of Perlin noise octaves (layers)
	PerlinPersistence = 0.5 // Perlin octave contribution (blending factor)

	ChunkWidth        = 16  // Chunk width in blocks
	ChunkHeight       = 256 // Chunk height in blocks
	WorldChunksY      = 20  // Number of chunks vertically in the world
	WorldChunksX      = 25  // Number of chunks horizontally in the world
	TileSize          = 42  // Tile size in pixels
	DefaultSeed       = 42  // Default world generation seed
	ChunkViewDistance = 3   // Chunks to keep loaded around the player
)
