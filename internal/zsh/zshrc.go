package zsh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/notasandworm/bzsh/internal/ui"
)

const (
	StartMarker       = "# >>> bzsh initialize >>>"
	EndMarker         = "# <<< bzsh initialize <<<"
	LegacyStartMarker = "# >>> bzhrc initialize >>>"
	LegacyEndMarker   = "# <<< bzhrc initialize <<<"
)


// ConfigOptions holds feature toggle choices selected during setup.
type ConfigOptions struct {
	PromptDecorator     bool
	Autocomplete        bool
	HistorySettings     bool
	KeybindingsSearch   bool
	ShellOptions        bool
	EzaAliases          bool
	BatAlias            string // Distro-specific bat alias string (e.g. alias bat="batcat" or "")
	Zoxide              bool
	FdSymlink           bool
	FdTarget            string
	NvimUpdater         bool
	AddToPath           bool
}

// ConfigPaths returns standard directory locations for bzsh.
type ConfigPaths struct {
	ConfigDir string // ~/.config/bzsh
	BinDir    string // ~/.local/bin
	ZshrcFile string // ~/.zshrc
}

// GetDefaultPaths calculates standard paths in the user's home directory.
func GetDefaultPaths() (*ConfigPaths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not resolve home directory: %w", err)
	}

	return &ConfigPaths{
		ConfigDir: filepath.Join(home, ".config", "bzsh"),
		BinDir:    filepath.Join(home, ".local", "bin"),
		ZshrcFile: filepath.Join(home, ".zshrc"),
	}, nil
}

// EnsureDirectories creates ~/.config/bzsh and ~/.local/bin if missing.
func EnsureDirectories(paths *ConfigPaths) error {
	if err := os.MkdirAll(paths.ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir %s: %w", paths.ConfigDir, err)
	}
	if err := os.MkdirAll(paths.BinDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin dir %s: %w", paths.BinDir, err)
	}
	return nil
}

// BackupZshrc creates a timestamped copy of .zshrc before we touch anything!
// Safety first, tinkering second.
func BackupZshrc(paths *ConfigPaths) error {
	if _, err := os.Stat(paths.ZshrcFile); os.IsNotExist(err) {
		return nil // Nothing to back up if .zshrc doesn't exist yet!
	}

	timestamp := time.Now().Format("20060102150405")
	backupPath := fmt.Sprintf("%s.bzsh-backup.%s", paths.ZshrcFile, timestamp)

	ui.PrintStep(fmt.Sprintf("Backing up existing %s to %s...", paths.ZshrcFile, backupPath))

	data, err := os.ReadFile(paths.ZshrcFile)
	if err != nil {
		return fmt.Errorf("failed to read %s for backup: %w", paths.ZshrcFile, err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file %s: %w", backupPath, err)
	}

	ui.PrintOK("Backup created successfully.")
	return nil
}

// WriteEmbeddedScript extracts an embedded script file to ~/.config/bzsh/
func WriteEmbeddedScript(filename string, destDir string) error {
	data, err := GetEmbeddedScript(filename)
	if err != nil {
		return err
	}

	destPath := filepath.Join(destDir, filename)
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed writing %s to %s: %w", filename, destPath, err)
	}
	return nil
}

