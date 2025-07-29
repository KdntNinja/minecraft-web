# Settings

Configuration and settings management for Webcraft.

- Stores all tunable game constants (tile size, chunk size, player speed, etc.)
- Used by engine, gameplay, rendering, and generation
- Supports both build-time and runtime configuration

## Structure

- `settings.go` - Main settings and constants

## Usage

- Import and use for all game-wide constants
- Change values here to tune gameplay or performance

## Notes

- Avoid hardcoding magic numbers elsewhere in the codebase
- Document all settings for clarity
