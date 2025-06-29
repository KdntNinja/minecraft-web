package game

import (
	"crypto/rand"
	"fmt"
	"image/color"
	"math/big"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/progress"
	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
	"github.com/KdntNinja/webcraft/internal/gameplay/world"
	"github.com/KdntNinja/webcraft/internal/generation"
	"github.com/KdntNinja/webcraft/internal/rendering"
)

var globalSeed int64

// init generates a random seed when the package is initialized
func init() {
	// Generate a truly random seed
	if randomBig, err := rand.Int(rand.Reader, big.NewInt(1000000)); err == nil {
		globalSeed = randomBig.Int64()
	} else {
		// Fallback to time-based seed if crypto/rand fails
		globalSeed = time.Now().UnixNano() % 1000000
	}
}

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
	fpsHistory []float64 // For FPS graph in debug overlay
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

	// Hide the cursor for better gameplay experience
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	progress.UpdateCurrentStepProgress(2, "Set up game configuration")

	// Always reset world generation with the new seed
	generation.ResetWorldGeneration(seed)
	progress.UpdateCurrentStepProgress(3, "Reset generation systems")

	// Pre-allocate player image to avoid recreating it every frame
	g.playerImage = ebiten.NewImage(settings.PlayerWidth, settings.PlayerHeight)
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
			g.CameraX = player.X + float64(settings.PlayerWidth)/2 - float64(g.LastScreenW)/2
			g.CameraY = player.Y + float64(settings.PlayerHeight)/2 - float64(g.LastScreenH)/2 - float64(settings.TileSize*2) // Offset upward slightly
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

	// Pre-calculate camera bounds once
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	// Update entities (reduce slice allocation by reusing)
	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			// Frustum culling for entities
			if p.X+float64(settings.PlayerColliderWidth) < camLeft || p.X > camRight ||
				p.Y+float64(settings.PlayerColliderHeight) < camTop || p.Y > camBottom {
				continue // Skip entities far from view
			}
			// Set the offset for the player's collision system
			p.AABB.GridOffsetX = g.physicsOffsetX
			p.AABB.GridOffsetY = g.physicsOffsetY

			// Update player movement
			p.Update()

			// Handle block interactions separately
			blockInteraction := p.HandleBlockInteractions(g.CameraX, g.CameraY)
			if blockInteraction != nil {
				g.handleBlockInteraction(p, blockInteraction)
			}

			// Use cached physics world
			if g.physicsWorld != nil {
				p.CollideBlocksAdvanced(g.physicsWorld)
			}
		}
	}

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
		// Fill with black background
		screen.Fill(color.RGBA{0, 0, 0, 255}) // Black
		return
	}

	// --- Terraria-like sky color based on player Y position ---
	// Get player Y (if available), else use camera Y
	playerY := g.CameraY
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			playerY = player.Y
		}
	}
	// Use settings for sky transition
	skyColor := color.RGBA{135, 206, 235, 255}      // Terraria-like sky blue
	undergroundColor := color.RGBA{10, 10, 30, 255} // Deep blue/black

	var bgColor color.RGBA
	if playerY <= settings.SkyTransitionStartY {
		bgColor = skyColor
	} else if playerY >= settings.SkyTransitionEndY {
		bgColor = undergroundColor
	} else {
		t := (playerY - settings.SkyTransitionStartY) / (settings.SkyTransitionEndY - settings.SkyTransitionStartY)
		bgR := uint8(float64(skyColor.R)*(1-t) + float64(undergroundColor.R)*t)
		bgG := uint8(float64(skyColor.G)*(1-t) + float64(undergroundColor.G)*t)
		bgB := uint8(float64(skyColor.B)*(1-t) + float64(undergroundColor.B)*t)
		bgColor = color.RGBA{bgR, bgG, bgB, 255}
	}

	screen.Fill(bgColor)

	// Render world directly to screen (avoid intermediate image allocation)
	rendering.DrawWithCamera(g.World.GetChunksForRendering(), screen, g.CameraX, g.CameraY)

	// Pre-calculate camera bounds for entity culling
	camLeft := g.CameraX - float64(settings.TileSize*2)
	camRight := g.CameraX + float64(g.LastScreenW) + float64(settings.TileSize*2)
	camTop := g.CameraY - float64(settings.TileSize*2)
	camBottom := g.CameraY + float64(g.LastScreenH) + float64(settings.TileSize*2)

	// Reusable draw options to reduce allocations
	var op ebiten.DrawImageOptions

	// Draw only entities near the camera/screen for performance
	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			// Use pre-calculated camera bounds
			if p.X+float64(settings.PlayerColliderWidth) < camLeft || p.X > camRight ||
				p.Y+float64(settings.PlayerColliderHeight) < camTop || p.Y > camBottom {
				continue // Skip entities far from view
			}
			// Center sprite horizontally and bottom-align vertically over collider
			px := int(p.X - g.CameraX - float64(settings.PlayerWidth-settings.PlayerColliderWidth)/2)
			py := int(p.Y - g.CameraY - float64(settings.PlayerHeight-settings.PlayerColliderHeight))
			if px > -settings.PlayerWidth && px < g.LastScreenW && py > -settings.PlayerHeight && py < g.LastScreenH {
				// Reuse the DrawImageOptions instead of creating new ones
				op.GeoM.Reset()
				op.GeoM.Translate(float64(px), float64(py))
				screen.DrawImage(g.playerImage, &op)
			}
		}
	}

	// Draw crosshair/target indicator
	g.drawCrosshair(screen)

	// --- F3 debug screen ---
	if g.ShowDebug {
		g.drawDebugInfo(screen)
		return // Hide all other UI except debug overlay
	}

	// Display FPS in the top left corner
	fpsText := fmt.Sprintf("FPS: %.1f", g.currentFPS)
	ebitenutil.DebugPrint(screen, fpsText)

	// Display currently selected block and controls
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			selectedBlockText := fmt.Sprintf("\nSelected Block: %v", player.SelectedBlock)
			controlsText := "\nControls:\n Left Click = Break\n Right Click = Place"
			numbersText := "\nBlocks:\n 1=Grass\n 2=Dirt\n 3=Clay\n 4=Stone\n 5=Copper\n 6=Iron\n 7=Gold\n 8=Ash\n 9=Wood\n 0=Leaves\n"
			uiText := fpsText + selectedBlockText + controlsText + numbersText
			ebitenutil.DebugPrint(screen, uiText)
		}
	} else {
		ebitenutil.DebugPrint(screen, fpsText)
	}
}

// GetFPSHistory returns a slice of the last n FPS values for the debug graph
func (g *Game) GetFPSHistory(n int) []float64 {
	if len(g.fpsHistory) < n {
		// Pad with zeros if not enough data
		pad := make([]float64, n-len(g.fpsHistory))
		return append(pad, g.fpsHistory...)
	}
	return g.fpsHistory[len(g.fpsHistory)-n:]
}
