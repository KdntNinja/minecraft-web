package game

import (
	"github.com/KdntNinja/webcraft/engine/block"
	"github.com/KdntNinja/webcraft/engine/render"
	"github.com/KdntNinja/webcraft/engine/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	World      [][]world.Chunk
	LoadedY    int // Number of vertical chunk rows loaded
	MaxChunksY int // Total vertical chunk rows to eventually load
	CenterX    int // Center chunk X
}

func NewGame() *Game {
	maxChunksY := block.ChunkHeight * 4 // Arbitrary, can be infinite
	centerChunkX := 0
	initialChunksY := 2 // Only load a few rows at first for fast startup
	w := world.GenerateWorld(initialChunksY, centerChunkX)
	return &Game{
		World:      w,
		LoadedY:    initialChunksY,
		MaxChunksY: maxChunksY,
		CenterX:    centerChunkX,
	}
}

func (g *Game) Update() error {
	// Load more chunks each frame until MaxChunksY is reached
	if g.LoadedY < g.MaxChunksY {
		g.World = append(g.World, world.GenerateWorld(1, g.CenterX)[0])
		g.LoadedY++
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	render.Draw(&g.World, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Calculate how many chunks are needed to fill the view
	chunksX := (outsideWidth + block.ChunkWidth*block.TileSize - 1) / (block.ChunkWidth * block.TileSize)
	chunksY := (outsideHeight + block.ChunkHeight*block.TileSize - 1) / (block.ChunkHeight * block.TileSize)

	// Ensure enough vertical chunks
	if len(g.World) < chunksY {
		for i := len(g.World); i < chunksY; i++ {
			g.World = append(g.World, world.GenerateWorld(1, g.CenterX)[0])
			g.LoadedY++
		}
	}

	// Ensure enough horizontal chunks in each row
	for y := 0; y < len(g.World); y++ {
		row := g.World[y]
		if len(row) < chunksX {
			newRow := make([]world.Chunk, chunksX)
			copy(newRow, row)
			for x := len(row); x < chunksX; x++ {
				newRow[x] = world.GenerateChunk(g.CenterX+x-len(row)/2, y)
			}
			g.World[y] = newRow
		}
	}

	return outsideWidth, outsideHeight
}
