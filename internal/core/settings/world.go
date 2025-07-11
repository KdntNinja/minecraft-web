package settings

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
