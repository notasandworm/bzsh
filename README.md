# 󰅶 Busy Shell (bzsh)

BzSh configures a beautiful prompt, sets up handy autocomplete, tunes history, and plugs in modern CLI tools like `eza`, `bat`, and `zoxide`.

Now powered by a fast, single Go binary with embedded scripts and native distro detection (Debian/Ubuntu & Arch Linux)!

---

## 🚀 Quick Install (One-Liner)

If you're on a fresh system and want to get set up immediately:

```bash
# 1. Download the bzsh binary
curl -fsSL https://raw.githubusercontent.com/notasandworm/bzsh/main/install.sh | sh

# 2. Run setup to configure your shell
bzsh setup
```

*(Note: `bzsh setup` runs non-interactively with auto-yes by default. Run `bzsh setup -i` if you prefer to toggle features interactively!)*

---

## 🛠️ Manual Installation (Git Clone or Local Build)

If you prefer building locally with Go:

```bash
# Clone the repository
git clone git@github.com:notasandworm/bzsh.git
cd bzsh

# Build the binary
go build -o bzsh main.go

# Run setup (non-interactive auto-yes by default)
./bzsh setup

# Or run interactively if you want to toggle individual features:
./bzsh setup -i
```

Once installed, the `bzsh` binary is available in your PATH at `~/.local/bin/bzsh`.

---

## ✨ Features

- **Custom Prompt Decorator**: A clean prompt showing your current directory, Git status, active runtime/SDK version (Go, Python, Rust, Node, Docker), and command execution duration.
- **Smart Autocomplete**: Colorized completions, partial-word matching, and process helper completions for commands like `kill`.
- **Sensible History Settings**: Command history sharing across terminal tabs, timestamps, and deduplication.
- **Keybindings**: Prefix history search with arrow keys (type `git` and press Up to search `git` commands).
- **Better Shell Defaults**: Suppresses annoyances like system beeps and enables automatic `cd`.
- **Modern CLI Tools (Distro-Aware)**:
  - `eza` aliases (`ls`, `la`, `ll` with icons and Git statuses).
  - `bat` syntax highlighting (`batcat` on Debian, `bat` on Arch).
  - `zoxide` integration for jumping around folders fast.
  - Automatically handles Debian's `fd-find` naming quirk while maintaining native `fd` on Arch.
- **Neovim Updater**: Installs an `update-nvim` helper function to pull and link the latest Neovim AppImage directly.

---

## 🔄 Keeping Things Updated

To update your scripts and binary to the latest version:

```bash
bzsh update
```

---

## 🧹 Zero-Footprint Uninstall

If you decide `bzsh` isn't for you:

```bash
bzsh uninstall
```

This will cleanly strip the `bzsh` configuration block from your `.zshrc` and prompt to wipe the config folder and binary for a 100% clean footprint.

---

## ✉️ Queries & Support

If you run into issues, have questions, or just want to chat about Zsh configs, drop me an email at **notasandworm@gmail.com**!
