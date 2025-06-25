#!/bin/sh

# Get GOROOT
GOROOT=$(go env GOROOT)

# Set source and destination paths
# Search for wasm_exec.js in GOROOT and list all found files
FOUND_FILES=$(find "$GOROOT" -name wasm_exec.js 2>/dev/null)
echo "Found the following wasm_exec.js files:"
echo "$FOUND_FILES"

# Use the first found file as the source
SRC=$(echo "$FOUND_FILES" | head -n 1)
DEST="$(dirname "$0")/wasm_exec.js"

# Check if source exists
if [ ! -f "$SRC" ]; then
    echo "wasm_exec.js not found in GOROOT."
    exit 1
fi

# Create destination directory if it doesn't exist
DEST_DIR="$(dirname "$DEST")"
if [ ! -d "$DEST_DIR" ]; then
    mkdir -p "$DEST_DIR"
fi

# Copy the file
cp -f "$SRC" "$DEST"
echo "Copied $SRC to $DEST"
