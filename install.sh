#!/bin/bash

# Configuration
REPO="cainseing/drop-cli"
INSTALL_PATH="/usr/local/bin/drop"
# Points to the raw file in the main branch's bin folder
BASE_URL="https://raw.githubusercontent.com/$REPO/main/bin"

echo "Checking system compatibility..."

# Detect OS
OS_TYPE=$(uname -s | tr '[:upper:]' '[:lower:]')
OS="linux"
[[ "$OS_TYPE" == "darwin" ]] && OS="darwin"

if [[ "$OS_TYPE" != "darwin" && "$OS_TYPE" != "linux" ]]; then
    echo "Error: Unsupported OS ($OS_TYPE)"
    exit 1
fi

# Detect Architecture
ARCH_TYPE=$(uname -m)
ARCH="amd64"
[[ "$ARCH_TYPE" == "arm"* || "$ARCH_TYPE" == "aarch64" ]] && ARCH="arm64"

if [[ "$ARCH_TYPE" != "x86_64" && "$ARCH_TYPE" != "arm"* && "$ARCH_TYPE" != "aarch64" ]]; then
    echo "Error: Unsupported architecture ($ARCH_TYPE)"
    exit 1
fi

# Construct URL
# Example: https://raw.githubusercontent.com/cainseing/drop-cli/main/bin/drop-linux-amd64
BINARY_NAME="drop-$OS-$ARCH"
BINARY_URL="$BASE_URL/$BINARY_NAME"

echo "Downloading $BINARY_NAME..."

# Download
curl -L -o "/tmp/drop_binary" "$BINARY_URL"

if [ $? -ne 0 ]; then
    echo "Error: Download failed. Verify that $BINARY_NAME exists in the /bin folder of the repo."
    exit 1
fi

# Permission and Install
chmod +x /tmp/drop_binary

echo "Installing to $INSTALL_PATH (requires sudo)..."
if ! sudo mv /tmp/drop_binary "$INSTALL_PATH"; then
    echo "Error: Installation failed during move."
    exit 1
fi

echo "-------------------------------------------"
echo "Success! 'drop' has been installed."
echo "Run 'drop --help' to get started."
echo "-------------------------------------------"