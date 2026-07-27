package distro

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/notasandworm/bzsh/internal/ui"
)

// ArchAdapter handles Arch Linux & derivatives (Manjaro, EndeavourOS, etc.).
// Matt's note: Arch keeps things simple and upstream! No 'batcat' or 'fdfind' weirdness here.
// Official core/extra packages are named 'bat' and 'fd' directly. Beautiful!
type ArchAdapter struct{}

func NewArchAdapter() *ArchAdapter {
	return &ArchAdapter{}
}

func (a *ArchAdapter) Name() string {
	return "Arch Linux"
}

// CheckMissingPackages inspects which optional tools are missing on Arch Linux.
func (a *ArchAdapter) CheckMissingPackages(features []string) ([]FeaturePackage, error) {
	var missing []FeaturePackage

	for _, feat := range features {
		switch feat {
		case "zsh":
			if !CommandExists("zsh") {
				missing = append(missing, FeaturePackage{LogicalName: "zsh", PackageName: "zsh", BinaryName: "zsh"})
			}
		case "eza":
			if !CommandExists("eza") {
				missing = append(missing, FeaturePackage{LogicalName: "eza", PackageName: "eza", BinaryName: "eza"})
			}
		case "bat":
			if !CommandExists("bat") {
				missing = append(missing, FeaturePackage{LogicalName: "bat", PackageName: "bat", BinaryName: "bat"})
			}
		case "fd-find":
			if !CommandExists("fd") {
				missing = append(missing, FeaturePackage{LogicalName: "fd-find", PackageName: "fd", BinaryName: "fd"})
			}
		case "zoxide":
			if !CommandExists("zoxide") {
				missing = append(missing, FeaturePackage{LogicalName: "zoxide", PackageName: "zoxide", BinaryName: "zoxide"})
			}
		}
	}

	return missing, nil
}

// InstallPackages invokes pacman to install missing packages from official repos.
func (a *ArchAdapter) InstallPackages(pkgs []string) error {
	if len(pkgs) == 0 {
		return nil
	}

	ui.PrintStep(fmt.Sprintf("Running pacman -S for official packages: %v...", pkgs))

	cmdStr := fmt.Sprintf("sudo pacman -S --needed --noconfirm %s", joinArgs(pkgs))

	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pacman installation failed: %w", err)
	}

	return nil
}

// GetBatAlias returns empty string on Arch because Arch uses native 'bat'!
func (a *ArchAdapter) GetBatAlias() string {
	return ""
}

// NeedsFdSymlink returns false on Arch because Arch installs native 'fd'!
func (a *ArchAdapter) NeedsFdSymlink() bool {
	return false
}

// GetFdTargetBinary is not needed on Arch, but satisfies the interface.
func (a *ArchAdapter) GetFdTargetBinary() string {
	path, err := exec.LookPath("fd")
	if err == nil {
		return path
	}
	return "/usr/bin/fd"
}
