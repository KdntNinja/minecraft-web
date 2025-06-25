# ---- Build Stage ----
FROM golang:1.24.3 AS builder
WORKDIR /app
COPY . .
RUN GOOS=js GOARCH=wasm go build -o /app/wasm/main.wasm ./main.go

# ---- Production Stage ----
FROM nginx:alpine
WORKDIR /usr/share/nginx/html
COPY --from=builder /app/wasm/main.wasm ./wasm/main.wasm
COPY wasm/wasm_exec.js ./wasm/wasm_exec.js
COPY wasm/index.html ./wasm/index.html
COPY wasm /usr/share/nginx/html/wasm

# Change Nginx to listen on port 3000
RUN sed -i 's/listen\s\+80;/listen 3000;/' /etc/nginx/conf.d/default.conf

EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
