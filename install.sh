#!/bin/sh
# ------------------------------------------------------------------------------
# 󰅶 Busy Shell (bzsh) - Binary Downloader
# ------------------------------------------------------------------------------
# Hey! This script downloads the pre-compiled bzsh binary directly into
# ~/.local/bin/bzsh so you don't have to compile anything manually.
#
# Once downloaded, just run:
#   bzsh setup

set -e

BIN_DIR="$HOME/.local/bin"
mkdir -p "$BIN_DIR"

echo "➜ Downloading Busy Shell (bzsh) binary..."

# Determine architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH_NAME="amd64" ;;
  aarch64|arm64) ARCH_NAME="arm64" ;;
  *)       ARCH_NAME="amd64" ;;
esac

BINARY_URL="https://github.com/notasandworm/bzsh/releases/latest/download/bzsh-linux-${ARCH_NAME}"
RAW_URL="https://raw.githubusercontent.com/notasandworm/bzsh/main/bzsh"

DOWNLOAD_SUCCESS=0

download_file() {
  url="$1"
  dest="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$dest" 2>/dev/null && return 0
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$dest" "$url" 2>/dev/null && return 0
  fi
  return 1
}

# 1. Attempt download from GitHub Release assets
if download_file "$BINARY_URL" "$BIN_DIR/bzsh"; then
  DOWNLOAD_SUCCESS=1
fi

# 2. Fallback: Attempt download from raw GitHub repository (main branch)
if [ "$DOWNLOAD_SUCCESS" -ne 1 ]; then
  if download_file "$RAW_URL" "$BIN_DIR/bzsh"; then
    DOWNLOAD_SUCCESS=1
  fi
fi

# 3. Fallback: Copy local workspace binary if available
if [ "$DOWNLOAD_SUCCESS" -ne 1 ] && [ -f ./bzsh ]; then
  echo "➜ Using local bzsh binary..."
  cp ./bzsh "$BIN_DIR/bzsh" && DOWNLOAD_SUCCESS=1
fi

if [ "$DOWNLOAD_SUCCESS" -ne 1 ] || [ ! -f "$BIN_DIR/bzsh" ]; then
  echo "❌ Error: Failed to download or locate bzsh binary." >&2
  exit 1
fi

chmod +x "$BIN_DIR/bzsh"

echo ""
echo "✔ Successfully installed bzsh binary to $BIN_DIR/bzsh!"
echo "➜ Run 'bzsh setup' to start the installation."
echo ""
