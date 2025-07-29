package coretypes

type ChunkCoord struct {
	X int
	Y int
}

type ChunkGenerationJob struct {
	Coord    ChunkCoord
	Priority int // Higher values = higher priority
}

type Chunk struct {
	// Define the chunk data structure here, e.g. 2D slice of BlockType
	Blocks [][]BlockType
}

type ChunkManager interface {
	GetChunk(chunkX, chunkY int) *Chunk
	UpdatePlayerPosition(playerX, playerY float64)
	GetAllChunks() map[ChunkCoord]*Chunk
	SetBlock(x, y int, blockType BlockType) bool
	GetBlock(x, y int) BlockType
	InitialLoadWithProgress(playerX, playerY float64)
	GetLoadedChunkCount() int
	Shutdown()
}
