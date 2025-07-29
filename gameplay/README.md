# Gameplay

Game-specific logic for Webcraft, including:

- **Player**: Handles player entity, input, movement, and physics integration (`player/`)
- **World**: Manages world state, dynamic chunk loading, block and entity management, and collision grid (`world/`)
- **Chunks**: Chunk coordinate math, chunk manager, and chunk loading logic (`world/chunks/`)
- **Entities**: Entity system and update logic for all in-game entities
- **Procedural Generation**: Terrain, caves, ores, and trees generation (`generation/`)
- **Settings**: Game settings and configuration (`settings/`)
- **Progress**: Progress tracking and reporting for world generation and loading (`progress/`)

## Structure

- `player/` - Player entity, input, and movement
- `world/` - World state, chunk manager, block and entity logic
- `world/chunks/` - Chunk coordinate math and chunk management
- `generation/` - Procedural world generation (terrain, caves, ores, trees)
- `settings/` - Game settings and configuration
- `progress/` - Progress tracking and reporting

## Dependencies

- Depends only on `coretypes` (shared interfaces/types) and engine interfaces
- No direct dependencies on rendering or physics packages

## Notes

- All shared types/interfaces are in `coretypes` to avoid import cycles
- Designed for modularity and testability
