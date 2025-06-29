package settings

// --- Player/Entity Physics ---
const (
	PlayerWidth        = TileSize     // Player width in pixels
	PlayerHeight       = TileSize * 2 // Player height in pixels
	PlayerMoveSpeed    = 4.3          // Player move speed
	PlayerJumpSpeed    = -12.0        // Player jump velocity (negative = up)
	PlayerGravity      = 0.7          // Player gravity per frame
	PlayerMaxFallSpeed = 15.0         // Player terminal velocity

	PlayerGroundFriction  = 0.6  // Ground friction multiplier
	PlayerAirResistance   = 0.98 // Air resistance multiplier
	PlayerGroundThreshold = 0.1  // Threshold for "on ground" state
)
