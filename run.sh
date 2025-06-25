#!/bin/zsh

clear

go fmt

# Build the WASM binary and serve for development
set -e

# Build the WASM binary using the module path for reproducibility
MODULE_PATH="github.com/KdntNinja/webcraft"
echo "Building wasm/main.wasm..."
env GOOS=js GOARCH=wasm go build -o wasm/main.wasm "$MODULE_PATH"

echo "Copying wasm_exec.js if needed..."
# Prefer Go 1.24+ location, fallback to older if not found
GOROOT=$(go env GOROOT)
if [ -f "$GOROOT/lib/wasm/wasm_exec.js" ]; then
  cp "$GOROOT/lib/wasm/wasm_exec.js" wasm/
elif [ -f "$GOROOT/misc/wasm/wasm_exec.js" ]; then
  cp "$GOROOT/misc/wasm/wasm_exec.js" wasm/
else
  ./wasm/scripts/find_wasm_exec.sh
fi

echo "Starting local server at http://localhost:8080/"
cd wasm
python3 -m http.server 8080