package physics

// AABB (Axis-Aligned Bounding Box) for collision detection and physics
type AABB struct {
	X, Y          float64 // Position in world coordinates
	Width, Height int     // Size in pixels
	VX, VY        float64 // Velocity in pixels per frame
	OnGround      bool    // Whether entity is touching ground
	GridOffsetX   int     // Offset for collision grid X (for infinite world)
	GridOffsetY   int     // Offset for collision grid Y (for infinite world)
}

// Entity interface implementations for AABB

func (a *AABB) ClampX(min, max float64) {
	if a.X < min {
		a.X = min
	}
	if a.X > max {
		a.X = max
	}
}

func (a *AABB) GetPosition() (float64, float64) {
	return a.X, a.Y
}

func (a *AABB) SetPosition(x, y float64) {
	a.X = x
	a.Y = y
}

func (a *AABB) Update() {
	// Default empty implementation - override in specific entities
}

// Entities type moved to coretypes/entity.go

// InputState tracks input state for entities
type InputState struct {
	JumpPressed      bool // Current jump key state
	WasJumpPressed   bool // Previous jump key state
	LastGroundedTime int  // Frames since last grounded
	LastJumpPressed  int  // Frames since jump was pressed (for jump buffer)
	SneakPressed     bool // Whether sneak (shift) is held
	JumpHoldTime     int  // Frames jump has been held for variable jump height
}

// UpdateInputState updates input state tracking
func (i *InputState) UpdateInputState(jumpPressed bool, onGround bool) {
	i.WasJumpPressed = i.JumpPressed
	i.JumpPressed = jumpPressed

	if onGround {
		i.LastGroundedTime = 0
	} else {
		i.LastGroundedTime++
	}
}

// CanJump checks if a fresh jump input was received
func (i *InputState) CanJump() bool {
	return i.JumpPressed && !i.WasJumpPressed
}
