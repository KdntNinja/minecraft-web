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
	World      *world.World
	LoadedY    int // Number of vertical chunk rows loaded
	MaxChunksY int // Total vertical chunk rows to eventually load
	CenterX    int // Center chunk X
}

func NewGame() *Game {
	maxChunksY := block.ChunkHeight * 4 // Arbitrary, can be infinite
	centerChunkX := 0
	initialChunksY := 2 // Only load a few rows at first for fast startup
	w := world.NewWorld(initialChunksY, centerChunkX)
	return &Game{
		World:      w,
		LoadedY:    initialChunksY,
		MaxChunksY: maxChunksY,
		CenterX:    centerChunkX,
	}
}

func (g *Game) Update() error {
	if g.LoadedY < g.MaxChunksY {
		// Add a new row of chunks to the world
		newRow := make([]world.Chunk, len(g.World.Blocks[0]))
		for cx := range newRow {
			newRow[cx] = world.GenerateChunk(g.CenterX+cx-len(newRow)/2, g.LoadedY)
		}
		g.World.Blocks = append(g.World.Blocks, newRow)
		g.LoadedY++
	}
	// Update all entities (including player)
	grid := g.World.ToIntGrid()
	for _, e := range g.World.Entities {
		e.Update()
		if p, ok := e.(interface {
			CollideBlocks([][]int)
			ClampX(float64, float64)
		}); ok {
			p.CollideBlocks(grid)
			p.ClampX(0, float64(len(grid[0])*block.TileSize-player.Width))
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	render.Draw(&g.World.Blocks, screen)
	// Draw all entities (player, etc.)
	for _, e := range g.World.Entities {
		if p, ok := e.(*player.Player); ok {
			playerColor := [4]uint8{255, 255, 0, 255} // Yellow
			px, py := int(p.X), int(p.Y)
			img := ebiten.NewImage(player.Width, player.Height)
			img.Fill(color.RGBA{playerColor[0], playerColor[1], playerColor[2], playerColor[3]})
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(px), float64(py))
			screen.DrawImage(img, &op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Calculate how many chunks are needed to fill the view
	chunksX := (outsideWidth + block.ChunkWidth*block.TileSize - 1) / (block.ChunkWidth * block.TileSize)
	chunksY := (outsideHeight + block.ChunkHeight*block.TileSize - 1) / (block.ChunkHeight * block.TileSize)

	// Ensure enough vertical chunks
	if len(g.World.Blocks) < chunksY {
		for i := len(g.World.Blocks); i < chunksY; i++ {
			newRow := make([]world.Chunk, len(g.World.Blocks[0]))
			for cx := range newRow {
				newRow[cx] = world.GenerateChunk(g.CenterX+cx-len(newRow)/2, i)
			}
			g.World.Blocks = append(g.World.Blocks, newRow)
			g.LoadedY++
		}
	}

	// Ensure enough horizontal chunks in each row
	for y := 0; y < len(g.World.Blocks); y++ {
		row := g.World.Blocks[y]
		if len(row) < chunksX {
			newRow := make([]world.Chunk, chunksX)
			copy(newRow, row)
			for x := len(row); x < chunksX; x++ {
				newRow[x] = world.GenerateChunk(g.CenterX+x-len(row)/2, y)
			}
			g.World.Blocks[y] = newRow
		}
	}

	return outsideWidth, outsideHeight
}
