package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/progress"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
	"github.com/KdntNinja/webcraft/internal/gameplay/world"
	"github.com/KdntNinja/webcraft/internal/generation"
	"github.com/KdntNinja/webcraft/internal/rendering"
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
	physicsWorld   *entity.PhysicsWorld // Cached physics world for collisions
	physicsGrid    [][]int              // Cached grid used to build physicsWorld
	physicsOffsetX int
	physicsOffsetY int

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
}

func NewGame() *Game {
	// Use view distance from settings
	viewDistance := settings.ChunkViewDistance
	totalChunks := (viewDistance*2 + 1) * (viewDistance*2 + 1)

	steps := []progress.ProgressStep{
		{Name: "Initializing", Weight: 1.0, SubSteps: 4, Description: "Starting game initialization..."},
		{Name: "World Setup", Weight: 1.0, SubSteps: 3, Description: "Setting up world structure..."},
		{Name: "Generating Terrain", Weight: 8.0, SubSteps: totalChunks, Description: "Generating world chunks..."},
		{Name: "Spawning Player", Weight: 1.0, SubSteps: 3, Description: "Creating player entity..."},
		{Name: "Finalizing", Weight: 1.0, SubSteps: 1, Description: "Finishing initialization..."},
	}
	progress.InitializeProgress(steps)

	seed := globalSeed
	progress.UpdateCurrentStepProgress(1, fmt.Sprintf("Using random world seed: %d", seed))

	g := &Game{
		LastScreenW:   800, // Default screen width
		LastScreenH:   600, // Default screen height
		Seed:          seed,
		lastFPSUpdate: time.Now(), // Initialize FPS tracking
		currentFPS:    60.0,       // Default FPS value
	}

	// Hide the cursor for better gameplay experience and use custom crosshair
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	progress.UpdateCurrentStepProgress(2, "Set up game configuration")

	// Always reset world generation with the new seed
	generation.ResetWorldGeneration(seed)
	progress.UpdateCurrentStepProgress(3, "Reset generation systems")

	// Pre-allocate player image to avoid recreating it every frame
	g.playerImage = ebiten.NewImage(settings.PlayerSpriteWidth, settings.PlayerSpriteHeight)
	g.playerImage.Fill(color.RGBA{255, 255, 0, 255}) // Yellow
	progress.UpdateCurrentStepProgress(4, "Created player graphics")

	// Complete initialization step
	progress.CompleteCurrentStep()

	// Create a simple world with fixed size, passing the seed
	// This will use the new progress system for world generation
	g.World = world.NewWorld(seed)

	// Initialize camera position to follow the player's spawn location with tighter centering
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			// Center camera more tightly on the player for a zoomed-in feel
			g.CameraX = player.X + float64(settings.PlayerSpriteWidth)/2 - float64(g.LastScreenW)/2
			g.CameraY = player.Y + float64(settings.PlayerSpriteHeight)/2 - float64(g.LastScreenH)/2 - float64(settings.TileSize*2) // Offset upward slightly
		}
	}

	return g
}

func (g *Game) Update() error {
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

	// Update world (handles dynamic chunk loading)
	g.World.Update()

	// Update FPS calculation every second
	now := time.Now()
	if now.Sub(g.lastFPSUpdate) >= time.Second {
		g.currentFPS = float64(g.fpsCounter) / now.Sub(g.lastFPSUpdate).Seconds()
		g.fpsCounter = 0
		g.lastFPSUpdate = now
	}

	// Update only entities near the camera/screen - cached grid for better performance
	g.UpdateEntitiesNearCamera()

	// Update camera to follow player more responsively for zoomed-in feel
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
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
		g.physicsGrid, g.physicsOffsetX, g.physicsOffsetY = g.World.ToIntGrid()
		g.physicsWorld = entity.NewPhysicsWorld(g.physicsGrid)
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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.World == nil {
		screen.Fill(color.RGBA{0, 0, 0, 255})
		return
	}

	// Sky color
	playerY := g.CameraY
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			playerY = player.Y
		}
	}
	bgColor := GetBackgroundColor(playerY)
	screen.Fill(bgColor)

	// World rendering
	rendering.DrawWithCamera(g.World.GetChunksForRendering(), screen, g.CameraX, g.CameraY)

	// Entity rendering
	rendering.DrawEntities(g.World.Entities, screen, g.CameraX, g.CameraY, g.LastScreenW, g.LastScreenH, g.playerImage)

	// Crosshair
	rendering.DrawCrosshair(screen, g.World, g.CameraX, g.CameraY)

	// UI
	selectedBlock := ""
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			selectedBlock = player.SelectedBlock.String()
		}
	}
	rendering.DrawGameUI(screen, g.currentFPS, selectedBlock)
}
