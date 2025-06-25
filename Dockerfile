# ---- Build Stage ----
FROM golang:1.24.3 AS builder
WORKDIR /app
COPY . .
RUN GOOS=js GOARCH=wasm go build -o /app/wasm/main.wasm ./main.go

# ---- Production Stage ----
FROM python:3.12-alpine
WORKDIR /app

# Install curl for healthchecks
RUN apk add --no-cache curl

COPY --from=builder /app/wasm/main.wasm ./wasm/main.wasm
COPY wasm/wasm_exec.js ./wasm/wasm_exec.js
COPY wasm/index.html ./index.html
COPY wasm /app/wasm

EXPOSE 3000
CMD ["python", "-m", "http.server", "3000"]
