#!/bin/sh

# Get GOROOT
GOROOT=$(go env GOROOT)
echo "GOROOT: $GOROOT"

# Set destination path relative to wasm directory (not scripts directory)
DEST="$(dirname "$(dirname "$0")")/wasm_exec.js"
echo "Destination: $DEST"

# Search for wasm_exec.js in GOROOT and list all found files
echo "Searching for wasm_exec.js in GOROOT..."
FOUND_FILES=$(find "$GOROOT" -name "wasm_exec.js" 2>/dev/null)

if [ -z "$FOUND_FILES" ]; then
    echo "No wasm_exec.js files found in GOROOT."
    echo "Trying alternative locations..."
    
    # Try common Go installation paths
    for path in /usr/local/go /opt/go /usr/lib/go; do
        if [ -d "$path" ]; then
            echo "Checking $path..."
            ALT_FILES=$(find "$path" -name "wasm_exec.js" 2>/dev/null)
            if [ -n "$ALT_FILES" ]; then
                FOUND_FILES="$ALT_FILES"
                break
            fi
        fi
    done
fi

if [ -z "$FOUND_FILES" ]; then
    echo "ERROR: wasm_exec.js not found in any location."
    echo "Please ensure Go is properly installed with WebAssembly support."
    exit 1
fi

echo "Found the following wasm_exec.js files:"
echo "$FOUND_FILES"

# Use the first found file as the source
SRC=$(echo "$FOUND_FILES" | head -n 1)
echo "Using: $SRC"

# Check if source exists and is readable
if [ ! -f "$SRC" ] || [ ! -r "$SRC" ]; then
    echo "ERROR: Source file $SRC is not accessible."
    exit 1
fi

# Create destination directory if it doesn't exist
DEST_DIR="$(dirname "$DEST")"
if [ ! -d "$DEST_DIR" ]; then
    echo "Creating destination directory: $DEST_DIR"
    mkdir -p "$DEST_DIR"
fi

# Copy the file
echo "Copying $SRC to $DEST"
if cp "$SRC" "$DEST"; then
    echo "SUCCESS: wasm_exec.js copied successfully."
    # Verify the copy
    if [ -f "$DEST" ]; then
        echo "Verified: Destination file exists and is $(wc -c < "$DEST") bytes."
    fi
else
    echo "ERROR: Failed to copy file."
    exit 1
fi