# ---- Build Stage ----
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .

# Update go.mod and build the WASM binary
RUN go mod tidy && \
    GOOS=js GOARCH=wasm EBITEN_GRAPHICS_LIBRARY=opengl go build -o wasm/main.wasm github.com/KdntNinja/webcraft && \
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" wasm/

# ---- Production Stage ----
FROM python:3.12-alpine
WORKDIR /app

# Install curl for healthchecks
RUN apk add --no-cache curl

# Copy the built wasm directory from builder stage
COPY --from=builder /app/wasm /app/wasm

WORKDIR /app/wasm

# Install curl for healthchecks
RUN apk add --no-cache curl

# Copy the entire wasm directory to preserve structure
COPY --from=builder /app/wasm /app/wasm
COPY wasm /app/wasm
COPY wasm/index.html /app/wasm/index.html

EXPOSE 3000
CMD ["python", "-m", "http.server", "3000"]
