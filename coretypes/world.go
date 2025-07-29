package coretypes

// World is the interface for the game world, used by engine and rendering to avoid import cycles.
type World interface {
	GetEntities() []Entity
	Update()
	IsGridDirty() bool
	ToIntGrid() ([][]int, int, int)
	GetChunksForRendering() interface{}
	BreakBlock(x, y int) bool
	PlaceBlock(x, y int, blockType BlockType) bool
	GetBlockAt(x, y int) BlockType
	Stop()
}
