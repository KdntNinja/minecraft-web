package rendering

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/gameplay/player"
)

// DrawHotbarUI draws a Minecraft-style hotbar in the top left of the screen
func DrawHotbarUI(screen *ebiten.Image, p *player.Player) {
	tileSize := 36   // UI tile size (pixels)
	padding := 10    // Padding between slots
	x0, y0 := 18, 18 // Top left corner

	hotbarWidth := 9*tileSize + 8*padding
	hotbarHeight := tileSize + 12

	// Draw faint background behind hotbar
	hotbarBg := ebiten.NewImage(hotbarWidth+16, hotbarHeight)
	hotbarBg.Fill(color.RGBA{20, 20, 20, 180})
	hotbarBgOpts := &ebiten.DrawImageOptions{}
	hotbarBgOpts.GeoM.Translate(float64(x0-8), float64(y0-6))
	screen.DrawImage(hotbarBg, hotbarBgOpts)

	for i := 0; i < len(p.Hotbar); i++ {
		x := x0 + i*(tileSize+padding)
		blockType := p.Hotbar[i]
		count := p.Inventory[blockType]

		// Draw slot with rounded corners and shadow
		slotImg := ebiten.NewImage(tileSize, tileSize)
		slotImg.Fill(color.RGBA{60, 60, 60, 220})
		if blockType == p.SelectedBlock {
			// Thicker, vibrant border for selected
			slotImg.Fill(color.RGBA{255, 215, 0, 220})
		}
		// Draw shadow
		shadowImg := ebiten.NewImage(tileSize, tileSize)
		shadowImg.Fill(color.RGBA{0, 0, 0, 60})
		shadowOpts := &ebiten.DrawImageOptions{}
		shadowOpts.GeoM.Translate(float64(x)+2, float64(y0)+4)
		screen.DrawImage(shadowImg, shadowOpts)

		// Draw slot
		slotOpts := &ebiten.DrawImageOptions{}
		slotOpts.GeoM.Translate(float64(x), float64(y0))
		screen.DrawImage(slotImg, slotOpts)

		// Draw block icon centered
		if blockType != block.Air {
			tile := getBlockTileImage(blockType)
			if tile != nil {
				iconOpts := &ebiten.DrawImageOptions{}
				iconOpts.GeoM.Translate(float64(x)+4, float64(y0)+4)
				screen.DrawImage(tile, iconOpts)
			}
		}

		// Draw border (rounded effect)
		borderImg := ebiten.NewImage(tileSize, tileSize)
		borderImg.Fill(color.RGBA{120, 120, 120, 255})
		borderOpts := &ebiten.DrawImageOptions{}
		borderOpts.GeoM.Translate(float64(x), float64(y0))
		borderOpts.CompositeMode = ebiten.CompositeModeSourceOver
		// Simulate rounded border by overlaying a smaller dark rect
		borderImg.DrawImage(ebiten.NewImage(tileSize-8, tileSize-8), &ebiten.DrawImageOptions{})

		// Draw block count with black outline for readability (bottom left)
		if count > 0 {
			countStr := fmt.Sprintf("%d", count)
			textX := x + 8
			textY := y0 + tileSize - 12
			DrawUITextOutline(screen, countStr, textX, textY, color.Black, color.White)
		}

	}
}

// getBlockTileImage returns the block's tile image for the hotbar
func getBlockTileImage(blockType block.BlockType) *ebiten.Image {
	if int(blockType) < len(tileImages) {
		return tileImages[blockType]
	}
	return nil
}

// DrawUIText draws simple text (for block count)
func DrawUIText(screen *ebiten.Image, text string, x, y int, clr color.Color) {
	// Use ebiten's built-in text drawing (or fallback to a simple method)
	// For now, use ebiten's debug text
	// TODO: Replace with proper font rendering
	if text != "" {
		// Use ebiten's debug print for simplicity
		// You can replace this with github.com/hajimehoshi/ebiten/text for real font
		ebitenutil.DebugPrintAt(screen, text, x, y)
	}
}
