#!/bin/sh
# ------------------------------------------------------------------------------
# 󰅶 Busy Zsh (bzsh) - Quick Installer
# ------------------------------------------------------------------------------
# Hey! This script fetches the compiled bzsh binary and runs setup.
# Minimal, fast, and easy to read.

set -e

BIN_DIR="$HOME/.local/bin"
mkdir -p "$BIN_DIR"

echo "➜ Fetching bzsh binary..."

# Determine architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH_NAME="amd64" ;;
  aarch64|arm64) ARCH_NAME="arm64" ;;
  *)       ARCH_NAME="amd64" ;;
esac

# Download binary release from GitHub (or local copy if running inside git clone)
BINARY_URL="https://github.com/notasandworm/bzsh/releases/latest/download/bzsh-linux-${ARCH_NAME}"

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$BINARY_URL" -o "$BIN_DIR/bzsh" 2>/dev/null || cp ./bzsh "$BIN_DIR/bzsh" 2>/dev/null || true
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$BIN_DIR/bzsh" "$BINARY_URL" 2>/dev/null || cp ./bzsh "$BIN_DIR/bzsh" 2>/dev/null || true
fi

chmod +x "$BIN_DIR/bzsh"

echo "✔ Installed bzsh binary to $BIN_DIR/bzsh"
echo "➜ Launching bzsh setup..."

"$BIN_DIR/bzsh" setup
