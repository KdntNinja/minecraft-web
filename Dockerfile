# Build stage
FROM golang:1.24.4-bookworm AS builder

WORKDIR /app/game
COPY ./game .

# Build Go WASM binary
RUN GOOS=js GOARCH=wasm go build -o wasm/main.wasm main.go

# Final image: use Python to serve the game directory
FROM python:alpine
WORKDIR /app/game
COPY --from=builder /app/game/static ./static
COPY --from=builder /app/game/wasm ./wasm
EXPOSE 8000
CMD ["python3", "-m", "http.server", "8000"]