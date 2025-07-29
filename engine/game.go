package engine

import (
	"fmt"
	"image/color"
	"runtime"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/generation"
	"github.com/KdntNinja/webcraft/rendering"
	"github.com/KdntNinja/webcraft/settings"
)

	// World rendering using DrawWithCamera from renderer.go
	rendering.DrawWithCamera(chunks, screen, g.CameraX, g.CameraY)
	// Entity rendering
	rendering.DrawEntities(entities, screen, g.CameraX, g.CameraY, g.LastScreenW, g.LastScreenH, g.playerImage)

	// Crosshair
	rendering.DrawCrosshair(screen, g.World, g.CameraX, g.CameraY)

	if g.ShowDebug {
		// Hide normal UI and show debug overlay
		rendering.DrawDebugOverlay(
			screen,
			g.fpsHistory,
			g.fpsHistoryMin, g.fpsHistoryMax,
			g.currentFPS,
			chunks,
			len(entities),
			0.0, 0.0, // memUsage, maxMem (implement as needed)
			"", "", "", "", "", "", // playerInfo, chunkInfo, playerStats, camInfo, seedInfo, worldInfo
			g.tickTimes, g.tickTimeMin, g.tickTimeMax,
			nil, nil, // renderedBlocksHistory, generatedBlocksHistory
		)
	} else {
		// Hotbar UI (top left)
		if len(entities) > 0 {
			rendering.DrawHotbarUI(screen, entities[0])
		}
	}

	// Track render performance
	g.renderTime = time.Since(renderStart)
}

func (g *Game) Update() error {
	g.frameStartTime = time.Now()

	if g.World == nil {
		return nil
	}

	// --- F3 debug toggle (edge-triggered) ---
	f3Pressed := ebiten.IsKeyPressed(ebiten.KeyF3)
	if f3Pressed && !g.prevF3Pressed {
		g.ShowDebug = !g.ShowDebug
	}
	g.prevF3Pressed = f3Pressed

	g.frameCount++
	g.fpsCounter++

	// Start parallel tasks
	g.parallelTasks.Add(1)
	go func() {
		defer g.parallelTasks.Done()
		// Update world (handles dynamic chunk loading)
		g.World.Update()
	}()

	// Update FPS calculation every second
	now := time.Now()
	if now.Sub(g.lastFPSUpdate) >= time.Second {
		g.currentFPS = float64(g.fpsCounter) / now.Sub(g.lastFPSUpdate).Seconds()
		g.fpsCounter = 0
		g.lastFPSUpdate = now
	}

	// Update entities using async physics system
	g.parallelTasks.Add(1)
	go func() {
		defer g.parallelTasks.Done()
		g.UpdateEntitiesNearCameraAsync()
	}()

	// Update camera to follow player more responsively for zoomed-in feel
	entities := g.World.GetEntities()
	if len(entities) > 0 {
		if player, ok := entities[0].(interface {
			GetX() float64
			GetY() float64
		}); ok {
			targetCameraX := player.GetX() + float64(settings.PlayerColliderWidth)/2 - float64(g.LastScreenW)/2
			targetCameraY := player.GetY() + float64(settings.PlayerColliderHeight)/2 - float64(g.LastScreenH)/2 - float64(settings.TileSize*2)

			lerpFactor := 0.12
			g.CameraX += (targetCameraX - g.CameraX) * lerpFactor
			g.CameraY += (targetCameraY - g.CameraY) * lerpFactor
		}
	}

	// Only regenerate collision grid and physics world when necessary
	if g.frameCount%60 == 0 || g.World.IsGridDirty() || g.physicsWorld == nil {
		g.parallelTasks.Add(1)
		go func() {
			defer g.parallelTasks.Done()
			g.updatePhysicsWorldAsync()
		}()
	}

	switch c := chunksIface.(type) {
	case map[coretypes.ChunkCoord]*coretypes.Chunk:
		chunks = c
	case map[generation.ChunkCoord]*coretypes.Chunk:
		chunks = make(map[coretypes.ChunkCoord]*coretypes.Chunk, len(c))
		for k, v := range c {
			chunks[coretypes.ChunkCoord{X: k.X, Y: k.Y}] = v
		}
	}
	var chunks map[coretypes.ChunkCoord]*coretypes.Chunk
	switch c := chunksIface.(type) {
	case map[coretypes.ChunkCoord]*coretypes.Chunk:
		chunks = c
	case map[generation.ChunkCoord]*coretypes.Chunk:
		chunks = make(map[coretypes.ChunkCoord]*coretypes.Chunk, len(c))
		for k, v := range c {
			chunks[coretypes.ChunkCoord{X: k.X, Y: k.Y}] = v
	if w, ok := g.World.(*world.World); ok {
		rendering.DrawCrosshair(screen, w, g.CameraX, g.CameraY)
	}
	}
	chunks, _ := chunksIface.(map[coretypes.ChunkCoord]*coretypes.Chunk)
			nil, nil, // renderedBlocksHistory, generatedBlocksHistory
		)
	} else {
		// Hotbar UI (top left)
		if len(entities) > 0 {
			rendering.DrawHotbarUI(screen, entities[0])
			len(chunks),
	if w, ok := g.World.(*world.World); ok {
		rendering.DrawCrosshair(screen, w, g.CameraX, g.CameraY)
	}
	worldPtr, _ := g.World.(*coretypes.World)
	rendering.DrawCrosshair(screen, worldPtr, g.CameraX, g.CameraY)
	// Track render performance
	g.renderTime = time.Since(renderStart)

	// Update performance tracking
			if p, ok := entities[0].(*gameplay.Player); ok {
				rendering.DrawHotbarUI(screen, p)
			}
}
			len(chunks),
// UpdateEntitiesNearCameraAsync updates entities using async physics system
			len(chunks),
	// Pre-calculate camera bounds once
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	// Filter entities within camera bounds
			if p, ok := entities[0].(*gameplay.Player); ok {
				rendering.DrawHotbarUI(screen, p)
			}
	for _, e := range g.World.GetEntities() {
		if p, ok := e.(interface {
			playerPtr, _ := entities[0].(*coretypes.Player)
			rendering.DrawHotbarUI(screen, playerPtr)
			GetY() float64
			GetColliderWidth() float64
			GetColliderHeight() float64
		}); ok {
			if p.GetX()+p.GetColliderWidth() < camLeft || p.GetX() > camRight ||
				p.GetY()+p.GetColliderHeight() < camTop || p.GetY() > camBottom {
				continue
			}
			nearbyEntities = append(nearbyEntities, e)
		}
	}

	// Process entities using async physics system
	g.asyncPhysics.ProcessEntitiesAsync(nearbyEntities, g.physicsWorld, func(ent coretypes.Entity) {
		// Use type assertion to access position
		if p, ok := ent.(interface {
			GetX() float64
			GetY() float64
		}); ok {
			x := int(p.GetX())
			y := int(p.GetY())
			if x >= 0 && x < len(g.physicsGrid[0]) && y >= 0 && y < len(g.physicsGrid) {
				block := g.physicsGrid[y][x]
				if block != 0 {
					fmt.Printf("[DEBUG] Entity at (%d,%d) collides with block type %d\n", x, y, block)
				}
			}
		}
	})
}

