#!/bin/bash

# Build for Intel (amd64)
GOOS=darwin GOARCH=amd64 go build -o myapp_intel

# Build for ARM64
GOOS=darwin GOARCH=arm64 go build -o myapp_arm

# Create universal binary
lipo -create -output myapp myapp_intel myapp_arm

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