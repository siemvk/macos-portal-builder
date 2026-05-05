#!/bin/bash

# macOS Source Builder Oneliner Installer
# This script downloads the latest release of the macOS Source Builder and runs it.

set -euo pipefail

# repo info
REPO="siemvk/macos-portal-builder"
INSTALL_DIR="$HOME/.macos-source-builder"
BINARY_NAME="source-game-builder-tool"

# colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

printf "${BLUE}==>${NC} macOS Source Builder Installer\n"

# check if running on macOS (for the people that somehow manage to think this will work on windows or linux?)
if [[ "$OSTYPE" != "darwin"* ]]; then
    printf "${RED}Error:${NC} This tool is only for macOS.\n"
    exit 1
fi

# create install directory
mkdir -p "$INSTALL_DIR"

printf "${BLUE}==>${NC} Fetching latest release information...\n"
RELEASE_JSON=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest")
TAG=$(echo "$RELEASE_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
    printf "${RED}Error:${NC} Could not find latest release. Please check $REPO\n"
    exit 1
fi

printf "${BLUE}==>${NC} Found version: ${GREEN}$TAG${NC}\n"

ASSET_NAME="Source-game-builder-tool-macos.zip"
CHECKSUM_NAME="${ASSET_NAME}.sha256"

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$TAG/$ASSET_NAME"
CHECKSUM_URL="https://github.com/$REPO/releases/download/$TAG/$CHECKSUM_NAME"

printf "${BLUE}==>${NC} Downloading $ASSET_NAME...\n"
curl -fsSL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$ASSET_NAME"

printf "${BLUE}==>${NC} Downloading checksum...\n"
if curl -fsSL "$CHECKSUM_URL" -o "$INSTALL_DIR/$CHECKSUM_NAME"; then
    printf "${BLUE}==>${NC} Verifying checksum...\n"
    # Navigate to INSTALL_DIR so shasum can find the file by its relative name
    cd "$INSTALL_DIR" || exit 1
    if ! shasum -a 256 -c "$CHECKSUM_NAME"; then
        printf "${RED}Error:${NC} Checksum verification failed. The downloaded file might be corrupted or compromised.\n"
        exit 1
    fi
    cd - > /dev/null || exit 1
else
    printf "${RED}Error:${NC} Could not download checksum file. Cannot verify the integrity of the downloaded package.\n"
    exit 1
fi

printf "${BLUE}==>${NC} Unzipping...\n"
unzip -q -o "$INSTALL_DIR/$ASSET_NAME" -d "$INSTALL_DIR"

# clean up zip
rm "$INSTALL_DIR/$ASSET_NAME"

# path to the binary inside the app bundle
BINARY_PATH="$INSTALL_DIR/Source-game-builder-tool-macos.app/Contents/MacOS/source-game-builder-tool"

if [ ! -f "$BINARY_PATH" ]; then
    BINARY_PATH="$INSTALL_DIR/source-game-builder-tool"
fi

if [ -f "$BINARY_PATH" ]; then
    chmod +x "$BINARY_PATH"
    printf "${GREEN}==>${NC} Successfully installed to $INSTALL_DIR\n"
    printf "${BLUE}==>${NC} Starting the builder...\n\n"

    # run binary
    "$BINARY_PATH" "$@"

    if [ -d "$INSTALL_DIR" ]; then
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
