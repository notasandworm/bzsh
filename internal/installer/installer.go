package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/notasandworm/bzsh/internal/distro"
	"github.com/notasandworm/bzsh/internal/ui"
	"github.com/notasandworm/bzsh/internal/zsh"
)

// RunSetup manages the entire installation flow.
// By default, autoYes is true unless the user specifies --interactive (-i).
func RunSetup(interactive bool) error {
	autoYes := !interactive

	ui.PrintTitle()

	// Initial confirmation prompt before tweaking anything!
	if !ui.AskYesNo("This will install custom prompt decorators, shell configurations, and quality-of-life CLI packages. Continue?", true, autoYes) {
		ui.PrintWarn("Setup aborted by user.")
		return nil
	}

	ui.PrintStep("Starting setup...")

	adapter := distro.DetectDistro()
	ui.PrintOK(fmt.Sprintf("Detected distribution: %s", adapter.Name()))

	paths, err := zsh.GetDefaultPaths()
	if err != nil {
		return err
	}

	// 1. Zsh check & installation
	missingZsh, err := adapter.CheckMissingPackages([]string{"zsh"})
	if err != nil {
		return err
	}
	if len(missingZsh) > 0 {
		ui.PrintWarn("Zsh is not installed! Zsh is required for Busy Zsh to function.")
		if ui.AskYesNo("Would you like to install Zsh using your package manager?", true, autoYes) {
			if err := adapter.InstallPackages([]string{missingZsh[0].PackageName}); err != nil {
				return err
			}
			ui.PrintOK("Zsh installed successfully!")
		} else {
			ui.PrintError("Zsh is required. Aborting setup.")
			os.Exit(1)
		}
	} else {
		ui.PrintOK("Zsh is installed.")
	}

	// 2. Check optional tools
	features := []string{"eza", "bat", "zoxide", "fd-find"}
	if !adapter.NeedsNvimUpdater() {
		// On Arch, neovim is available directly in the official 'extra' repo!
		features = append(features, "neovim")
	}

	missingPkgs, err := adapter.CheckMissingPackages(features)
	if err != nil {
		return err
	}

	if len(missingPkgs) > 0 {
		var pkgNames []string
		for _, m := range missingPkgs {
			ui.PrintWarn(fmt.Sprintf("Missing optional tool: %s (package: %s)", m.LogicalName, m.PackageName))
			pkgNames = append(pkgNames, m.PackageName)
		}

		if ui.AskYesNo(fmt.Sprintf("Would you like to install missing packages (%v)?", pkgNames), true, autoYes) {
			if err := adapter.InstallPackages(pkgNames); err != nil {
				ui.PrintError(fmt.Sprintf("Package installation failed: %v", err))
			} else {
				ui.PrintOK("Finished installing dependencies!")
			}
		}
	} else {
		ui.PrintOK("All optional dependencies are satisfied.")
	}

	// 3. Backup existing .zshrc
	if err := zsh.BackupZshrc(paths); err != nil {
		ui.PrintError(fmt.Sprintf("Backup failed: %v", err))
	}

	// 4. Ensure config & bin dirs exist
	if err := zsh.EnsureDirectories(paths); err != nil {
		return err
	}

	// 5. Gather feature choices
	var nvimUpdaterChoice bool
	if adapter.NeedsNvimUpdater() {
		nvimUpdaterChoice = ui.AskYesNo("Add Neovim updater (update-nvim) helper function?", true, autoYes)
	}

	if !autoYes {
		ui.PrintNerdFontTest()
	}
	nerdFontChoice := ui.AskYesNo("Enable Nerd Font icons for shell prompt and eza aliases?", true, autoYes)

	opts := zsh.ConfigOptions{
		NerdFonts:         nerdFontChoice,
		PromptDecorator:   ui.AskYesNo("Set Busy Shell as prompt decorator?", true, autoYes),
		Autocomplete:      ui.AskYesNo("Enable Busy Shell Autocomplete system?", true, autoYes),
		HistorySettings:   ui.AskYesNo("Allow Busy Shell to configure Zsh History settings?", true, autoYes),
		KeybindingsSearch: ui.AskYesNo("Enable arrow-key history prefix searching?", true, autoYes),
		ShellOptions:      ui.AskYesNo("Enable auto-cd and suppress beep alerts?", true, autoYes),
		EzaAliases:        ui.AskYesNo("Set up ls aliases using eza?", true, autoYes),
		BatAlias:          adapter.GetBatAlias(),
		Zoxide:            ui.AskYesNo("Enable zoxide directory jumper?", true, autoYes),
		NvimUpdater:       nvimUpdaterChoice,
		AddToPath:         ui.AskYesNo(fmt.Sprintf("Add %s to your shell PATH?", paths.BinDir), true, autoYes),
	}

	// 6. Handle fd symlink (e.g. Debian fdfind -> fd quirk)
	if adapter.NeedsFdSymlink() {
		if ui.AskYesNo("Fix Debian fd-find quirk (creating symlink to fd)?", true, autoYes) {
			targetBinary := adapter.GetFdTargetBinary()
			symlinkPath := filepath.Join(paths.BinDir, "fd")
			if _, err := os.Stat(targetBinary); err == nil {
				ui.PrintStep(fmt.Sprintf("Creating symlink %s -> %s...", symlinkPath, targetBinary))
				_ = os.Remove(symlinkPath) // Clear previous link if present
				if err := os.Symlink(targetBinary, symlinkPath); err == nil {
					ui.PrintOK("Symlink created.")
				} else {
					ui.PrintWarn(fmt.Sprintf("Could not create symlink: %v", err))
				}
			}
		}
	}

	// 7. Install self binary to ~/.local/bin/bzsh
	selfExecutable, err := os.Executable()
	if err == nil {
		destBinary := filepath.Join(paths.BinDir, "bzsh")
		if selfExecutable != destBinary {
			ui.PrintStep(fmt.Sprintf("Copying bzsh binary to %s...", destBinary))
			data, err := os.ReadFile(selfExecutable)
			if err == nil {
				_ = os.WriteFile(destBinary, data, 0755)
				ui.PrintOK("Installed bzsh binary.")
			}
		}
	}

	// 8. Generate and write config block
	block, err := zsh.GenerateConfigBlock(opts, paths)
	if err != nil {
		return fmt.Errorf("failed to generate config block: %w", err)
	}

	if err := zsh.UpdateZshrc(paths, block); err != nil {
		return err
	}

	ui.PrintOK("Setup completed successfully!")
	ui.PrintShellChangePrompt()
	ui.PrintFooter()
	return nil
}

