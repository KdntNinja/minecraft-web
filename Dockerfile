
# ---- Build Stage ----
FROM golang:1.24.3-alpine AS builder
WORKDIR /app

# Copy source code (flat structure, no internal/)
COPY go.mod go.sum ./
COPY *.go ./
COPY wasm/ ./wasm/
COPY coretypes/ ./coretypes/
COPY engine/ ./engine/
COPY gameplay/ ./gameplay/
COPY generation/ ./generation/
COPY rendering/ ./rendering/
COPY physics/ ./physics/
COPY progress/ ./progress/
COPY settings/ ./settings/

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
COPY --from=builder /app/*.go /app/
COPY --from=builder /app/coretypes /app/coretypes
COPY --from=builder /app/engine /app/engine
COPY --from=builder /app/gameplay /app/gameplay
COPY --from=builder /app/generation /app/generation
COPY --from=builder /app/rendering /app/rendering
COPY --from=builder /app/physics /app/physics
COPY --from=builder /app/progress /app/progress
COPY --from=builder /app/settings /app/settings

EXPOSE 3000
CMD ["go", "run", "main.go"]
