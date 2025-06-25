# ---- Build Stage ----
FROM golang:1.22.4 AS builder
WORKDIR /app
COPY . .
RUN GOOS=js GOARCH=wasm go build -o /app/wasm/main.wasm ./main.go

# ---- Production Stage ----
FROM python:3.12-alpine
WORKDIR /app/wasm

# Install curl for healthchecks
RUN apk add --no-cache curl

# Copy the entire wasm directory to preserve structure
COPY --from=builder /app/wasm /app/wasm
COPY wasm /app/wasm
COPY wasm/index.html /app/wasm/index.html

EXPOSE 3000
CMD ["python", "-m", "http.server", "3000"]
