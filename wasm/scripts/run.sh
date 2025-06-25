#!/bin/zsh
# Build the WASM binary and serve for development
set -e

echo "Building main.wasm..."
GOOS=js GOARCH=wasm go build -o main.wasm ../main.go

echo "Copying wasm_exec.js if needed..."
if [ ! -f wasm_exec.js ]; then
  cp $(go env GOROOT)/misc/wasm/wasm_exec.js wasm_exec.js
fi

echo "Starting local server at http://localhost:8080/wasm/index.html"
./serve.sh
