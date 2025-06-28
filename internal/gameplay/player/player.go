package player

import (
	"github.com/KdntNinja/webcraft/internal/core/engine/block"
	"github.com/KdntNinja/webcraft/internal/core/physics/entity"
	"github.com/KdntNinja/webcraft/internal/core/settings"
)

type Player struct {
	entity.AABB
	entity.InputState
	wasOnGround         bool            // Previous frame ground state
	SelectedBlock       block.BlockType // Currently selected block type for placing
	InteractionRange    float64         // Maximum range for block interaction
	LastInteractionTime int             // Frame counter for interaction cooldown
	InteractionCooldown int             // Cooldown frames between interactions (faster than inpututil)
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		AABB: entity.AABB{
			X: x, Y: y, Width: settings.PlayerWidth, Height: settings.PlayerHeight,
		},
		SelectedBlock:       block.Dirt,                     // Default to dirt blocks
		InteractionRange:    float64(settings.TileSize * 4), // 4 block radius
		InteractionCooldown: 0,                              // 3 frames cooldown (about 0.05 seconds at 60fps)
	}
}

func (p *Player) Update() {
	p.wasOnGround = p.OnGround // Store previous ground state

	// Increment interaction timer each frame
	p.LastInteractionTime++

	// Process input and update movement (without camera-dependent interactions)
	isMoving, targetVX, jumpKeyPressed, _ := p.HandleInput(0, 0) // Pass dummy camera values for basic input

	// Update input state tracking
	p.InputState.UpdateInputState(jumpKeyPressed, p.OnGround)

	p.ApplyMovement(isMoving, targetVX)
	p.HandleJump()
	p.ApplyGravity()
}

// HandleBlockInteractions processes block interactions with camera coordinates
func (p *Player) HandleBlockInteractions(cameraX, cameraY float64) *BlockInteraction {
	// Process input for block interactions
	_, _, _, blockInteraction := p.HandleInput(cameraX, cameraY)

	// Return the interaction directly since range checking is done in HandleInput
	return blockInteraction
}

// SetSelectedBlock changes the currently selected block type
func (p *Player) SetSelectedBlock(blockType block.BlockType) {
	p.SelectedBlock = blockType
}

// Entity interface implementations (delegate to AABB)
func (p *Player) CollideBlocks(blocks [][]int) {
	p.AABB.CollideBlocks(blocks)
}

func (p *Player) ClampX(min, max float64) {
	p.AABB.ClampX(min, max)
}

func (p *Player) GetPosition() (float64, float64) {
	return p.AABB.GetPosition()
}

func (p *Player) SetPosition(x, y float64) {
	p.AABB.SetPosition(x, y)
}
