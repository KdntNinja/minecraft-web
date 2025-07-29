package world

func (w *World) startUpdateWorkers() {
	for i := 0; i < w.numUpdateWorkers; i++ {
		w.updateWorkers.Add(1)
		go w.updateWorker(i)
	}
}

func (w *World) updateWorker(workerID int) {
	defer w.updateWorkers.Done()

	for {
		select {
		case <-w.updateCtx.Done():
			return
		case task := <-w.updateTasks:
			// Process different types of update tasks
			switch task.taskType {
			case "grid_generation":
				w.processGridGeneration(task.data, task.callback)
			case "entity_update":
				w.processEntityUpdate(task.data, task.callback)
			case "physics_update":
				w.processPhysicsUpdate(task.data, task.callback)
			}
		}
	}
}
