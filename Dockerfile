# Use the official Rust image as the build environment
FROM rust:latest as builder

WORKDIR /app

# Copy source code and manifest
COPY . .

# Install wasm32 target for Rust
RUN rustup target add wasm32-unknown-unknown

# Build the project with the custom profile for wasm
RUN cargo build --release --target wasm32-unknown-unknown --profile wasm-release

# Final image (optional: use a minimal image if you want to serve the WASM)
FROM debian:bullseye-slim as final
WORKDIR /app

# Copy the built WASM file(s) from the builder
COPY --from=builder /app/target/wasm32-unknown-unknown/wasm-release /app/wasm-release

# Set default command (list files for demonstration)
CMD ["ls", "/app/wasm-release"]
