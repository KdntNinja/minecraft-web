package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/loading"
	"github.com/KdntNinja/webcraft/engine/player"
	"github.com/KdntNinja/webcraft/engine/render"
	"github.com/KdntNinja/webcraft/engine/world"
)

type Game struct {
	World          *world.World
	LoadedChunks   map[string]bool        // Track loaded chunks for efficient management
	CenterChunkX   int                    // Current center chunk X coordinate
	CenterChunkY   int                    // Current center chunk Y coordinate
	LastScreenW    int                    // Cache last screen width
	LastScreenH    int                    // Cache last screen height
	ChunksPerFrame int                    // Limit chunks generated per frame
	CameraX        float64                // Camera X position
	CameraY        float64                // Camera Y position
	LoadDistance   int                    // Distance in chunks to load around player
	UnloadDistance int                    // Distance in chunks to unload from player
	LoadingScreen  *loading.LoadingScreen // Loading screen for initial setup
	IsLoading      bool                   // Whether we're still in loading phase
}

func NewGame() *Game {
	g := &Game{
		LoadedChunks:   make(map[string]bool),
		CenterChunkX:   0,
		CenterChunkY:   0,
		LastScreenW:    800, // Default screen width
		LastScreenH:    600, // Default screen height
		ChunksPerFrame: 1,   // Generate only 1 chunk per frame for smooth performance
		LoadDistance:   3,   // Load chunks within 3 chunk radius (much smaller)
		UnloadDistance: 6,   // Unload chunks beyond 6 chunk radius
		LoadingScreen:  loading.NewLoadingScreen(),
		IsLoading:      true,
	}

	return g
}

// markInitialChunksAsLoaded marks the initially generated chunks as loaded
func (g *Game) markInitialChunksAsLoaded() {
	if g.World == nil {
		return
	}
	for y := 0; y < len(g.World.Blocks); y++ {
		for x := 0; x < len(g.World.Blocks[y]); x++ {
			chunkX, chunkY := g.arrayToWorldCoords(x, y)
			g.LoadedChunks[chunkKey(chunkX, chunkY)] = true
		}
	}
}

func (g *Game) Update() error {
	// Handle loading phase
	if g.IsLoading {
		g.LoadingScreen.Update()
		if g.LoadingScreen.IsComplete() {
			// Loading complete, transfer world and setup camera
			g.World = g.LoadingScreen.GetWorld()
			g.IsLoading = false

			// Initialize camera position based on player position
			if len(g.World.Entities) > 0 {
				if player, ok := g.World.Entities[0].(*player.Player); ok {
					g.CameraX = player.X + float64(player.Width)/2 - float64(g.LastScreenW)/2
					g.CameraY = player.Y + float64(player.Height)/2 - float64(g.LastScreenH)/2
					g.CenterChunkX = int(player.X) / (block.ChunkWidth * block.TileSize)
					g.CenterChunkY = int(player.Y) / (block.ChunkHeight * block.TileSize)
				}
			}

			// Mark initial chunks as loaded
			g.markInitialChunksAsLoaded()
		}
		return nil
	}

	// Normal game update logic
	if g.World == nil {
		return nil
	}

	// Update all entities (including player)
	grid := g.World.ToIntGrid()
	for _, e := range g.World.Entities {
		e.Update()
		if p, ok := e.(interface {
			CollideBlocks([][]int)
		}); ok {
			p.CollideBlocks(grid)
		}
	}

	// Update camera to follow player and manage chunk loading
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			// Smooth camera following with some easing
			targetCameraX := player.X + float64(player.Width)/2 - float64(g.LastScreenW)/2
			targetCameraY := player.Y + float64(player.Height)/2 - float64(g.LastScreenH)/2

			// Smooth camera movement (lerp)
			lerpFactor := 0.1
			g.CameraX += (targetCameraX - g.CameraX) * lerpFactor
			g.CameraY += (targetCameraY - g.CameraY) * lerpFactor

			// Update center chunk based on player position
			newCenterChunkX := int(player.X) / (block.ChunkWidth * block.TileSize)
			newCenterChunkY := int(player.Y) / (block.ChunkHeight * block.TileSize)

			// Load new chunks if player moved to a different chunk
			if newCenterChunkX != g.CenterChunkX || newCenterChunkY != g.CenterChunkY {
				g.CenterChunkX = newCenterChunkX
				g.CenterChunkY = newCenterChunkY
			}

			// Always try to load chunks gradually (limited per frame)
			g.loadChunksAroundPlayer()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Show loading screen while loading
	if g.IsLoading {
		g.LoadingScreen.Draw(screen)
		return
	}

	// Normal game rendering
	if g.World == nil {
		return
	}

	// Apply camera transform
	var cameraOp ebiten.DrawImageOptions
	cameraOp.GeoM.Translate(-g.CameraX, -g.CameraY)

	// Create a temporary image for world rendering
	worldImg := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	// Render world with camera offset
	render.DrawWithCamera(&g.World.Blocks, worldImg, g.CameraX, g.CameraY)

	// Draw all entities with camera offset
	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			playerColor := [4]uint8{255, 255, 0, 255} // Yellow
			px, py := int(p.X-g.CameraX), int(p.Y-g.CameraY)
			img := ebiten.NewImage(player.Width, player.Height)
			img.Fill(color.RGBA{playerColor[0], playerColor[1], playerColor[2], playerColor[3]})
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(px), float64(py))
			worldImg.DrawImage(img, &op)
		}
	}

	// Draw the world image to the screen
	screen.DrawImage(worldImg, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.LastScreenW = outsideWidth
	g.LastScreenH = outsideHeight
	return outsideWidth, outsideHeight
}
