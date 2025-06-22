# Build stage
FROM golang:1.24.4-slim AS builder

WORKDIR /app/game
COPY ./game .

# Build Go WASM binary
RUN GOOS=js GOARCH=wasm go build -o main.wasm main.go

WORKDIR /app/server
COPY ./server .

# Final image: use Go to serve files
FROM golang:1.24.4-slim
WORKDIR /app
COPY --from=builder /app/game ./game
COPY --from=builder /app/server/serve.go ./server/serve.go

EXPOSE 8000
CMD ["go", "run", "server/serve.go"]