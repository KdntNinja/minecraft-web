# ---- Build Stage ----
FROM golang:1.22.4 AS builder
WORKDIR /app
COPY . .

# Build the WASM binary using the module path and copy wasm_exec.js
RUN GOOS=js GOARCH=wasm EBITEN_GRAPHICS_LIBRARY=opengl go build -o wasm/main.wasm github.com/KdntNinja/webcraft && \
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" wasm/

# ---- Production Stage ----
FROM python:3.12-alpine
WORKDIR /app

# Install curl for healthchecks
RUN apk add --no-cache curl

# Copy the built wasm directory from builder stage
COPY --from=builder /app/wasm /app/wasm

WORKDIR /app/wasm
EXPOSE 3000
CMD ["python", "-m", "http.server", "3000"]
