# 󰅶 Busy Shell (bzsh)

BzSh configures a beautiful prompt, sets up handy autocomplete, tunes history, and plugs in modern CLI tools like `eza`, `bat`, and `zoxide`.

Now powered by a fast, single Go binary with embedded scripts and native distro detection (Debian/Ubuntu & Arch Linux)!

---

*(Note for Matt: `go build` after making changes!!!)*

## 💡 Why bzsh?

I build started this project after encountering a few issues:
1. There was a prompt offset bug with starship where whenever i SSH'd into another machine, the prompt would be offset by a few characters.
2. I found prompt artifacts were common and could only be cleared when clearing my terminal via `^L`

`bzsh` was built as a lightweight personal alternative:
- **Flawless SSH & Remote Compatibility**: Built with native Zsh prompt expansions and hooks, eliminating prompt offset rendering bugs over SSH i encountered.
- **Fast & Self-Contained**: Powered by a fast Go binary with embedded Zsh scripts, providing automatic package detection for Debian/Ubuntu and Arch Linux.
- **All-in-One Shell Enhancement**: Installs a clean multi-line prompt decorator, smart autocomplete, synchronized history, arrow-key history search, and modern CLI tool integrations.

---

## 🚀 Quick Install (One-Liner)

If you're on a fresh system and want to get set up immediately:

```bash
# 1. Download the bzsh binary
curl -fsSL https://raw.githubusercontent.com/notasandworm/bzsh/main/install.sh | sh
```

```bash
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

## 🎨 Customizing & Moving Prompt Components

`bzsh` manages prompt layout inside `~/.config/bzsh/prompt.bzsh`. The default multi-line prompt is built from modular helper components:

```zsh
PROMPT='$(_os_info)%F{blue}%~%f $(_git_status)$(_sdk_info)${_CMD_ELAPSED}
$(_venv_info)%F{green}%n@%m%f%# '
```

### Available Components

| Component | Description | Example Output |
| :--- | :--- | :--- |
| `$(_os_info)` | OS / Distro badge | `%F{cyan}[ARCH]%f ` or `%F{red}[DEB]%f ` |
| `%F{blue}%~%f` | Current directory path | `~/projects/bzsh` |
| `$(_git_status)` | Active Git branch & dirty state | `%F{magenta}+main*%f ` |
| `$(_sdk_info)` | Detected project runtime / SDK | `%F{cyan}[Go v1.22]%f `, `[Py v3.12]`, etc. |
| `${_CMD_ELAPSED}` | Command execution time (if ≥0.1s) | `%F{yellow}took 1.2s%f` |
| `$(_venv_info)` | Python venv / Conda environment | `%F{yellow}(myenv)%f ` |
| `%F{green}%n@%m%f%# ` | User, host, and prompt symbol (`%` / `#`) | `user@host% ` |

### Example: Moving `[OS_info]` Before Prompt Symbol `%`

If you want to move prompt components around—for instance, moving `[OS_info]` from the top-left corner to right before the `%` prompt symbol on the second line—edit `~/.config/bzsh/prompt.bzsh` (or set `PROMPT` in your `.zshrc`):

```zsh
# Modified PROMPT with [OS_info] placed right before the prompt symbol '%'
PROMPT='%F{blue}%~%f $(_git_status)$(_sdk_info)${_CMD_ELAPSED}
$(_venv_info)%F{green}%n@%m%f $(_os_info)%# '
```

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
