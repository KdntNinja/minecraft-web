# Webcraft

A 2D sandbox game engine inspired by Terraria and Minecraft, written in Go and powered by Ebiten.

## Features

- Dynamic chunk loading and procedural world generation
- WASM and native builds
- Modular, cycle-free architecture
- Modern static file server for web builds

## Project Structure

- `engine/` - Core game loop, rendering, and engine logic
- `gameplay/` - Game-specific logic (player, world, chunks, etc.)
- `coretypes/` - Shared interfaces and types for decoupling
- `assets/images/` - Game image assets
- `wasm/` - WASM build and static web files

## Build & Run

```sh
# Native
./run.sh

# WebAssembly
GOOS=js GOARCH=wasm go build -o wasm/main.wasm
```

## License

MIT
