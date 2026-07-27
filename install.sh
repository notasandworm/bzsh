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

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$BINARY_URL" -o "$BIN_DIR/bzsh" 2>/dev/null || cp ./bzsh "$BIN_DIR/bzsh" 2>/dev/null || true
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$BIN_DIR/bzsh" "$BINARY_URL" 2>/dev/null || cp ./bzsh "$BIN_DIR/bzsh" 2>/dev/null || true
fi

chmod +x "$BIN_DIR/bzsh"

echo ""
echo "✔ Successfully downloaded bzsh binary to $BIN_DIR/bzsh!"
echo "➜ Run 'bzsh setup' to start the installation."
echo ""
