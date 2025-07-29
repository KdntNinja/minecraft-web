package engine

import (
	"fmt"
	"image/color"
	"runtime"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/gameplay"
	"github.com/KdntNinja/webcraft/gameplay/world"
	"github.com/KdntNinja/webcraft/generation"
	"github.com/KdntNinja/webcraft/physics"
	"github.com/KdntNinja/webcraft/progress"
	"github.com/KdntNinja/webcraft/rendering"
	"github.com/KdntNinja/webcraft/settings"
	"github.com/KdntNinja/webcraft/worldgen"
)

type Game struct {
	World       *world.World
	LastScreenW int     // Cache last screen width
	LastScreenH int     // Cache last screen height
	CameraX     float64 // Camera X position
	CameraY     float64 // Camera Y position
	Seed        int64   // World seed for deterministic generation

	// Pre-allocated images to reduce memory allocation
	playerImage *ebiten.Image
	frameCount  int // For frame rate limiting

	// Performance monitoring
	fpsCounter    int       // Frame counter for FPS calculation
	lastFPSUpdate time.Time // Last time FPS was calculated
	currentFPS    float64   // Current FPS value to display

	// --- Physics cache ---
	physicsWorld   *physics.PhysicsWorld // Cached physics world for collisions
	physicsGrid    [][]int               // Cached grid used to build physicsWorld
	physicsOffsetX int
	physicsOffsetY int

	// Async physics system
	asyncPhysics *physics.AsyncPhysicsSystem

	// Debug
	ShowDebug     bool // Show debug screen when F3 is pressed
	prevF3Pressed bool // Track previous F3 key state for toggle

	// Debug graph
	fpsHistory    []float64 // For FPS graph in debug overlay
	fpsHistoryMin float64   // Track min FPS seen for relative graph
	fpsHistoryMax float64   // Track max FPS seen for relative graph

	// Additional debug tracking
	tickTimes   []float64 // Track frame/tick times for performance monitoring
	tickTimeMin float64   // Minimum tick time
	tickTimeMax float64   // Maximum tick time

	// Performance optimization
	frameStartTime time.Time
	updateTime     time.Duration
	renderTime     time.Duration
	parallelTasks  sync.WaitGroup
}