// GenerateConfigBlock builds the full shell script block to insert into .zshrc.
func GenerateConfigBlock(opts ConfigOptions, paths *ConfigPaths) (string, error) {
	var sb strings.Builder

	sb.WriteString(StartMarker + "\n")
	sb.WriteString("# This block is managed by bzsh. Do not edit manually.\n\n")

	// 1. Prompt Decorator
	if opts.PromptDecorator {
		if err := WriteEmbeddedScript("prompt.bzsh", paths.ConfigDir); err != nil {
			return "", err
		}
		promptPath := filepath.Join(paths.ConfigDir, "prompt.bzsh")
		sb.WriteString(fmt.Sprintf("# Busy Shell (bzsh) Prompt Decorator\n[[ -f \"%s\" ]] && source \"%s\"\n\n", promptPath, promptPath))
	} else {
		sb.WriteString("# Simple prompt setup\nautoload -Uz promptinit && promptinit\nprompt adam1\n\n")
	}

	// 2. Autocomplete
	if opts.Autocomplete {
		if err := WriteEmbeddedScript("autocomplete.bzsh", paths.ConfigDir); err != nil {
			return "", err
		}
		autoPath := filepath.Join(paths.ConfigDir, "autocomplete.bzsh")
		sb.WriteString(fmt.Sprintf("# Busy Shell (bzsh) Autocomplete\n[[ -f \"%s\" ]] && source \"%s\"\n\n", autoPath, autoPath))
	}

	// 3. History Settings
	if opts.HistorySettings {
		sb.WriteString("# History configuration\n" +
			"HISTFILE=\"$HOME/.zsh_history\"\n" +
			"HISTSIZE=50000\n" +
			"SAVEHIST=50000\n" +
			"setopt APPEND_HISTORY          # Append history to file rather than overwrite\n" +
			"setopt SHARE_HISTORY           # Share history across shell instances immediately\n" +
			"setopt EXTENDED_HISTORY        # Save timestamp and execution duration\n" +
			"setopt HIST_IGNORE_ALL_DUPS    # Purge duplicate entries when adding new ones\n" +
			"setopt HIST_IGNORE_SPACE       # Do not record commands preceded by a space\n\n")
	}

	// 4. Keybindings & History Search
	if opts.KeybindingsSearch {
		sb.WriteString("# Keybindings & History Search\n" +
			"autoload -Uz up-line-or-beginning-search down-line-or-beginning-search\n" +
			"zle -N up-line-or-beginning-search\n" +
			"zle -N down-line-or-beginning-search\n" +
			"bindkey '^[[A' up-line-or-beginning-search\n" +
			"bindkey '^[[B' down-line-or-beginning-search\n\n")
	}

	// 5. Shell Options
	if opts.ShellOptions {
		sb.WriteString("# Shell options\n" +
			"setopt AUTO_CD                 # Change directory without typing 'cd'\n" +
			"setopt NO_BEEP                 # Suppress terminal audio alerts\n\n")
	}

	// 6. Eza Aliases
	if opts.EzaAliases {
		sb.WriteString("# Eza Aliases\n" +
			"if command -v eza &>/dev/null; then\n" +
			"  alias ls=\"eza --icons=always --colour=always\"\n" +
			"  alias la=\"eza -alog --git --icons=always --colour=always\"\n" +
			"  alias ll=\"eza -log --git --icons=always --colour=always\"\n" +
			"fi\n\n")
	}

	// 7. Bat Alias (distro specific)
	if opts.BatAlias != "" {
		sb.WriteString(fmt.Sprintf("# Bat/Batcat alias\nif command -v batcat &>/dev/null; then\n  %s\nfi\n\n", opts.BatAlias))
	}

	// 8. Zoxide Init
	if opts.Zoxide {
		sb.WriteString("# Zoxide init\n" +
			"if command -v zoxide &>/dev/null; then\n" +
			"  eval \"$(zoxide init --cmd z zsh)\"\n" +
			"fi\n\n")
	}

	// 9. Neovim Updater Helper
	if opts.NvimUpdater {
		if err := WriteEmbeddedScript("nvim-update.bzsh", paths.ConfigDir); err != nil {
			return "", err
		}
		nvimPath := filepath.Join(paths.ConfigDir, "nvim-update.bzsh")
		sb.WriteString(fmt.Sprintf("# Neovim Updater Helper\n[[ -f \"%s\" ]] && source \"%s\"\n\n", nvimPath, nvimPath))
	}

	// 10. Local bin path
	if opts.AddToPath {
		sb.WriteString(fmt.Sprintf("# Include local bin in PATH\nexport PATH=\"$PATH:%s\"\n\n", paths.BinDir))
	}

	sb.WriteString(EndMarker)
	return sb.String(), nil
}

