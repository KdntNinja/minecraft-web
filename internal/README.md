# Internal Engine Components

This directory contains the internal engine components for Webcraft, organized into logical modules:

## Directory Structure

- `block/` - Block definitions and chunk management
- `entity/` - Base entity system, physics, and collision detection
- `game/` - Main game loop and state management
- `noise/` - Procedural generation algorithms (terrain, caves, ores, biomes)
- `player/` - Player-specific logic, input handling, and movement
- `render/` - Rendering system and graphics utilities
- `world/` - World generation, terrain creation, and chunk management

## Architecture

The engine follows a component-based architecture where:

- Entities implement common interfaces for physics and rendering
- World generation uses sophisticated noise algorithms for varied terrain
- The render system optimizes for performance with batching and culling
- Physics system provides AABB collision detection with sub-pixel precision
