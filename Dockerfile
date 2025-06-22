# Build stage
FROM golang:1.24.4-bookworm AS builder

WORKDIR /app/game
COPY ./game .

# Build Go WASM binary
RUN GOOS=js GOARCH=wasm go build -o main.wasm main.go

WORKDIR /app/server
COPY ./server .

# Final image: use Go to serve built game assets
FROM golang:1.24.4-bookworm
WORKDIR /app
# Copy only the built WASM and static assets (not Go source)
COPY --from=builder /app/game/index.html ./game/index.html
COPY --from=builder /app/game/main.wasm ./game/main.wasm
COPY --from=builder /app/game/wasm_exec.js ./game/wasm_exec.js
COPY --from=builder /app/server/serve.go ./server/serve.go

EXPOSE 8000
CMD ["go", "run", "server/serve.go"]