// updatePhysicsWorldAsync updates the physics world asynchronously
func (g *Game) updatePhysicsWorldAsync() {
	g.physicsGrid, g.physicsOffsetX, g.physicsOffsetY = g.World.ToIntGrid()
	g.physicsWorld = physics.NewPhysicsWorld(g.physicsGrid)

	// Update spatial grid for physics system
	g.asyncPhysics.UpdateSpatialGrid(g.World.GetEntities())

	// Sanity check: if grid is all air, log a warning (for debugging)
	allAir := true
	for _, row := range g.physicsGrid {
		for _, v := range row {
			if v != 0 {
				allAir = false
				break
			}
		}
		if !allAir {
			break
		}
	}
	if allAir {
		fmt.Println("[WARN] Physics grid is all air! Player will float.")
	}
}

// updatePerformanceMetrics updates performance tracking metrics
func (g *Game) updatePerformanceMetrics() {
	// Update tick times for debug overlay
	totalFrameTime := g.updateTime + g.renderTime
	g.tickTimes = append(g.tickTimes, totalFrameTime.Seconds()*1000) // Convert to milliseconds

	// Keep only last 120 frames for performance
	if len(g.tickTimes) > 120 {
		g.tickTimes = g.tickTimes[1:]
	}

	// Update min/max tick times
	if len(g.tickTimes) > 0 {
		frameTime := g.tickTimes[len(g.tickTimes)-1]
		if g.tickTimeMin == 0 || frameTime < g.tickTimeMin {
			g.tickTimeMin = frameTime
		}
		if frameTime > g.tickTimeMax {
			g.tickTimeMax = frameTime
		}
	}

	// Update FPS history for debug overlay
	g.fpsHistory = append(g.fpsHistory, g.currentFPS)
	if len(g.fpsHistory) > 120 {
		g.fpsHistory = g.fpsHistory[1:]
	}

	// Update FPS min/max
	if len(g.fpsHistory) > 0 {
		if g.fpsHistoryMin == 0 || g.currentFPS < g.fpsHistoryMin {
			g.fpsHistoryMin = g.currentFPS
		}
		if g.currentFPS > g.fpsHistoryMax {
			g.fpsHistoryMax = g.currentFPS
		}
	}
}

// Shutdown cleanly shuts down all async systems
func (g *Game) Shutdown() {
	fmt.Println("GAME: Shutting down async systems...")
	if g.World != nil {
		g.World.Stop()
	}
	if g.asyncPhysics != nil {
		g.asyncPhysics.Shutdown()
	}
	fmt.Println("GAME: Shutdown complete")
}
