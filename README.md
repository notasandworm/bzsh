# 󰅶 Busy Shell (bzsh)

BzSh configures a beautiful prompt, sets up handy autocomplete, tunes history, and plugs in modern CLI tools like `eza`, `bat`, and `zoxide`.

---

*(Note for Matt: `go build` after making changes!!!)*

## 💡 Why bzsh?

I started this project after using nerd font capable prompt decorators like [Starship](https://starship.rs/) and encountering a few issues:
1. There is a prompt offset bug with nerd fonts on the same line as the prompt starship where the prompt would be offset by a few characters.
2. When this happened, I found prompt artifacts were common and could only be cleared when clearing my terminal via `[^l]`.

`bzsh` was built as a lightweight personal alternative:
- **SSH & Remote Compatibility**: Built with native Zsh prompt expansions and hooks, eliminating prompt offset rendering bugs over SSH i encountered.
- **Fast & Self-Contained**: Powered by a fast Go binary with embedded Zsh scripts, providing automatic package detection for Debian/Ubuntu and Arch Linux.
- **All-in-One Shell Enhancement**: Installs a clean multi-line prompt decorator, smart autocomplete, synchronized history, arrow-key history search, and modern CLI tool integrations.

---

## 🚀 Quick Install 

Installer installs Bzsh, [Eza](https://github.com/eza-community/eza), [bat](https://github.com/sharkdp/bat), [zoxide](https://github.com/ajeetdsouza/zoxide) as needed

1. Download the bzsh binary
```bash
curl -fsSL https://raw.githubusercontent.com/notasandworm/bzsh/main/install.sh | sh
```

2. Run setup to configure your shell
```bash
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
  - `eza` aliases (`ls`, `la`, `ll` with dynamic icon support based on Nerd Font preference).
  - `bat` syntax highlighting (`batcat` on Debian, `bat` on Arch).
  - `zoxide` integration for jumping around folders fast.
  - Automatically handles Debian's `fd-find` naming quirk while maintaining native `fd` on Arch.
- **Nerd Fonts & Dual-Mode Symbol Support**: Rich icon mode with automatic clean text fallback for systems without Nerd Fonts installed (`bzsh font` command included).
- **Neovim Updater**: Installs an `update-nvim` helper function to pull and link the latest Neovim AppImage directly.

---

## 🔤 Nerd Font Support & Symbol Downloader

`bzsh` supports a dual-mode symbol architecture to ensure your prompt and `eza` look great with rich Nerd Font icons when available, while completely avoiding missing character boxes or squares (`[?]`) on basic terminals:

- **Nerd Font Mode (`BZSH_NERD_FONTS=1`)**: Renders rich glyphs for OS (`󰣇`, ``, `🐧`), Git (``), and project runtimes (Go ``, Python ``, Rust ``, Node ``, C++ ``, C ``, CMake/ESP-IDF ``, Docker `󰡨`).
- **Clean Text Mode (`BZSH_NERD_FONTS=0`)**: Automatically falls back to clean, readable text tags (`[ARCH]`, `[DEB]`, `+main`, `[C++]`, `[C]`, `[ESP-IDF]`, `[Go]`, etc.).

### Automatic Font Downloader
To download and install the official `SymbolsNerdFontMono` font into `~/.local/share/fonts/` and update your desktop font cache:

```bash
bzsh font
```

---

## 🎨 Customizing & Moving Prompt Components

`bzsh` manages prompt layout inside `~/.config/bzsh/prompt.bzsh`. The default multi-line prompt is built from modular helper components:

```zsh
PROMPT='%F{blue}%~%f$(_git_status)$(_sdk_info)${_CMD_ELAPSED}
$(_venv_info)%F{green}%n@%m%f$(_os_info)%# '
```

### Available Components

| Component | Description | Rich Icon Output (`BZSH_NERD_FONTS=1`) | Fallback Output (`BZSH_NERD_FONTS=0`) |
| :--- | :--- | :--- | :--- |
| `%F{blue}%~%f` | Current working directory | `~/projects/bzsh` | `~/projects/bzsh` |
| `$(_git_status)` | Active Git branch & dirty state | `  main*` | ` +main*` |
| `$(_sdk_info)` | Detected project runtime / SDK | `  v1.22`, `  C++`, `  ESP-IDF` | ` [Go v1.22]`, ` [C++]`, ` [ESP-IDF]` |
| `${_CMD_ELAPSED}` | Command execution time (if ≥0.1s) | ` took 1.2s` | ` took 1.2s` |
| `$(_venv_info)` | Python venv / Conda environment | `(myenv) ` | `(myenv) ` |
| `%F{green}%n@%m%f` | User and host name | `user@host` | `user@host` |
| `$(_os_info)` | OS / Distro badge | `󰣇 `, ` `, or `🐧` | `[ARCH]`, `[DEB]`, or `[LINUX]` |
| `%# ` | Prompt symbol (`%` for user, `#` for root) | `% ` | `% ` |

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
