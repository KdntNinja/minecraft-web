package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/player"
	"github.com/KdntNinja/webcraft/engine/render"
	"github.com/KdntNinja/webcraft/engine/world"
)

type Game struct {
	World          *world.World
	LoadedY        int // Number of vertical chunk rows loaded
	MaxChunksY     int // Total vertical chunk rows to eventually load
	CenterX        int // Center chunk X
	LastScreenW    int // Cache last screen width
	LastScreenH    int // Cache last screen height
	ChunksPerFrame int // Limit chunks generated per frame
	CameraX        float64 // Camera X position
	CameraY        float64 // Camera Y position
}

func NewGame() *Game {
	centerChunkX := 0
	initialChunksY := 10 // Start with enough chunks to fill most screens
	w := world.NewWorld(initialChunksY, centerChunkX)
	return &Game{
		World:          w,
		LoadedY:        initialChunksY,
		MaxChunksY:     0, // Will be set dynamically based on window height
		CenterX:        centerChunkX,
		LastScreenW:    0,
		LastScreenH:    0,
		ChunksPerFrame: 2, // Generate max 2 chunks per frame to reduce lag
	}
}

func (g *Game) Update() error {
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

	// Update camera to follow player
	if len(g.World.Entities) > 0 {
		if player, ok := g.World.Entities[0].(*player.Player); ok {
			// Smooth camera following with some easing
			targetCameraX := player.X + float64(player.Width)/2 - float64(g.LastScreenW)/2
			targetCameraY := player.Y + float64(player.Height)/2 - float64(g.LastScreenH)/2
			
			// Smooth camera movement (lerp)
			lerpFactor := 0.1
			g.CameraX += (targetCameraX - g.CameraX) * lerpFactor
			g.CameraY += (targetCameraY - g.CameraY) * lerpFactor
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
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
	// Only regenerate chunks if screen size changed significantly to avoid constant recalculation
	if abs(outsideWidth-g.LastScreenW) < 10 && abs(outsideHeight-g.LastScreenH) < 10 {
		return outsideWidth, outsideHeight
	}

	g.LastScreenW = outsideWidth
	g.LastScreenH = outsideHeight

	// Calculate required chunks to fill the entire screen
	chunkWidthPixels := block.ChunkWidth * block.TileSize
	chunkHeightPixels := block.ChunkHeight * block.TileSize

	chunksX := (outsideWidth + chunkWidthPixels - 1) / chunkWidthPixels
	chunksY := (outsideHeight + chunkHeightPixels - 1) / chunkHeightPixels

	// Ensure we have enough chunks to fill the screen completely
	currentChunksY := len(g.World.Blocks)
	currentChunksX := 0
	if currentChunksY > 0 {
		currentChunksX = len(g.World.Blocks[0])
	}

	// Limit chunk generation per frame to reduce lag spikes
	chunksGenerated := 0

	// Expand vertically if needed (limited per frame)
	if currentChunksY < chunksY && chunksGenerated < g.ChunksPerFrame {
		rowsToAdd := min(chunksY-currentChunksY, g.ChunksPerFrame-chunksGenerated)
		for i := 0; i < rowsToAdd; i++ {
			newRow := make([]world.Chunk, max(currentChunksX, chunksX))
			for cx := 0; cx < len(newRow); cx++ {
				newRow[cx] = world.GenerateChunk(g.CenterX+cx-len(newRow)/2, currentChunksY+i)
				chunksGenerated++
				if chunksGenerated >= g.ChunksPerFrame {
					break
				}
			}
			g.World.Blocks = append(g.World.Blocks, newRow)
			if chunksGenerated >= g.ChunksPerFrame {
				break
			}
		}
	}

	// Expand horizontally if needed - ensure ALL rows have the same width
	targetChunksX := max(currentChunksX, chunksX)
	if chunksGenerated < g.ChunksPerFrame && currentChunksX < targetChunksX {
		for y := 0; y < len(g.World.Blocks); y++ {
			row := g.World.Blocks[y]
			if len(row) < targetChunksX {
				newRow := make([]world.Chunk, targetChunksX)
				copy(newRow, row)
				// Generate missing chunks up to the limit
				chunksToGenerate := min(targetChunksX-len(row), g.ChunksPerFrame-chunksGenerated)
				for x := len(row); x < len(row)+chunksToGenerate; x++ {
					newRow[x] = world.GenerateChunk(g.CenterX+x-len(row)/2, y)
					chunksGenerated++
				}
				// Fill remaining slots with air chunks if we hit the generation limit
				for x := len(row)+chunksToGenerate; x < targetChunksX; x++ {
					newRow[x] = world.Chunk{} // Empty chunk (all air)
				}
				g.World.Blocks[y] = newRow
			}
		}
	}

	return outsideWidth, outsideHeight
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