// RunUpdate refreshes embedded scripts and updates the bzsh binary.
func RunUpdate() error {
	ui.PrintTitle()
	ui.PrintStep("Starting update check...")

	paths, err := zsh.GetDefaultPaths()
	if err != nil {
		return err
	}

	// Refresh embedded scripts in ~/.config/bzsh/
	ui.PrintStep(fmt.Sprintf("Refreshing script files in %s...", paths.ConfigDir))
	_ = zsh.WriteEmbeddedScript("prompt.bzsh", paths.ConfigDir)
	_ = zsh.WriteEmbeddedScript("autocomplete.bzsh", paths.ConfigDir)
	_ = zsh.WriteEmbeddedScript("nvim-update.bzsh", paths.ConfigDir)

	// Copy current executable to ~/.local/bin/bzsh if running from elsewhere
	selfExecutable, err := os.Executable()
	if err == nil {
		destBinary := filepath.Join(paths.BinDir, "bzsh")
		if selfExecutable != destBinary {
			data, err := os.ReadFile(selfExecutable)
			if err == nil {
				_ = os.WriteFile(destBinary, data, 0755)
			}
		}
	}

	ui.PrintOK("Updated successfully!")
	ui.PrintFooter()
	return nil
}

// RunUninstall cleans up bzsh configuration and files.
func RunUninstall(autoYes bool) error {
	ui.PrintTitle()
	ui.PrintStep("Starting uninstallation...")

	paths, err := zsh.GetDefaultPaths()
	if err != nil {
		return err
	}

	// Remove .zshrc block
	if err := zsh.RemoveZshrcBlock(paths); err != nil {
		ui.PrintError(fmt.Sprintf("Error removing zshrc block: %v", err))
	}

	// Ask to remove ~/.config/bzsh
	if ui.AskYesNo(fmt.Sprintf("Remove configuration folder %s for zero footprint?", paths.ConfigDir), true, autoYes) {
		ui.PrintStep(fmt.Sprintf("Deleting %s...", paths.ConfigDir))
		_ = os.RemoveAll(paths.ConfigDir)
		ui.PrintOK("Configuration folder deleted.")
	}

	// Ask to remove binary from ~/.local/bin/bzsh
	destBinary := filepath.Join(paths.BinDir, "bzsh")
	if ui.AskYesNo(fmt.Sprintf("Remove bzsh binary from %s?", paths.BinDir), true, autoYes) {
		ui.PrintStep(fmt.Sprintf("Deleting %s...", destBinary))
		_ = os.Remove(destBinary)
		ui.PrintOK("Binary deleted.")
	}

	// Clean up fd symlink if it exists
	symlinkPath := filepath.Join(paths.BinDir, "fd")
	if fi, err := os.Lstat(symlinkPath); err == nil && (fi.Mode()&os.ModeSymlink != 0) {
		ui.PrintStep(fmt.Sprintf("Removing fd-find symlink at %s...", symlinkPath))
		_ = os.Remove(symlinkPath)
		ui.PrintOK("Symlink deleted.")
	}

	ui.PrintOK("Uninstall completed!")
	ui.PrintFooter()
	return nil
}

