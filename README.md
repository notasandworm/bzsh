# 󰅶 Busy Zsh (bzsh)

Hey there! Welcome to **Busy Zsh** (or `bzsh` for short). 

This is my personal playground for setting up a super-charged, comfortable Zsh terminal environment on fresh system installs. It configures a beautiful prompt, sets up handy autocomplete, tunes history, and plugs in a few of my favorite CLI tools (like `eza`, `bat`, and `zoxide`).

I built an interactive installer in Bash to handle everything. It checks your dependencies, backs up your existing configuration, and lets you toggle features on and off so you only get what you want.

---

## 🚀 Quick Install (One-Liner)

If you're on a fresh system and want to get set up immediately, just copy and paste this command into your terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/notasandworm/bzsh/main/bzsh -o /tmp/bzsh && chmod +x /tmp/bzsh && /tmp/bzsh setup && rm /tmp/bzsh
```

*(Note: If Zsh isn't installed yet, the script will notice and ask if you want to install it first!)*

---

## 🛠️ Manual Installation (Git Clone)

If you'd rather clone the repo and run it locally, you can do this:

```bash
# Clone the repository
git clone git@github.com:notasandworm/bzsh.git
cd bzsh

# Run the installer
./bzsh setup
```

Once installed, the `bzsh` command will be added to your path at `~/.local/bin/bzsh`, so you can run it from anywhere.

---

## ✨ Features You Can Toggle

- **Custom Prompt Decorator**: A clean prompt showing your current directory, Git status, active runtime/SDK version (Go, Python, Rust, Node), and command execution duration.
- **Smart Autocomplete**: Colorized completions, partial-word matching, and process helper completions for commands like `kill`.
- **Sensible History Settings**: Sets up command history sharing across terminal tabs, timestamps, and deduplication.
- **Keybindings**: Enables prefix history search with up/down arrow keys (type `git` and press Up to search only commands starting with `git`).
- **Better Shell Defaults**: Suppresses annoyances like system beeps and enables automatic `cd` (just type a directory name without `cd`).
- **Modern CLI Tools**:
  - `eza` aliases (`ls`, `la`, `ll` with icons and Git statuses).
  - `bat` syntax highlighting.
  - `zoxide` integration for jumping around folders fast.
  - Fixes Debian's `fd-find` naming quirk automatically.
- **Neovim Updater**: Installs an `update-nvim` helper function to pull and link the latest Neovim AppImage directly.

---

## 🔄 Keeping Things Updated

To update your scripts and aliases to the latest version, just run:

```bash
bzsh update
```

---

## 🧹 Zero-Footprint Uninstall

If you decide `bzsh` isn't for you, removing it is completely pain-free. Run:

```bash
bzsh uninstall
```

This will cleanly strip the `bzsh` configuration block from your `.zshrc` and ask if you want to wipe the config folder and commands for a 100% clean footprint.

---

## ✉️ Queries & Support

If you run into issues, have questions, or just want to chat about Zsh configs, drop me an email at **notasandworm@gmail.com**!
