package settings

const (
	// --- Player Physics ---
	PlayerSpriteWidth    = TileSize             // Visual width of the player sprite.
	PlayerSpriteHeight   = TileSize * 2         // Visual height of the player sprite.
	PlayerColliderWidth  = (TileSize * 9) / 10  // Physics bounding box width for collision.
	PlayerColliderHeight = (TileSize * 18) / 10 // Physics bounding box height for collision.
	PlayerMoveSpeed      = 4.3                  // Maximum horizontal walking speed (blocks/sec).
	PlayerJumpSpeed      = -9.0                 // Initial vertical jump velocity (upwards).
	PlayerGravity        = 0.45                 // Gravity force applied each frame.
	PlayerMaxFallSpeed   = 16.0                  // Maximum downward velocity (terminal velocity).

	// --- Movement Tuning ---
	PlayerWalkAccel      = 0.25                  // Acceleration when walking on the ground.
	PlayerAirAccel       = 0.04                  // Acceleration when in the air.
	PlayerGroundFriction = 0.55                  // Friction applied when on the ground and not moving.
	PlayerAirFriction    = 0.985                 // Air resistance applied when airborne.
	PlayerSneakSpeed     = PlayerMoveSpeed * 0.3 // Movement speed when sneaking.
	PlayerSneakAccel     = PlayerWalkAccel * 0.5 // Acceleration when sneaking.
	PlayerSprintSpeed    = PlayerMoveSpeed * 1.3 // Movement speed when sprinting.
	PlayerSprintAccel    = PlayerWalkAccel * 1.2 // Acceleration when sprinting.

	// --- Jump Mechanics ---
	PlayerCoyoteFrames     = 8    // Grace period (frames) to jump after leaving a ledge.
	PlayerJumpBufferFrames = 8    // Grace period (frames) to buffer a jump before landing.
	PlayerJumpHoldMax      = 12   // Max duration (frames) to hold jump for variable height.
	PlayerJumpHoldForce    = 0.32 // Upward force applied each frame when holding jump.
)
