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
	time.Sleep(100 * time.Millisecond) // Small delay to make progress visible

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
	time.Sleep(100 * time.Millisecond)

	// Always reset world generation with the new seed
	generation.ResetWorldGeneration(seed)
	progress.UpdateCurrentStepProgress(3, "Reset generation systems")
	time.Sleep(100 * time.Millisecond)

	// Pre-allocate player image to avoid recreating it every frame
	g.playerImage = ebiten.NewImage(settings.PlayerWidth, settings.PlayerHeight)
	g.playerImage.Fill(color.RGBA{255, 255, 0, 255}) // Yellow
	progress.UpdateCurrentStepProgress(4, "Created player graphics")
	time.Sleep(100 * time.Millisecond)

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

	// Clear screen with black background
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Black

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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.LastScreenW = outsideWidth
	g.LastScreenH = outsideHeight
	return outsideWidth, outsideHeight
}

// handleBlockInteraction processes block interaction events from the player
func (g *Game) handleBlockInteraction(p *player.Player, interaction *player.BlockInteraction) {
	switch interaction.Type {
	case player.BreakBlock:
		g.World.BreakBlock(interaction.BlockX, interaction.BlockY)
	case player.PlaceBlock:
		g.World.PlaceBlock(interaction.BlockX, interaction.BlockY, p.SelectedBlock)
	}
}

// drawCrosshair draws a targeting reticle and highlights the block under the cursor
func (g *Game) drawCrosshair(screen *ebiten.Image) {
	if len(g.World.Entities) == 0 {
		return
	}

	player, ok := g.World.Entities[0].(*player.Player)
	if !ok {
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()

	// Convert screen coordinates to world coordinates
	worldX := float64(mouseX) + g.CameraX
	worldY := float64(mouseY) + g.CameraY

	// Convert to block coordinates
	blockX := int(worldX / float64(settings.TileSize))
	blockY := int(worldY / float64(settings.TileSize))

	// Handle negative coordinates properly
	if worldX < 0 {
		blockX = int(worldX/float64(settings.TileSize)) - 1
	}
	if worldY < 0 {
		blockY = int(worldY/float64(settings.TileSize)) - 1
	}

	// Check if block is in range
	playerCenterX := player.X + float64(settings.PlayerColliderWidth)/2
	playerCenterY := player.Y + float64(settings.PlayerColliderHeight)/2
	blockCenterX := float64(blockX)*float64(settings.TileSize) + float64(settings.TileSize)/2
	blockCenterY := float64(blockY)*float64(settings.TileSize) + float64(settings.TileSize)/2

	dx := blockCenterX - playerCenterX
	dy := blockCenterY - playerCenterY
	distance := dx*dx + dy*dy
	inRange := distance <= player.InteractionRange*player.InteractionRange

	// Calculate screen position of the target block
	blockScreenX := float64(blockX*settings.TileSize) - g.CameraX
	blockScreenY := float64(blockY*settings.TileSize) - g.CameraY

	// Only draw if block is on screen
	if blockScreenX >= -float64(settings.TileSize) && blockScreenX < float64(screen.Bounds().Dx()) &&
		blockScreenY >= -float64(settings.TileSize) && blockScreenY < float64(screen.Bounds().Dy()) {

		// Create highlight color based on whether block is in range
		var highlightColor color.RGBA
		if inRange {
			highlightColor = color.RGBA{255, 255, 255, 128} // White semi-transparent
		} else {
			highlightColor = color.RGBA{255, 0, 0, 128} // Red semi-transparent (out of range)
		}

		// Draw block outline
		g.drawBlockOutline(screen, int(blockScreenX), int(blockScreenY), highlightColor)
	}

	// Draw simple crosshair at cursor
	crosshairSize := 8
	crosshairColor := color.RGBA{255, 255, 255, 200}

	// Horizontal line
	for i := -crosshairSize; i <= crosshairSize; i++ {
		px, py := mouseX+i, mouseY
		if px >= 0 && px < screen.Bounds().Dx() && py >= 0 && py < screen.Bounds().Dy() {
			screen.Set(px, py, crosshairColor)
		}
	}

	// Vertical line
	for i := -crosshairSize; i <= crosshairSize; i++ {
		px, py := mouseX, mouseY+i
		if px >= 0 && px < screen.Bounds().Dx() && py >= 0 && py < screen.Bounds().Dy() {
			screen.Set(px, py, crosshairColor)
		}
	}
}

// drawBlockOutline draws an outline around a block
func (g *Game) drawBlockOutline(screen *ebiten.Image, x, y int, outlineColor color.RGBA) {
	tileSize := settings.TileSize

	// Top edge
	for i := 0; i < tileSize; i++ {
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y >= 0 && y < screen.Bounds().Dy() {
			screen.Set(x+i, y, outlineColor)
		}
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y+1 >= 0 && y+1 < screen.Bounds().Dy() {
			screen.Set(x+i, y+1, outlineColor)
		}
	}

	// Bottom edge
	for i := 0; i < tileSize; i++ {
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y+tileSize-1 >= 0 && y+tileSize-1 < screen.Bounds().Dy() {
			screen.Set(x+i, y+tileSize-1, outlineColor)
		}
		if x+i >= 0 && x+i < screen.Bounds().Dx() && y+tileSize-2 >= 0 && y+tileSize-2 < screen.Bounds().Dy() {
			screen.Set(x+i, y+tileSize-2, outlineColor)
		}
	}

	// Left edge
	for i := 0; i < tileSize; i++ {
		if x >= 0 && x < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x, y+i, outlineColor)
		}
		if x+1 >= 0 && x+1 < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x+1, y+i, outlineColor)
		}
	}

	// Right edge
	for i := 0; i < tileSize; i++ {
		if x+tileSize-1 >= 0 && x+tileSize-1 < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x+tileSize-1, y+i, outlineColor)
		}
		if x+tileSize-2 >= 0 && x+tileSize-2 < screen.Bounds().Dx() && y+i >= 0 && y+i < screen.Bounds().Dy() {
			screen.Set(x+tileSize-2, y+i, outlineColor)
		}
	}
}
