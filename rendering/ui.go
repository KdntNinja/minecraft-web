package rendering

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// DrawGameUI draws the FPS, selected block, and controls UI.
// DrawGameUI draws the FPS, selected block, and controls UI.
func DrawGameUI(screen *ebiten.Image, fps float64, selectedBlock string) {
	// Draw FPS
	fpsText := fmt.Sprintf("FPS: %.1f", fps)
	ebitenutil.DebugPrintAt(screen, fpsText, 16, 12)

	// Draw selected block
	if selectedBlock != "" {
		blockText := fmt.Sprintf("Selected Block: %s", selectedBlock)
		ebitenutil.DebugPrintAt(screen, blockText, 16, 36)
	}

	// Draw controls
	controls := []string{
		"Controls:",
		"  Left Click  = Break",
		"  Right Click = Place",
		"  1-9,0       = Select Block",
	}
	for i, line := range controls {
		ebitenutil.DebugPrintAt(screen, line, 16, 64+i*18)
	}

	// Draw block numbers
	blocks := []string{
		"Blocks:",
		"  1 = Grass   2 = Dirt   3 = Clay   4 = Stone",
		"  5 = Copper  6 = Iron   7 = Gold   8 = Ash",
		"  9 = Wood    0 = Leaves",
	}
	for i, line := range blocks {
		ebitenutil.DebugPrintAt(screen, line, 16, 140+i*16)
	}
}
