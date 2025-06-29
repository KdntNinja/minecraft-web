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

			// --- Use CanBreakBlock for break logic and UI ---
			if interactionType == BreakBlock && !p.CanBreakBlock(blockX, blockY) {
				return isMoving, targetVX, jumpKeyPressed, nil
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

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

// Add a helper to check if a block is breakable from the player's perspective
func (p *Player) CanBreakBlock(blockX, blockY int) bool {
	playerCenterX := int((p.AABB.X + float64(p.AABB.Width)/2) / float64(settings.TileSize))
	playerCenterY := int((p.AABB.Y + float64(p.AABB.Height)/2) / float64(settings.TileSize))

	x0, y0 := playerCenterX, playerCenterY
	x1, y1 := blockX, blockY
	dx := absInt(x1 - x0)
	dy := absInt(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	x, y := x0, y0
	for {
		if !(x == x0 && y == y0) && !(x == x1 && y == y1) {
			if p.World != nil {
				if p.World.GetBlockAt(x, y) != block.Air {
					return false
				}
				// If moving diagonally, check both adjacent cells to prevent corner breaking
				if x != x0 && y != y0 {
					if p.World.GetBlockAt(x, y0) != block.Air || p.World.GetBlockAt(x0, y) != block.Air {
						return false
					}
				}
			}
		}
		if x == x1 && y == y1 {
			break
		}
		err2 := 2 * err
		if err2 > -dy {
			err -= dy
			x += sx
		}
		if err2 < dx {
			err += dx
			y += sy
		}
	}
	return true
}
