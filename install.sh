#!/bin/bash
set -euo pipefail

REPO="shubh-io/dockmate"
BINARY_NAME="dockmate"
INSTALL_DIR="/usr/local/bin"

# Better architecture detection
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) RELEASE_ARCH="amd64" ;;
    aarch64|arm64) RELEASE_ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "==> Preparing to install dockmate from GitHub releases..."

# Better JSON parsing - fetch entire response first
API_URL="https://api.github.com/repos/$REPO/releases/latest"
echo "==> Checking GitHub for the latest release..."

# Download the full API response to avoid pipe issues
API_RESPONSE=$(curl -fsSL "$API_URL" 2>&1) || {
    echo "Error: Failed to fetch release info from GitHub"
    echo "This might be due to rate limiting or network issues"
    exit 1
}

# Parse tag name more reliably
LATEST_TAG=$(echo "$API_RESPONSE" | grep -o '"tag_name": *"[^"]*"' | head -1 | sed 's/.*"\([^"]*\)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "Error: Could not determine latest release version"
    echo "GitHub API might be rate limited. Try again in a few minutes."
    exit 1
fi

echo "✔ Latest version found: $LATEST_TAG"

ASSET_NAME="dockmate-linux-${RELEASE_ARCH}"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$ASSET_NAME"

echo "==> Downloading release binary..."
echo "==> From: $DOWNLOAD_URL"

TMP_BIN=$(mktemp /tmp/dockmate.XXXXXX)

# Download with better error handling
if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_BIN"; then
    echo "Error: Failed to download binary"
    rm -f "$TMP_BIN"
    exit 1
fi

# Checksum verification (optional)
CHECKSUM_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$ASSET_NAME.sha256"
if curl -fsSL -o "$TMP_BIN.sha256" "$CHECKSUM_URL" 2>/dev/null; then
    echo "==> Verifying checksum..."
    cd $(dirname "$TMP_BIN")
    if sha256sum -c "$TMP_BIN.sha256" 2>/dev/null | grep -q OK; then
        echo "✔ Checksum verified"
    else
        echo "Warning: Checksum verification failed"
    fi
else
    echo "==> No checksum file found for this release; skipping verification."
fi

chmod +x "$TMP_BIN"

echo "==> Installing dockmate to $INSTALL_DIR..."

# Use sudo only if needed
if [ "$(id -u)" -eq 0 ]; then
    mkdir -p "$INSTALL_DIR"
    mv "$TMP_BIN" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo mkdir -p "$INSTALL_DIR"
    sudo mv "$TMP_BIN" "$INSTALL_DIR/$BINARY_NAME"
fi

echo "✔ dockmate is installed to $INSTALL_DIR/$BINARY_NAME."
echo ""
echo "Next steps:"
echo "  - Run: dockmate"
echo "  - Make sure $INSTALL_DIR is on your PATH if needed."
echo ""
echo "To update later, use 'dockmate update' or re-run this installer."
