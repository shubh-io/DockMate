#!/usr/bin/env bash
#
# Dockmate installer script
# Repository: https://github.com/shubh-io/DockMate
#
# This script downloads the latest Dockmate release binary from GitHub
# and installs it into /usr/local/bin (or a custom directory).
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/shubh-io/DockMate/main/install.sh | bash
#
# Optional env vars:
#   INSTALL_DIR=/some/dir   # default: /usr/local/bin
#   REPO=shubh-io/dockmate  # override repo if forked
#

set -euo pipefail



#config
REPO="shubh-io/dockmate"
BINARY_NAME="dockmate"
INSTALL_DIR="/usr/local/bin"


# Logging helpers (friendly, minimal)

info()    { echo "==> $*"; }
success() { echo "✔ $*"; }
error()   { echo "Error: $*" >&2; }


# Pre-flight checks

if ! command -v curl >/dev/null 2>&1; then
  error "curl is required but not installed. Please install curl and try again."
  exit 1
fi

if [ "$(id -u)" -ne 0 ] && ! command -v sudo >/dev/null 2>&1; then
  error "I need permission to move the binary into ${INSTALL_DIR}."
  error "You are not root and 'sudo' is not available."
  error "Run this as root or install sudo, then try again."
  exit 1
fi


# Detect architecture

ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)
    RELEASE_ARCH="amd64"
    ;;
  *)
    error "Unsupported architecture: ${ARCH}"
    error "Currently this installer only supports x86_64 (amd64)."
    exit 1
    ;;
esac


# Fetch latest release tag

info "Preparing to install ${BINARY_NAME} from GitHub releases..."
API_URL="https://api.github.com/repos/${REPO}/releases/latest"

info "Checking GitHub for the latest release..."
LATEST_TAG="$(curl -fsSL "${API_URL}" | grep -m1 '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')"

if [ -z "${LATEST_TAG}" ]; then
  error "Could not determine the latest release tag from GitHub."
  error "Make sure the repository has at least one published release."
  exit 1
fi

success "Latest version found: ${LATEST_TAG}"


# Build download URL

# Expected asset name: dockmate-linux-amd64
ASSET_NAME="${BINARY_NAME}-linux-${RELEASE_ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ASSET_NAME}"

info "Downloading release binary..."
info "From: ${DOWNLOAD_URL}"

TMP_BIN="$(mktemp "/tmp/${BINARY_NAME}.XXXXXX")"

if ! curl -fsSL "${DOWNLOAD_URL}" -o "${TMP_BIN}"; then
  error "Download failed."
  error "Check that the release ${LATEST_TAG} includes an asset named '${ASSET_NAME}'."
  rm -f "${TMP_BIN}"
  exit 1
fi


# Verify checksum (if provided by the release)

# We expect a release asset named "${ASSET_NAME}.sha256" containing a single
# checksum line. If it's not present, we skip verification and continue.
CHECKSUM_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ASSET_NAME}.sha256"
if curl -fsSL -o "${TMP_BIN}.sha256" "${CHECKSUM_URL}"; then
  info "Verifying checksum..."
  expected="$(awk '{print $1}' "${TMP_BIN}.sha256" || true)"
  if [ -z "$expected" ]; then
    error "Checksum file is empty or malformed; aborting."
    rm -f "${TMP_BIN}" "${TMP_BIN}.sha256"
    exit 1
  fi

  if command -v sha256sum >/dev/null 2>&1; then
    actual="$(sha256sum "${TMP_BIN}" | awk '{print $1}')"
  else
    actual="$(shasum -a 256 "${TMP_BIN}" | awk '{print $1}')"
  fi

  if [ "$actual" != "$expected" ]; then
    error "Checksum verification failed."
    error "Expected: $expected" 
    error "Actual:   $actual"
    rm -f "${TMP_BIN}" "${TMP_BIN}.sha256"
    exit 1
  fi

  info "Checksum OK."
else
  info "No checksum file found for this release; skipping verification."
fi

chmod +x "${TMP_BIN}"


# Install binary

info "Installing ${BINARY_NAME} to ${INSTALL_DIR}..."

if [ "$(id -u)" -eq 0 ]; then
  # Already root
  mkdir -p "${INSTALL_DIR}"
  mv "${TMP_BIN}" "${INSTALL_DIR}/${BINARY_NAME}"
else
  # using sudo for privileged operations
  sudo mkdir -p "${INSTALL_DIR}"
  sudo mv "${TMP_BIN}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

success "${BINARY_NAME} is installed to ${INSTALL_DIR}/${BINARY_NAME}."

echo
echo "Next steps:"
echo "  - Run: ${BINARY_NAME}"
echo "  - Make sure ${INSTALL_DIR} is on your PATH if needed."
echo
echo "To update later, use 'dockmate update' or re-run this installer."
