package player

import "github.com/KdntNinja/webcraft/internal/block"

// Player physics and movement constants
const (
	Width        = block.TileSize     // Player hitbox width (42px)
	Height       = block.TileSize * 2 // Player hitbox height (84px, 2 blocks tall)
	MoveSpeed    = 4.3                // Horizontal movement speed in pixels/frame
	JumpSpeed    = -12.0              // Initial jump velocity (negative = upward)
	Gravity      = 0.7                // Gravity acceleration per frame
	MaxFallSpeed = 15.0               // Terminal velocity cap

	// Movement feel adjustments
	GroundFriction = 0.6  // Friction multiplier when on ground (0.6 = 40% speed reduction)
	AirResistance  = 0.98 // Air resistance multiplier (0.98 = 2% speed reduction)

	// Collision precision
	GroundThreshold = 0.1 // Maximum distance from ground to allow jumping
)
