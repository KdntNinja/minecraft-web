#!/bin/zsh

clear

go fmt

# Build the WASM binary and serve for development
set -e

# Build the WASM binary using the module path for reproducibility
MODULE_PATH="github.com/KdntNinja/webcraft"
# Force Ebiten to use WebGL1 for maximum browser compatibility (including Firefox)
env GOOS=js GOARCH=wasm EBITEN_GRAPHICS_LIBRARY=opengl go build -o wasm/main.wasm "$MODULE_PATH"

# Prefer Go 1.24+ location, fallback to older if not found
GOROOT=$(go env GOROOT)
if [ -f "$GOROOT/lib/wasm/wasm_exec.js" ]; then
  cp "$GOROOT/lib/wasm/wasm_exec.js" wasm/
elif [ -f "$GOROOT/misc/wasm/wasm_exec.js" ]; then
  cp "$GOROOT/misc/wasm/wasm_exec.js" wasm/
else
  ./wasm/scripts/find_wasm_exec.sh
fi

cd wasm
go run serve.go