#!/bin/zsh
# Build the WASM binary and serve for development
set -e

echo "Building wasm/main.wasm..."
GOOS=js GOARCH=wasm go build -o wasm/main.wasm main.go

echo "Copying wasm_exec.js if needed..."
if [ ! -f wasm/wasm_exec.js ]; then
  ./wasm/scripts/find_wasm_exec.sh
fi

echo "Starting local server at http://localhost:8080/"
cd wasm
./scripts/serve.sh