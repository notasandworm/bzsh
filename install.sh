#!/bin/sh
# ------------------------------------------------------------------------------
# 󰅶 Busy Shell (bzsh) - Installer
# ------------------------------------------------------------------------------
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/notasandworm/bzsh/main/install.sh | sh
#   curl -fsSL ... | sh -s -- -i          # Interactive setup
#   curl -fsSL ... | sh -s -- --no-setup  # Download binary only

set -e

BIN_DIR="$HOME/.local/bin"
mkdir -p "$BIN_DIR"

echo "➜ Downloading Busy Shell (bzsh) binary..."

# Determine architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)        ARCH_NAME="amd64" ;;
  aarch64|arm64) ARCH_NAME="arm64" ;;
  *)             ARCH_NAME="amd64" ;;
esac

BRANCH="${BZSH_BRANCH:-main}"
BINARY_URL="https://github.com/notasandworm/bzsh/releases/latest/download/bzsh-linux-${ARCH_NAME}"
RAW_URL="https://raw.githubusercontent.com/notasandworm/bzsh/${BRANCH}/bzsh"

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

# 1. Local workspace fallback (for testing inside repo)
if [ -f "./bzsh" ]; then
  echo "➜ Using local bzsh binary..."
  cp ./bzsh "$BIN_DIR/bzsh" && DOWNLOAD_SUCCESS=1
fi

# 2. Local Go build fallback (if running inside repo and go is available)
if [ "$DOWNLOAD_SUCCESS" -ne 1 ] && [ -f "./main.go" ] && command -v go >/dev/null 2>&1; then
  echo "➜ Compiling bzsh from local source..."
  if go build -o "$BIN_DIR/bzsh" main.go; then
    DOWNLOAD_SUCCESS=1
  fi
fi

# 3. If BZSH_BRANCH is explicitly set (e.g. BZSH_BRANCH=dev), try raw branch binary
if [ "$DOWNLOAD_SUCCESS" -ne 1 ] && [ -n "$BZSH_BRANCH" ]; then
  echo "➜ Fetching bzsh binary from '${BZSH_BRANCH}' branch..."
  if download_file "$RAW_URL" "$BIN_DIR/bzsh"; then
    DOWNLOAD_SUCCESS=1
  fi
fi

# 4. Attempt download from GitHub Release assets (standard release install)
if [ "$DOWNLOAD_SUCCESS" -ne 1 ]; then
  if download_file "$BINARY_URL" "$BIN_DIR/bzsh"; then
    DOWNLOAD_SUCCESS=1
  elif download_file "https://github.com/notasandworm/bzsh/releases/latest/download/bzsh" "$BIN_DIR/bzsh"; then
    DOWNLOAD_SUCCESS=1
  fi
fi

# 5. Fallback: Attempt download from raw GitHub repository
if [ "$DOWNLOAD_SUCCESS" -ne 1 ]; then
  if download_file "$RAW_URL" "$BIN_DIR/bzsh"; then
    DOWNLOAD_SUCCESS=1
  fi
fi

if [ "$DOWNLOAD_SUCCESS" -ne 1 ] || [ ! -f "$BIN_DIR/bzsh" ]; then
  echo "❌ Error: Failed to download or build bzsh binary." >&2
  exit 1
fi

chmod +x "$BIN_DIR/bzsh"

INSTALLED_VER=$("$BIN_DIR/bzsh" --version 2>/dev/null | grep -i 'bzsh version' | head -n 1)
if [ -n "$INSTALLED_VER" ]; then
  echo "✔ Successfully installed $INSTALLED_VER to $BIN_DIR/bzsh!"
else
  echo "✔ Successfully installed bzsh binary to $BIN_DIR/bzsh!"
fi

# Parse flags passed to install.sh
NO_SETUP=0
INTERACTIVE=0

for arg in "$@"; do
  case "$arg" in
    --no-setup|-s) NO_SETUP=1 ;;
    -i|--interactive) INTERACTIVE=1 ;;
  esac
done

if [ "$NO_SETUP" -eq 1 ]; then
  echo "➜ Binary installed at $BIN_DIR/bzsh."
  echo "➜ To configure your shell at any time, run: $BIN_DIR/bzsh setup"
  echo ""
  exit 0
fi

echo ""
echo "➜ Running bzsh setup..."

# Pass /dev/tty to setup if available so interactive prompts & sudo work cleanly when piped
SETUP_TTY=""
if [ -r /dev/tty ] && (exec < /dev/tty) 2>/dev/null; then
  SETUP_TTY="/dev/tty"
fi

if [ -n "$SETUP_TTY" ]; then
  if [ "$INTERACTIVE" -eq 1 ]; then
    "$BIN_DIR/bzsh" setup -i <"$SETUP_TTY"
  else
    "$BIN_DIR/bzsh" setup <"$SETUP_TTY"
  fi
else
  if [ "$INTERACTIVE" -eq 1 ]; then
    "$BIN_DIR/bzsh" setup -i
  else
    "$BIN_DIR/bzsh" setup
  fi
fi
