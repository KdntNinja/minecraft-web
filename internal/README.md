# Webcraft Internal Structure

This directory contains the core engine and gameplay components for Webcraft, organized into modular packages for clarity and maintainability.

## Directory Overview

- `core/`
  - `engine/`
    - `block/`   — Block definitions, chunk/tile sizing (now uses centralized settings)
    - `game/`    — Main game loop, state management, camera
    - `graphics/`— Graphics utilities (WIP)
    - `physics/`
      - `entity/` — Entity base, movement, collision, physics resolution
  - `settings/`  — Centralized game constants (player, world, rendering, noise, etc.)
- `gameplay/`
  - `player/`    — Player logic, input, movement (uses settings for all constants)
  - `world/`     — World and chunk management
- `generation/`
  - `noise/`     — Procedural generation: terrain, caves, ores, biomes
  - `terrain/`   — Terrain chunk generation and caching
- `systems/`
  - `rendering/`
    - `render/`  — Rendering system, camera, color utilities
- `wasm/`        — WebAssembly entrypoint, HTML, and scripts

## Architecture Notes

- **Component-based design:** Entities implement interfaces for physics, rendering, and interaction.
- **Centralized settings:** All major constants (player, world, rendering, etc.) are defined in `core/settings/settings.go` and used throughout the codebase.
- **Procedural world generation:** Uses advanced noise algorithms for varied, infinite terrain.
- **Optimized rendering:** Batching, culling, and camera logic for smooth performance.
- **Robust physics:** AABB collision detection, sub-pixel precision, and entity movement.

---

For more details, see the documentation in each subdirectory and the main project README.
