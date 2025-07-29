#!/bin/zsh

# Clear the terminal for a clean build log
clear

echo "Formatting Go code..."
go fmt

# Build the WASM binary and serve for development
set -e


# Ensure build directory exists
BUILD_DIR="web/build"
mkdir -p "$BUILD_DIR/css" "$BUILD_DIR/js"

MODULE_PATH="github.com/KdntNinja/webcraft"
# Force Ebiten to use WebGL1 for maximum browser compatibility (including Firefox)
echo "Building WebAssembly binary..."
env GOOS=js GOARCH=wasm EBITEN_GRAPHICS_LIBRARY=opengl go build -o "$BUILD_DIR/main.wasm" "$MODULE_PATH"

# Find and copy wasm_exec.js to the build directory
./scripts/find_wasm_exec.sh "$BUILD_DIR"

# Copy static assets to build directory
cp web/static/index.html "$BUILD_DIR/index.html"
cp web/css/style.css "$BUILD_DIR/css/style.css"
cp web/js/main.js "$BUILD_DIR/js/main.js"

# Run the native Go server

echo "Starting native server..."
go run web/serve.go