// ExtractUserContent parses a .zshrc string for bzsh or bzhrc markers
// and returns user content prior to the start marker and user content following the end marker.
func ExtractUserContent(content string) (before string, after string) {
	startPairs := []struct{ start, end string }{
		{StartMarker, EndMarker},
		{LegacyStartMarker, LegacyEndMarker},
	}

	for _, pair := range startPairs {
		if strings.Contains(content, pair.start) && strings.Contains(content, pair.end) {
			startIndex := strings.Index(content, pair.start)
			endIndex := strings.Index(content, pair.end) + len(pair.end)

			before = content[:startIndex]
			after = content[endIndex:]
			return before, after
		}
	}

	return content, ""
}

// FindLatestBackup finds the path to the most recent .zshrc.bzsh-backup.* file if present.
func FindLatestBackup(paths *ConfigPaths) (string, error) {
	dir := filepath.Dir(paths.ZshrcFile)
	pattern := filepath.Join(dir, ".zshrc.bzsh-backup.*")
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return "", err
	}

	var latestFile string
	var latestTime time.Time

	for _, match := range matches {
		info, err := os.Stat(match)
		if err == nil {
			if latestFile == "" || info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
				latestFile = match
			}
		}
	}

	return latestFile, nil
}

// UpdateZshrc replaces or appends the bzsh config block in .zshrc,
// preserving any lines existing before and after the bzsh/bzhrc indicators.
func UpdateZshrc(paths *ConfigPaths, newBlock string) error {
	ui.PrintStep(fmt.Sprintf("Updating %s...", paths.ZshrcFile))

	var existingContent string
	if data, err := os.ReadFile(paths.ZshrcFile); err == nil && len(data) > 0 {
		existingContent = string(data)
	} else {
		if backupPath, err := FindLatestBackup(paths); err == nil && backupPath != "" {
			if data, err := os.ReadFile(backupPath); err == nil {
				ui.PrintStep(fmt.Sprintf("Restoring pre/post indicator content from latest backup %s...", backupPath))
				existingContent = string(data)
			}
		}
	}

	userBefore, userAfter := ExtractUserContent(existingContent)

	var sb strings.Builder
	if strings.TrimSpace(userBefore) != "" {
		sb.WriteString(strings.TrimRight(userBefore, "\n"))
		sb.WriteString("\n\n")
	}

	sb.WriteString(newBlock)

	if strings.TrimSpace(userAfter) != "" {
		sb.WriteString("\n\n")
		sb.WriteString(strings.TrimLeft(userAfter, "\n"))
		if !strings.HasSuffix(userAfter, "\n") {
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("\n")
	}

	if err := os.WriteFile(paths.ZshrcFile, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed writing to %s: %w", paths.ZshrcFile, err)
	}

	ui.PrintOK("Updated bzsh configuration block in .zshrc!")
	return nil
}

// RemoveZshrcBlock strips the bzsh block from .zshrc during uninstall.
func RemoveZshrcBlock(paths *ConfigPaths) error {
	data, err := os.ReadFile(paths.ZshrcFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	content := string(data)
	before, after := ExtractUserContent(content)
	if before == content && after == "" {
		ui.PrintWarn(fmt.Sprintf("No bzsh block found in %s.", paths.ZshrcFile))
		return nil
	}

	ui.PrintStep(fmt.Sprintf("Removing bzsh block from %s...", paths.ZshrcFile))

	var sb strings.Builder
	if strings.TrimSpace(before) != "" {
		sb.WriteString(strings.TrimRight(before, "\n"))
		sb.WriteString("\n")
	}
	if strings.TrimSpace(after) != "" {
		sb.WriteString(strings.TrimLeft(after, "\n"))
	}

	if err := os.WriteFile(paths.ZshrcFile, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to save stripped %s: %w", paths.ZshrcFile, err)
	}
	ui.PrintOK("Removed bzsh block.")
	return nil
}

