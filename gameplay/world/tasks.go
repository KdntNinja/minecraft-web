package world

import (
	"sync"

	"github.com/KdntNinja/webcraft/coretypes"
)

func (w *World) processGridGeneration(data interface{}, callback func(interface{})) {
	if _, ok := data.(map[string]interface{}); ok {
		// Generate grid asynchronously
		allChunks := w.ChunkManager.GetAllChunks()
		result := w.generateGridDataAsync(allChunks)
		if callback != nil {
			callback(result)
		}
	}
}

func (w *World) processEntityUpdate(data interface{}, callback func(interface{})) {
	if entities, ok := data.([]coretypes.Entity); ok {
		// Update entities in parallel
		var wg sync.WaitGroup
		for _, e := range entities {
			wg.Add(1)
			go func(ent coretypes.Entity) {
				defer wg.Done()
				ent.Update()
			}(e)
		}
		wg.Wait()
		if callback != nil {
			callback(nil)
		}
	}
}

func (w *World) processPhysicsUpdate(data interface{}, callback func(interface{})) {
	if _, ok := data.(map[string]interface{}); ok {
		// Process physics updates asynchronously
		if callback != nil {
			callback(nil)
		}
	}
}
