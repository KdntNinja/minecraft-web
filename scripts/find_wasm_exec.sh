#!/bin/sh
# Copy wasm_exec.js from /usr/local/go/misc/wasm/wasm_exec.js to web/build/wasm_exec.js (or custom dir)
DEST_DIR="web/build"
if [ -n "$1" ]; then
    DEST_DIR="$1"
fi
DEST="$DEST_DIR/wasm_exec.js"
SRC="/usr/local/go/misc/wasm/wasm_exec.js"
echo "Copying $SRC to $DEST"
mkdir -p "$DEST_DIR"
if cp "$SRC" "$DEST"; then
    echo "SUCCESS: wasm_exec.js copied successfully."
    echo "Verified: Destination file exists and is $(wc -c < "$DEST") bytes."
else
    echo "ERROR: Failed to copy file."
    exit 1
fi
