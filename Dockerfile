#########################################################
# Build Stage: compile WASM and prepare static build
#########################################################
FROM golang:1.24.3 AS builder
WORKDIR /app

# Fetch modules
COPY go.mod go.sum ./
RUN go mod download

# Copy all source, assets, and scripts
COPY . ./

# Format and build WASM + static assets
RUN go fmt ./...
RUN mkdir -p web/build/css web/build/js && \
    GOOS=js GOARCH=wasm EBITEN_GRAPHICS_LIBRARY=opengl go build -o web/build/main.wasm . && \
    chmod +x scripts/find_wasm_exec.sh && scripts/find_wasm_exec.sh web/build && \
    cp web/static/index.html web/build/index.html && \
    cp web/css/style.css web/build/css/style.css && \
    cp web/js/main.js web/build/js/main.js

#########################################################
# Production Stage: run Go static server for WASM build
#########################################################
FROM golang:1.24.3-alpine AS production
WORKDIR /app

# Install minimal tools
RUN apk add --no-cache curl

# Copy static build and server code
COPY --from=builder /app/web/build web/build
COPY web/serve.go ./serve.go

EXPOSE 3000
# Launch the Go static file server
CMD ["go", "run", "serve.go"]