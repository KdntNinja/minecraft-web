package world

import "fmt"

// Stop stops the world and cleans up resources
func (w *World) Stop() {
	// Cancel the update context to stop all workers
	w.updateCancel()

	// Wait for all update workers to finish
	w.updateWorkers.Wait()
	close(w.updateTasks)
	w.ChunkManager.Shutdown()
	fmt.Println("WORLD: Shutdown complete")
}
