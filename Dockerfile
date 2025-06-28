# ---- Build Stage ----
FROM golang:1.24.3-alpine AS builder
WORKDIR /app

# Copy source code
COPY . .

# Make script executable and run it to get wasm_exec.js
RUN chmod +x wasm/scripts/find_wasm_exec.sh && \
    wasm/scripts/find_wasm_exec.sh

# Update go.mod and build the WASM binary
RUN go mod tidy && \
    GOOS=js GOARCH=wasm EBITEN_GRAPHICS_LIBRARY=opengl go build -o wasm/main.wasm .

# ---- Production Stage ----
FROM golang:1.24.3-alpine
WORKDIR /app

# Install curl for healthchecks
RUN apk add --no-cache curl

# Copy the built wasm directory from builder stage
COPY --from=builder /app/wasm /app/wasm

WORKDIR /app/wasm

EXPOSE 3000
CMD ["go", "run", "serve.go"]
