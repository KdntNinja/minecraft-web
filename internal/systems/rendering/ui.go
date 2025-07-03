package rendering

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// DrawGameUI draws the FPS, selected block, and controls UI.
// DrawGameUI draws the FPS, selected block, and controls UI.
func DrawGameUI(screen *ebiten.Image, fps float64, selectedBlock string) {
	fpsText := fmt.Sprintf("FPS: %.1f", fps)
	if selectedBlock != "" {
		selectedBlockText := fmt.Sprintf("\nSelected Block: %s", selectedBlock)
		controlsText := "\nControls:\n Left Click = Break\n Right Click = Place"
		numbersText := "\nBlocks:\n 1=Grass\n 2=Dirt\n 3=Clay\n 4=Stone\n 5=Copper\n 6=Iron\n 7=Gold\n 8=Ash\n 9=Wood\n 0=Leaves\n"
		uiText := fpsText + selectedBlockText + controlsText + numbersText
		ebitenutil.DebugPrint(screen, uiText)
		return
	}
	ebitenutil.DebugPrint(screen, fpsText)
}
