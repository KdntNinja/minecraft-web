package game

import (
	"fmt"
	"image/color"
	"runtime"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/container"
	"github.com/ebitenui/ebitenui/theme/basic"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KdntNinja/webcraft/internal/core/settings"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
)

var debugUI *ebitenui.UI
var debugPanel *container.Container
var debugLabels []*widget.Label
var debugFontFace = basic.DefaultFontFace

func (g *Game) initDebugUI() {
	if debugUI != nil {
		return
	}
	debugPanel = container.New(
		container.Layout(widget.NewAnchorLayout()),
		container.BackgroundImage(basic.NewNineSliceColor(color.RGBA{20, 20, 30, 220})),
	)
	debugUI = &ebitenui.UI{
		Container: debugPanel,
	}
	// Add labels for each debug line
	for i := 0; i < 10; i++ {
		lbl := widget.NewLabel(
			widget.LabelOpts.Text("", debugFontFace, &widget.LabelColor{Idle: color.White}),
		)
		debugPanel.AddChild(lbl)
		debugLabels = append(debugLabels, lbl)
	}
}

// drawDebugInfo draws a modern debug overlay using ebitenui
func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	g.initDebugUI()

	// --- FPS graph buffer ---
	if g.fpsHistory == nil {
		g.fpsHistory = make([]float64, settings.DebugGraphSamples)
	}
	g.fpsHistory = append(g.fpsHistory[1:], g.currentFPS)

	// --- Memory usage ---
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memMB := float64(m.Alloc) / 1024.0 / 1024.0

	// --- Player info ---
	var px, py, pvx, pvy float64
	var chunkX, chunkY int
	entityCount := len(g.World.Entities)
	if entityCount > 0 {
		if p, ok := g.World.Entities[0].(*player.Player); ok {
			px, py = p.X, p.Y
			pvx, pvy = p.VX, p.VY
			chunkX = int(px) / (settings.ChunkWidth * settings.TileSize)
			chunkY = int(py) / (settings.ChunkHeight * settings.TileSize)
		}
	}

	// --- Compose debug lines ---
	lines := []string{
		fmt.Sprintf("[DEBUG - F3]"),
		fmt.Sprintf("FPS: %.1f | Entities: %d | Mem: %.1f MB", g.currentFPS, entityCount, memMB),
		fmt.Sprintf("Player: X=%.2f Y=%.2f  Chunk: %d,%d", px, py, chunkX, chunkY),
		fmt.Sprintf("Velocity: X=%.2f Y=%.2f", pvx, pvy),
		fmt.Sprintf("Camera: X=%.2f Y=%.2f", g.CameraX, g.CameraY),
		fmt.Sprintf("Loaded Chunks: %d", g.World.ChunkManager.GetLoadedChunkCount()),
		fmt.Sprintf("Seed: %d", g.Seed),
		"", // Spacer
		"FPS Graph:",
		"", // Graph will be drawn below
	}
	for i, lbl := range debugLabels {
		if i < len(lines) {
			lbl.Label = lines[i]
		} else {
			lbl.Label = ""
		}
	}

	debugUI.Update()
	debugUI.Draw(screen)

	// --- Draw FPS graph below the panel ---
	graphX := 24
	graphY := settings.DebugOverlayHeight + 8
	graphW := settings.DebugOverlayWidth - 48
	maxFPS := 120.0
	for i := 1; i < settings.DebugGraphSamples; i++ {
		x0 := graphX + (i-1)*graphW/settings.DebugGraphSamples
		y0 := graphY + settings.DebugGraphHeight - int(g.fpsHistory[i-1]/maxFPS*float64(settings.DebugGraphHeight))
		x1 := graphX + i*graphW/settings.DebugGraphSamples
		y1 := graphY + settings.DebugGraphHeight - int(g.fpsHistory[i]/maxFPS*float64(settings.DebugGraphHeight))
		col := color.RGBA{100, 220, 100, 255}
		if g.fpsHistory[i] < 30 {
			col = color.RGBA{220, 100, 100, 255}
		}
		for dx := x0; dx <= x1; dx++ {
			if y0 >= graphY && y0 < graphY+settings.DebugGraphHeight && dx >= graphX && dx < graphX+graphW {
				screen.Set(dx, y0, col)
			}
			if y1 >= graphY && y1 < graphY+settings.DebugGraphHeight && dx >= graphX && dx < graphX+graphW {
				screen.Set(dx, y1, col)
			}
		}
	}
}