// RunFontInstall downloads and installs SymbolsNerdFontMono-Regular.ttf into ~/.local/share/fonts/
func RunFontInstall() error {
	ui.PrintTitle()
	ui.PrintStep("Starting Nerd Font download and installation...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	fontsDir := filepath.Join(homeDir, ".local", "share", "fonts")
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		return fmt.Errorf("failed to create fonts directory %s: %w", fontsDir, err)
	}

	fontURL := "https://github.com/ryanoasis/nerd-fonts/raw/main/patched-fonts/NerdFontsSymbolsOnly/SymbolsNerdFontMono-Regular.ttf"
	targetFontPath := filepath.Join(fontsDir, "SymbolsNerdFontMono-Regular.ttf")

	ui.PrintStep(fmt.Sprintf("Downloading SymbolsNerdFontMono-Regular.ttf to %s...", fontsDir))

	if err := downloadFontFile(fontURL, targetFontPath); err != nil {
		return fmt.Errorf("failed downloading font: %w", err)
	}

	ui.PrintOK("Font downloaded successfully!")

	// Refresh font cache if fc-cache is present
	if _, err := exec.LookPath("fc-cache"); err == nil {
		ui.PrintStep("Updating desktop font cache with fc-cache -fv...")
		cmd := exec.Command("fc-cache", "-fv")
		_ = cmd.Run()
		ui.PrintOK("Font cache refreshed.")
	}

	ui.PrintOK("Nerd Font symbols installed to ~/.local/share/fonts/!")
	fmt.Println()
	ui.PrintWarn("Note: Terminal emulators (Alacritty, Kitty, VS Code, iTerm2, etc.) require selecting")
	fmt.Println("      the Nerd Font or font fallback in their app settings for icons to display!")
	ui.PrintFooter()
	return nil
}

func downloadFontFile(url string, dest string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error status %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

