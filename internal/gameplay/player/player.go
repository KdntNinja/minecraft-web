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
	World               WorldBlockGetter // Use concrete interface for better performance
	Health              int             // Player health
	MaxHealth           int             // Maximum health
	Inventory           [block.NumBlockTypes]int // Use array for fast inventory access
	IsSprinting         bool            // Sprinting state
}

// WorldBlockGetter is a minimal interface for world block access (concrete, not anonymous)
type WorldBlockGetter interface {
	GetBlockAt(x, y int) block.BlockType
}

func NewPlayer(x, y float64, world WorldBlockGetter) *Player {
	// Center collider horizontally in sprite, bottom-aligned
	colliderX := x + float64(settings.PlayerWidth-settings.PlayerColliderWidth)/2
	colliderY := y + float64(settings.PlayerHeight-settings.PlayerColliderHeight)
	p := &Player{
		AABB: entity.AABB{
			X:      colliderX,
			Y:      colliderY,
			Width:  settings.PlayerColliderWidth,
			Height: settings.PlayerColliderHeight,
		},
		SelectedBlock:       block.Grass,                    // Default to grass blocks (block 1)
		InteractionRange:    float64(settings.TileSize * 4), // 4 block radius
		InteractionCooldown: 0,                              // 3 frames cooldown (about 0.05 seconds at 60fps)
		World:               world,
		Health:              100,
		MaxHealth:           100,
		IsSprinting:         false,
	}
	// Inventory array is zeroed by default
	return p
}

// Update processes player state each frame
func (p *Player) Update() {
	// Only update wasOnGround if state changed
	if p.wasOnGround != p.OnGround {
		p.wasOnGround = p.OnGround
	}

	// Increment interaction timer each frame (no-op if already maxed)
	if p.LastInteractionTime < 1<<30 {
		p.LastInteractionTime++
	}

	// Process input and update movement (without camera-dependent interactions)
	isMoving, targetVX, jumpKeyPressed, _ := p.HandleInput(0, 0)

	// Sprinting mechanic: hold Shift to sprint
	p.IsSprinting = false // (Handled in input, but reset here for safety)

	// Only update input state if jump key state changed
	if jumpKeyPressed != p.InputState.JumpPressed {
		p.InputState.UpdateInputState(jumpKeyPressed, p.OnGround)
	}

	// Apply sprint speed boost
	if isMoving && p.IsSprinting {
		targetVX *= 1.5 // Sprint speed multiplier
	}

	p.ApplyMovement(isMoving, targetVX)
	p.HandleJump()
	p.ApplyGravity()
}
// TakeDamage reduces player health and clamps to zero
func (p *Player) TakeDamage(amount int) {
	p.Health -= amount
	if p.Health < 0 {
		p.Health = 0
	}
}

// Heal increases player health up to MaxHealth
func (p *Player) Heal(amount int) {
	p.Health += amount
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
}

// AddToInventory adds a block to the player's inventory (array version, fast)
func (p *Player) AddToInventory(blockType block.BlockType, count int) {
	if int(blockType) >= 0 && int(blockType) < len(p.Inventory) {
		p.Inventory[blockType] += count
	}
}

// RemoveFromInventory removes a block from the player's inventory (array version, fast)
func (p *Player) RemoveFromInventory(blockType block.BlockType, count int) bool {
	if int(blockType) >= 0 && int(blockType) < len(p.Inventory) && p.Inventory[blockType] >= count {
		p.Inventory[blockType] -= count
		return true
	}
	return false
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

// CanInteract returns true if the player can interact (based on cooldown)
func (p *Player) CanInteract() bool {
	return p.LastInteractionTime >= p.InteractionCooldown
}

// Entity interface implementations (delegate to AABB)
// CollideBlocksAdvanced: Use robust sub-stepping collision with PhysicsWorld
type PhysicsWorldProvider interface {
	GetPhysicsWorld() *entity.PhysicsWorld
}

func (p *Player) CollideBlocksAdvanced(world *entity.PhysicsWorld) {
	p.AABB.CollideBlocksAdvanced(world)
}

// ResetInteractionCooldown resets the cooldown timer after an interaction
func (p *Player) ResetInteractionCooldown() {
	p.LastInteractionTime = 0
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
