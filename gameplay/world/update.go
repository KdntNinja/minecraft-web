package world

import (
	"sync"

	"github.com/KdntNinja/webcraft/coretypes"
)

// Update method to handle dynamic chunk loading
func (w *World) Update() {
	// Update chunk loading based on player position
	if len(w.Entities) > 0 {
		posX, posY := w.Entities[0].GetPosition()
		w.ChunkManager.UpdatePlayerPosition(posX, posY)
	}

	// Update entities directly in parallel (more efficient than task queue for this)
	var wg sync.WaitGroup
	for _, e := range w.Entities {
		wg.Add(1)
		go func(ent coretypes.Entity) {
			defer wg.Done()
			ent.Update()
		}(e)
	}
	wg.Wait()
}
