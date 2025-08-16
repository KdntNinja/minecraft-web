package gameplay

import (
	"github.com/KdntNinja/webcraft/coretypes"
	"github.com/KdntNinja/webcraft/physics"
	"github.com/KdntNinja/webcraft/settings"
)

type Player struct {
	physics.AABB
	physics.InputState
	wasOnGround           bool                         // Previous frame ground state
	SelectedBlock         coretypes.BlockType          // Currently selected block type for placing
	InteractionRange      float64                      // Maximum range for block interaction
	LastInteractionTime   int                          // Frame counter for interaction cooldown
	InteractionCooldown   int                          // Cooldown frames between interactions (faster than inpututil)
	World                 WorldBlockGetter             // Use concrete interface for better performance
	Health                int                          // Player health
	MaxHealth             int                          // Maximum health
	Inventory             [coretypes.NumBlockTypes]int // Use array for fast inventory access
	Hotbar                []coretypes.BlockType        // Dynamic hotbar (up to 9 blocks)
	IsSprinting           bool                         // Sprinting state
	lastEmptiedHotbarSlot int                          // -1 if none
	// ...existing code...
}

// WorldBlockGetter is a minimal interface for world block access (concrete, not anonymous)
type WorldBlockGetter interface {
	GetBlockAt(x, y int) coretypes.BlockType
}

func NewPlayer(x, y float64, world WorldBlockGetter) *Player {
	// Center collider horizontally in sprite, bottom-aligned
	colliderX := x + float64(settings.PlayerSpriteWidth-settings.PlayerColliderWidth)/2
	colliderY := y + float64(settings.PlayerSpriteHeight-settings.PlayerColliderHeight)
	p := &Player{
		AABB: physics.AABB{
			X:      colliderX,
			Y:      colliderY,
			Width:  settings.PlayerColliderWidth,
			Height: settings.PlayerColliderHeight,
		},
		SelectedBlock:         coretypes.Grass,                // Default to grass blocks (block 1)
		InteractionRange:      float64(settings.TileSize * 4), // 4 block radius
		InteractionCooldown:   0,                              // 3 frames cooldown (about 0.05 seconds at 60fps)
		World:                 world,
		Health:                100,
		MaxHealth:             100,
		Hotbar:                make([]coretypes.BlockType, 9), // Always 9 slots, filled with coretypes.Air
		lastEmptiedHotbarSlot: -1,
		IsSprinting:           false,
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

	// Only update input state if jump key state changed
	if jumpKeyPressed != p.InputState.JumpPressed {
		p.InputState.UpdateInputState(jumpKeyPressed, p.OnGround)
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
func (p *Player) AddToInventory(blockType coretypes.BlockType, count int) {
	if int(blockType) >= 0 && int(blockType) < len(p.Inventory) {
		p.Inventory[blockType] += count
		// Add to hotbar if not already present and not air
		if blockType != coretypes.Air {
			found := false
			for _, b := range p.Hotbar {
				if b == blockType {
					found = true
					break
				}
			}
			if !found {
				slot := -1
				// Prefer to fill the most recently emptied slot
				if p.lastEmptiedHotbarSlot >= 0 && p.lastEmptiedHotbarSlot < len(p.Hotbar) && p.Hotbar[p.lastEmptiedHotbarSlot] == coretypes.Air {
					slot = p.lastEmptiedHotbarSlot
				} else {
					// Otherwise, find first empty slot
					for i, b := range p.Hotbar {
						if b == coretypes.Air {
							slot = i
							break
						}
					}
				}
				// If no empty slot, replace first slot
				if slot == -1 && len(p.Hotbar) > 0 {
					slot = 0
				}
				if slot != -1 {
					p.Hotbar[slot] = blockType
					p.lastEmptiedHotbarSlot = -1 // Reset after use
				}
			}
		}
	}
}

// RemoveFromInventory removes a block from the player's inventory (array version, fast)
func (p *Player) RemoveFromInventory(blockType coretypes.BlockType, count int) bool {
	if int(blockType) >= 0 && int(blockType) < len(p.Inventory) && p.Inventory[blockType] >= count {
		p.Inventory[blockType] -= count
		// If the count is now zero, set hotbar slot to coretypes.Air (do not shift others)
		if p.Inventory[blockType] == 0 {
			for i, b := range p.Hotbar {
				if b == blockType {
					p.Hotbar[i] = coretypes.Air
					p.lastEmptiedHotbarSlot = i
					break
				}
			}
		}
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
func (p *Player) SetSelectedBlock(blockType coretypes.BlockType) {
	p.SelectedBlock = blockType
}

// CanInteract returns true if the player can interact (based on cooldown)
func (p *Player) CanInteract() bool {
	return p.LastInteractionTime >= p.InteractionCooldown
}

// CollideBlocks resolves collisions using the provided physics world
func (p *Player) CollideBlocks(world *physics.PhysicsWorld) {
	p.AABB.CollideBlocks(world)
}

// SetGridOffset sets the collision grid offset for the player
func (p *Player) SetGridOffset(x, y int) {
	p.AABB.SetGridOffset(x, y)
}

// GetX returns the player's X position
func (p *Player) GetX() float64 {
	return p.AABB.X
}

// GetY returns the player's Y position
func (p *Player) GetY() float64 {
	return p.AABB.Y
}

// GetColliderWidth returns the player's collider width
func (p *Player) GetColliderWidth() float64 {
	return float64(p.AABB.Width)
}

// GetColliderHeight returns the player's collider height
func (p *Player) GetColliderHeight() float64 {
	return float64(p.AABB.Height)
}

// GetSelectedBlock returns the currently selected block index
func (p *Player) GetSelectedBlock() int {
	return int(p.SelectedBlock)
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
