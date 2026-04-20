#!/bin/bash

# macOS Source Builder Oneliner Installer
# This script downloads the latest release of the macOS Source Builder and runs it.

set -e

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
RELEASE_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")
TAG=$(echo "$RELEASE_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
    printf "${RED}Error:${NC} Could not find latest release. Please check $REPO\n"
    exit 1
fi

printf "${BLUE}==>${NC} Found version: ${GREEN}$TAG${NC}\n"