func NewGame() *Game {
	// Use view distance from settings
	viewDistance := settings.ChunkViewDistance
	totalChunks := (viewDistance*2 + 1) * (viewDistance*2 + 1)

	steps := []progress.ProgressStep{
		{Name: "Initializing", Weight: 1.0, SubSteps: 6, Description: "Starting game initialization..."},
		{Name: "World Setup", Weight: 1.0, SubSteps: 3, Description: "Setting up world structure..."},
		{Name: "Generating Terrain", Weight: 8.0, SubSteps: totalChunks, Description: "Generating world chunks..."},
		{Name: "Spawning Player", Weight: 1.0, SubSteps: 3, Description: "Creating player entity..."},
		{Name: "Finalizing", Weight: 1.0, SubSteps: 1, Description: "Finishing initialization..."},
	}
	progress.InitializeProgress(steps)

	seed := globalSeed
	progress.UpdateCurrentStepProgress(1, fmt.Sprintf("Using random world seed: %d", seed))

	g := &Game{
		LastScreenW:    800, // Default screen width
		LastScreenH:    600, // Default screen height
		Seed:           seed,
		lastFPSUpdate:  time.Now(), // Initialize FPS tracking
		currentFPS:     60.0,       // Default FPS value
		frameStartTime: time.Now(),
	}

	// Initialize async systems
	progress.UpdateCurrentStepProgress(2, "Initializing async physics system...")
	g.asyncPhysics = physics.GetAsyncPhysicsSystem()

	// Hide the cursor for better gameplay experience and use custom crosshair
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	progress.UpdateCurrentStepProgress(4, "Set up game configuration")

	// Always reset world generation with the new seed
	generation.ResetWorldGeneration(seed)
	progress.UpdateCurrentStepProgress(5, "Reset generation systems")

	// Pre-allocate player image to avoid recreating it every frame
	g.playerImage = ebiten.NewImage(settings.PlayerSpriteWidth, settings.PlayerSpriteHeight)
	g.playerImage.Fill(color.RGBA{255, 255, 0, 255}) // Yellow
	progress.UpdateCurrentStepProgress(6, "Created player graphics")

	// Complete initialization step
	progress.CompleteCurrentStep()

	// Create a simple world with fixed size, passing the seed
	// This will use the new progress system for world generation
	// Find a spawn point and create a chunk manager
	spawn := worldgen.FindSpawnPoint()
	chunkManager := generation.NewChunkManager(settings.ChunkViewDistance)
	// Center chunk manager on spawn location before world creation
	chunkManager.UpdatePlayerPosition(spawn.X, spawn.Y)
	g.World = world.NewWorld(seed, chunkManager, spawn)

	// Initialize camera position to follow the player's spawn location with tighter centering
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*gameplay.Player); ok {
			// Center camera more tightly on the player for a zoomed-in feel
			g.CameraX = player.X + float64(settings.PlayerSpriteWidth)/2 - float64(g.LastScreenW)/2
			g.CameraY = player.Y + float64(settings.PlayerSpriteHeight)/2 - float64(g.LastScreenH)/2 - float64(settings.TileSize*2) // Offset upward slightly
		}
	}

	runtime.GC() // Force garbage collection after initialization
	fmt.Printf("GAME: Initialized with %d CPU cores available\n", runtime.NumCPU())
	return g
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
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*gameplay.Player); ok {
			// Tighter camera following with offset for better view ahead
			targetCameraX := player.X + float64(settings.PlayerColliderWidth)/2 - float64(g.LastScreenW)/2
			targetCameraY := player.Y + float64(settings.PlayerColliderHeight)/2 - float64(g.LastScreenH)/2 - float64(settings.TileSize*2)

			// More responsive camera movement for zoomed-in feel
			lerpFactor := 0.12 // Increased from 0.05 for more responsive following
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

	// Wait for all parallel tasks to complete
	g.parallelTasks.Wait()

	// Track update performance
	g.updateTime = time.Since(g.frameStartTime)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	renderStart := time.Now()

	if g.World == nil {
		screen.Fill(color.RGBA{0, 0, 0, 255})
		return
	}

	// Sky color
	playerY := g.CameraY
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*gameplay.Player); ok {
			playerY = player.Y
		}
	}
	bgColor := GetBackgroundColor(playerY)
	screen.Fill(bgColor)

	// World rendering using DrawWithCamera from renderer.go
	chunks, ok := g.World.GetChunksForRendering().(map[coretypes.ChunkCoord]*coretypes.Chunk)
	if ok {
		rendering.DrawWithCamera(chunks, screen, g.CameraX, g.CameraY)
	}

	// Entity rendering
	rendering.DrawEntities(g.World.Entities, screen, g.CameraX, g.CameraY, g.LastScreenW, g.LastScreenH, g.playerImage)

	// Crosshair
	rendering.DrawCrosshair(screen, g.World, g.CameraX, g.CameraY)

	if g.ShowDebug {
		// Hide normal UI and show debug overlay
		memStats := new(runtime.MemStats)
		runtime.ReadMemStats(memStats)
		memUsage := float64(memStats.Alloc) / (1024 * 1024)
		maxMem := float64(memStats.Sys) / (1024 * 1024)
		// Placeholder debug strings
		playerInfo := "Player: N/A"
		chunkInfo := "Chunks: N/A"
		playerStats := "Stats: N/A"
		camInfo := fmt.Sprintf("Camera: (%.1f, %.1f)", g.CameraX, g.CameraY)
		seedInfo := fmt.Sprintf("Seed: %d", g.Seed)
		worldInfo := "World: N/A"
		gcPercent := float64(memStats.GCCPUFraction) * 100
		// Use empty slices for block metrics
		renderedBlocksHistory := []int{}
		generatedBlocksHistory := []int{}
		loadedChunks := 0 // Could not determine loaded chunk count
		rendering.DrawDebugOverlay(
			screen,
			g.fpsHistory,
			g.fpsHistoryMin, g.fpsHistoryMax,
			g.currentFPS,
			loadedChunks,
			len(g.World.Entities),
			memUsage, maxMem,
			playerInfo, chunkInfo, playerStats, camInfo, seedInfo, worldInfo,
			g.tickTimes, g.tickTimeMin, g.tickTimeMax,
			gcPercent,
			renderedBlocksHistory,
			generatedBlocksHistory,
		)
	} else {
		// Hotbar UI (top left)
		if len(g.World.Entities) > 0 {
			if player, ok := g.World.Entities[0].(*gameplay.Player); ok {
				rendering.DrawHotbarUI(screen, player)
			}
		}
	}

	// Track render performance
	g.renderTime = time.Since(renderStart)

	// Update performance tracking
	g.updatePerformanceMetrics()
}

// UpdateEntitiesNearCameraAsync updates entities using async physics system
func (g *Game) UpdateEntitiesNearCameraAsync() {
	// Pre-calculate camera bounds once
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	// Filter entities within camera bounds
	var nearbyEntities []coretypes.Entity
	for _, e := range g.World.Entities {
		if p, ok := e.(*gameplay.Player); ok {
			// Frustum culling for entities
			if p.X+float64(settings.PlayerColliderWidth) < camLeft || p.X > camRight ||
				p.Y+float64(settings.PlayerColliderHeight) < camTop || p.Y > camBottom {
				continue // Skip entities far from view
			}
			nearbyEntities = append(nearbyEntities, e)
		}
	}

	// Process entities using async physics system
	g.asyncPhysics.ProcessEntitiesAsync(nearbyEntities, g.physicsWorld, func(ent coretypes.Entity) {
		if p, ok := ent.(*gameplay.Player); ok {
			// Set the offset for the player's collision system
			p.AABB.GridOffsetX = g.physicsOffsetX
			p.AABB.GridOffsetY = g.physicsOffsetY

			// Handle block interactions separately
			blockInteraction := p.HandleBlockInteractions(g.CameraX, g.CameraY)
			if blockInteraction != nil {
				switch blockInteraction.Type {
				case gameplay.BreakBlock:
					// Get the block type before breaking
					blockType := g.World.GetBlockAt(blockInteraction.BlockX, blockInteraction.BlockY)
					if blockType != coretypes.Air {
						// Add block to inventory
						p.AddToInventory(blockType, 1)
						g.World.BreakBlock(blockInteraction.BlockX, blockInteraction.BlockY)
					}
				case gameplay.PlaceBlock:
					// Only place if player has block in inventory
					if p.Inventory[p.SelectedBlock] > 0 {
						placed := g.World.PlaceBlock(blockInteraction.BlockX, blockInteraction.BlockY, p.SelectedBlock)
						if placed {
							p.RemoveFromInventory(p.SelectedBlock, 1)
						}
					}
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
	g.asyncPhysics.UpdateSpatialGrid(g.World.Entities)

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
