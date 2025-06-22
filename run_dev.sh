# Build the WASM binary
GOOS=js GOARCH=wasm go build -o game/wasm/main.wasm game/main.go

# Serve static and wasm files from project root
python3 -m http.server 8000