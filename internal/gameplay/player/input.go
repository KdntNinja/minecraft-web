package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

// BlockInteraction represents a block interaction event
type BlockInteraction struct {
	Type    BlockInteractionType
	BlockX  int // Block coordinate
	BlockY  int // Block coordinate
	ScreenX int // Screen pixel coordinate
	ScreenY int // Screen pixel coordinate
}

type BlockInteractionType int

const (
	BreakBlock BlockInteractionType = iota
	PlaceBlock
)

// HandleInput processes keyboard and mouse input and returns movement intentions and block interactions
func (p *Player) HandleInput(cameraX, cameraY float64) (isMoving bool, targetVX float64, jumpKeyPressed bool, blockInteraction *BlockInteraction) {
	// Sneak (hold Shift)
	p.InputState.SneakPressed = ebiten.IsKeyPressed(ebiten.KeyShift) || ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	isMoving = false
	targetVX = 0.0

	// Check horizontal movement keys (WASD + arrows)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		targetVX = -settings.PlayerMoveSpeed
		isMoving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		targetVX = settings.PlayerMoveSpeed
		isMoving = true
	}

	// Check jump keys (multiple options for accessibility)
	jumpKeyPressed = ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeySpace)

	// Handle block selection with number keys
	p.handleBlockSelection()

	// Handle mouse input for block interaction (instant response)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		mouseX, mouseY := ebiten.CursorPosition()
		// Use the original mapping: mouseX/mouseY are relative to the top-left of the screen.
		worldX := float64(mouseX) + cameraX
		worldY := float64(mouseY) + cameraY
		blockX := int(worldX / float64(settings.TileSize))
		blockY := int(worldY / float64(settings.TileSize))
		if worldX < 0 {
			blockX = int(worldX/float64(settings.TileSize)) - 1
		}
		if worldY < 0 {
			blockY = int(worldY/float64(settings.TileSize)) - 1
		}

		// Check if block is within interaction range (from player center to block center, matching crosshair)
		playerCenterX := p.AABB.X + float64(p.AABB.Width)/2
		playerCenterY := p.AABB.Y + float64(p.AABB.Height)/2
		blockCenterX := float64(blockX)*float64(settings.TileSize) + float64(settings.TileSize)/2
		blockCenterY := float64(blockY)*float64(settings.TileSize) + float64(settings.TileSize)/2
		dx := blockCenterX - playerCenterX
		dy := blockCenterY - playerCenterY
		distance := dx*dx + dy*dy

		if distance <= p.InteractionRange*p.InteractionRange {
			// Determine interaction type
			var interactionType BlockInteractionType
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				interactionType = BreakBlock
			} else {
				interactionType = PlaceBlock
			}

			// Reset cooldown timer
			p.LastInteractionTime = 0

			blockInteraction = &BlockInteraction{
				Type:    interactionType,
				BlockX:  blockX,
				BlockY:  blockY,
				ScreenX: mouseX,
				ScreenY: mouseY,
			}
		}
	}

	return isMoving, targetVX, jumpKeyPressed, blockInteraction
}

// handleBlockSelection processes number key input to change selected block type
func (p *Player) handleBlockSelection() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		p.SelectedBlock = block.Grass
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		p.SelectedBlock = block.Dirt
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		p.SelectedBlock = block.Clay
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		p.SelectedBlock = block.Stone
	} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
		p.SelectedBlock = block.CopperOre
	} else if inpututil.IsKeyJustPressed(ebiten.Key6) {
		p.SelectedBlock = block.IronOre
	} else if inpututil.IsKeyJustPressed(ebiten.Key7) {
		p.SelectedBlock = block.GoldOre
	} else if inpututil.IsKeyJustPressed(ebiten.Key8) {
		p.SelectedBlock = block.Ash
	} else if inpututil.IsKeyJustPressed(ebiten.Key9) {
		p.SelectedBlock = block.Wood
	} else if inpututil.IsKeyJustPressed(ebiten.Key0) {
		p.SelectedBlock = block.Leaves
	}
}
