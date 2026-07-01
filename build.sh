#!/bin/bash
set -e

VERSION="$1"

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.2.0"
    exit 1
fi

# Build for Intel (amd64)
GOOS=darwin GOARCH=amd64 go build -o myapp_intel

# Build for ARM64
GOOS=darwin GOARCH=arm64 go build -o myapp_arm

# Create universal binary
lipo -create -output myapp myapp_intel myapp_arm

# Create directory structure
mkdir -p dist/Source-game-builder-tool-macos.app/Contents/MacOS/

# Remove old binary from app bundle
rm -f dist/Source-game-builder-tool-macos.app/Contents/MacOS/source-game-builder-tool

# Remove intermediate builds
rm -f myapp_arm myapp_intel

# Move new binary into app bundle
mv myapp dist/Source-game-builder-tool-macos.app/Contents/MacOS/source-game-builder-tool

# Make executable
chmod +x dist/Source-game-builder-tool-macos.app/Contents/MacOS/source-game-builder-tool

# Verify architectures
lipo -info dist/Source-game-builder-tool-macos.app/Contents/MacOS/source-game-builder-tool

# Create release archive
ZIP="Source-game-builder-tool-macos.zip"
rm -f "$ZIP"

ditto -c -k --sequesterRsrc --keepParent \
    dist/Source-game-builder-tool-macos.app \
    "$ZIP"

# -----------------------------
# CREATE CHECKSUM
# -----------------------------
SHA_FILE="${ZIP}.sha256"
shasum -a 256 "$ZIP" | awk '{print $1}' > "$SHA_FILE"

echo "Checksum created: $SHA_FILE"

# -----------------------------
# CREATE GITHUB RELEASE
# -----------------------------
gh release create "$VERSION" \
    "$ZIP" \
    "$SHA_FILE" \
    --title "$VERSION" \
    --generate-notes

echo "Release $VERSION created successfully!"
