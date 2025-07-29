package world

// GetChunksForRendering returns all currently loaded chunks for rendering
func (w *World) GetChunksForRendering() interface{} {
	return w.ChunkManager.GetAllChunks()
}
