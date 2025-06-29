package settings

// --- Terrain Height Generation Constants ---
const (
	TerrainBaseScale       = 120.0  // Base scale for low-frequency terrain
	TerrainHillScale       = 40.0   // Scale for hills (medium-frequency noise)
	TerrainHillOffset      = 500.0  // Offset for hill noise
	TerrainDetailScale     = 10.0   // Scale for small terrain details (high-frequency noise)
	TerrainDetailOffset    = 2000.0 // Offset for detail noise
	TerrainBaseWeight      = 0.6    // Weight for base terrain layer
	TerrainHillWeight      = 0.3    // Weight for hill layer
	TerrainDetailWeight    = 0.1    // Weight for detail layer
	TerrainMinHeight       = 8      // Minimum allowed terrain height (blocks)
	TerrainMaxHeightBuffer = 8      // Buffer from world bottom for max height
)

// --- Tree/Biome Noise Constants (used in surface.go) ---
const (
	TreeBiomeNoiseScale = 80.0 // Scale for biome noise (surface variation)
	TreeNoiseScale      = 12.0 // Scale for tree placement noise
	TreeClayBiomeThresh = 0.45 // Threshold for clay biome
	TreeClayNoiseThresh = 0.5  // Threshold for clay noise
)

// --- Terraria-like Sky Transition Constants ---
const (
	SkyTopY             = -100.0 // Y above this is always sky blue
	SkyTransitionStartY = -50.0  // Start fading here
	SkyTransitionEndY   = 150.0  // Fully dark here
)
