package settings

// --- Cave Generation Parameters ---
const (
	CaveSurfaceEntranceMinDepth = -4     // Min depth (relative to surface) for cave entrances
	CaveSurfaceEntranceMaxDepth = 12     // Max depth (relative to surface) for cave entrances
	CaveSurfaceEntranceScale    = 30.0   // Surface cave entrance noise scale
	CaveSurfaceEntranceOffset   = 8000.0 // Surface cave entrance noise offset
	CaveSurfaceEntranceThresh   = 0.55   // Surface cave entrance threshold (higher = fewer)

	CaveLargeScale        = 50.0   // Large cavern noise scale (broad spaces)
	CaveHorizontalScale   = 18.0   // Horizontal tunnel noise scale
	CaveHorizontalYOffset = 1000.0 // Horizontal tunnel noise Y offset
	CaveHorizontalYScale  = 35.0   // Horizontal tunnel Y blending scale
	CaveVerticalScale     = 35.0   // Vertical tunnel noise scale
	CaveVerticalYOffset   = 2000.0 // Vertical tunnel noise Y offset
	CaveVerticalYScale    = 15.0   // Vertical tunnel Y blending scale
	CaveSmallScale        = 10.0   // Small cave noise scale (pockets)
	CaveSmallYOffset      = 3000.0 // Small cave noise Y offset
	CaveAirPocketScale    = 6.0    // Air pocket noise scale (tiny holes)
	CaveAirPocketYOffset  = 4000.0 // Air pocket noise Y offset

	CaveVeryDeepDepth   = 150 // Very deep caves start below this depth
	CaveDeepDepth       = 100 // Deep caves start below this depth
	CaveMediumDepth     = 50  // Medium caves start below this depth
	CaveShallowDepth    = 15  // Shallow caves start below this depth
	CaveMinShallowDepth = 2   // Minimum depth for any cave

	CaveVeryDeepLargeWeight  = 0.4  // Very deep: large cavern weight
	CaveVeryDeepHorizWeight  = 0.3  // Very deep: horizontal tunnel weight
	CaveVeryDeepVertWeight   = 0.2  // Very deep: vertical tunnel weight
	CaveVeryDeepSmallWeight  = 0.1  // Very deep: small cave weight
	CaveVeryDeepThresh       = 0.15 // Very deep: cave generation threshold
	CaveVeryDeepTunnelWeight = 0.6  // Very deep: interconnected tunnel weight
	CaveVeryDeepPocketWeight = 0.4  // Very deep: air pocket weight
	CaveVeryDeepTunnelThresh = 0.35 // Very deep: tunnel generation threshold

	CaveDeepLargeWeight  = 0.3  // Deep: large cavern weight
	CaveDeepHorizWeight  = 0.3  // Deep: horizontal tunnel weight
	CaveDeepSmallWeight  = 0.3  // Deep: small cave weight
	CaveDeepPocketWeight = 0.1  // Deep: air pocket weight
	CaveDeepThresh       = 0.18 // Deep: cave generation threshold
	CaveDeepVertThresh   = 0.45 // Deep: vertical tunnel threshold

	CaveMediumHorizWeight  = 0.4  // Medium: horizontal tunnel weight
	CaveMediumSmallWeight  = 0.4  // Medium: small cave weight
	CaveMediumPocketWeight = 0.2  // Medium: air pocket weight
	CaveMediumThresh       = 0.22 // Medium: cave generation threshold
	CaveMediumVertThresh   = 0.5  // Medium: vertical tunnel threshold

	CaveShallowHorizWeight  = 0.3  // Shallow: horizontal tunnel weight
	CaveShallowSmallWeight  = 0.5  // Shallow: small cave weight
	CaveShallowPocketWeight = 0.2  // Shallow: air pocket weight
	CaveShallowThresh       = 0.18 // Shallow: cave generation threshold
	CaveShallowVertThresh   = 0.45 // Shallow: vertical tunnel threshold

	CaveMinShallowSmallWeight  = 0.5  // Min-depth: small cave weight
	CaveMinShallowPocketWeight = 0.3  // Min-depth: air pocket weight
	CaveMinShallowHorizWeight  = 0.2  // Min-depth: horizontal tunnel weight
	CaveMinShallowThresh       = 0.22 // Min-depth: cave generation threshold
)
