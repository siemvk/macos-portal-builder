#!/bin/bash

# macOS Source Builder Oneliner Installer
# This script downloads the latest release of the macOS Source Builder and runs it.

set -euo pipefail

# Repository information
REPO="SSoggyTacoMan/macos-portal-builder"
INSTALL_DIR="$HOME/.macos-source-builder"
BINARY_NAME="source-game-builder-tool"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

printf "${BLUE}==>${NC} macOS Source Builder Installer\n"

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    printf "${RED}Error:${NC} This tool is only for macOS.\n"
    exit 1
fi

# Create install directory
mkdir -p "$INSTALL_DIR"

printf "${BLUE}==>${NC} Fetching latest release information...\n"
RELEASE_JSON=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest")
TAG=$(echo "$RELEASE_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
    printf "${RED}Error:${NC} Could not find latest release. Please check $REPO\n"
    exit 1
fi

printf "${BLUE}==>${NC} Found version: ${GREEN}$TAG${NC}\n"

# In a real scenario, you'd have specific binaries for arm64 and amd64, or a universal one.
# For now, let's assume there's a zip containing the app bundle as per README.
ASSET_NAME="Source-game-builder-tool-macos.zip"

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$TAG/$ASSET_NAME"

printf "${BLUE}==>${NC} Downloading $ASSET_NAME...\n"
curl -fsSL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$ASSET_NAME"

printf "${BLUE}==>${NC} Unzipping...\n"
unzip -q -o "$INSTALL_DIR/$ASSET_NAME" -d "$INSTALL_DIR"

# Clean up zip
rm "$INSTALL_DIR/$ASSET_NAME"

# Path to the binary inside the app bundle
BINARY_PATH="$INSTALL_DIR/Source-game-builder-tool-macos.app/Contents/MacOS/source-game-builder-tool"

if [ ! -f "$BINARY_PATH" ]; then
    # Maybe it was just a raw binary?
    BINARY_PATH="$INSTALL_DIR/source-game-builder-tool"
fi

if [ -f "$BINARY_PATH" ]; then
    chmod +x "$BINARY_PATH"
    printf "${GREEN}==>${NC} Successfully installed to $INSTALL_DIR\n"
    printf "${BLUE}==>${NC} Starting the builder...\n\n"

    # Run the binary
    "$BINARY_PATH" "$@"

    # After the tool runs, it might have deleted itself (if it was the binary).
    # However, since it's in an app bundle inside $INSTALL_DIR, let's see if it's still there.

    if [ -d "$INSTALL_DIR" ]; then
        # If the app bundle is gone, let's remove the whole install dir if it's empty
        # or just notify the user.
        if [ ! -f "$BINARY_PATH" ]; then
            rm -rf "$INSTALL_DIR"
            printf "\n${GREEN}==>${NC} Cleanup complete. Tool and temporary files removed.\n"
        else
            printf "\n${BLUE}==>${NC} Installation process finished.\n"
            printf "${BLUE}==>${NC} You can find the builder at: $INSTALL_DIR\n"
        fi
    fi
else
    printf "${RED}Error:${NC} Could not find the executable binary in the downloaded package.\n"
    exit 1
fi
