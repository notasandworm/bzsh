package distro

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/notasandworm/bzsh/internal/ui"
)

// DebianAdapter handles Debian & Ubuntu derivatives.
// Matt's note: Ah, Debian! The classic. But why did Debian rename 'bat' to 'batcat' and 'fd' to 'fdfind'?!
// Don't worry, our adapter handles those quirks seamlessly.
type DebianAdapter struct{}

func NewDebianAdapter() *DebianAdapter {
	return &DebianAdapter{}
}

func (d *DebianAdapter) Name() string {
	return "Debian / Ubuntu Derivative"
}

// CheckMissingPackages inspects which optional tools are missing on Debian systems.
func (d *DebianAdapter) CheckMissingPackages(features []string) ([]FeaturePackage, error) {
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
			if !CommandExists("bat") && !CommandExists("batcat") {
				missing = append(missing, FeaturePackage{LogicalName: "bat", PackageName: "bat", BinaryName: "batcat"})
			}
		case "fd-find":
			if !CommandExists("fd") && !CommandExists("fdfind") {
				missing = append(missing, FeaturePackage{LogicalName: "fd-find", PackageName: "fd-find", BinaryName: "fdfind"})
			}
		case "zoxide":
			if !CommandExists("zoxide") {
				missing = append(missing, FeaturePackage{LogicalName: "zoxide", PackageName: "zoxide", BinaryName: "zoxide"})
			}
		}
	}

	return missing, nil
}

// InstallPackages invokes apt-get to install missing tools.
func (d *DebianAdapter) InstallPackages(pkgs []string) error {
	if len(pkgs) == 0 {
		return nil
	}

	ui.PrintStep(fmt.Sprintf("Running apt-get update & install for: %v...", pkgs))

	cmdStr := fmt.Sprintf("sudo apt-get update && sudo apt-get install -y %s", joinArgs(pkgs))

	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("apt-get installation failed: %w", err)
	}

	return nil
}

// GetBatAlias returns Debian's specific batcat alias!
func (d *DebianAdapter) GetBatAlias() string {
	if CommandExists("batcat") {
		return "alias bat=\"batcat\""
	}
	return ""
}

// NeedsFdSymlink returns true for Debian because 'fdfind' needs a symlink to 'fd'.
func (d *DebianAdapter) NeedsFdSymlink() bool {
	return true
}

// GetFdTargetBinary finds where fdfind is installed on Debian.
func (d *DebianAdapter) GetFdTargetBinary() string {
	path, err := exec.LookPath("fdfind")
	if err == nil {
		return path
	}
	return "/usr/bin/fdfind"
}

// NeedsNvimUpdater returns true for Debian (since Debian repos often ship older Neovim versions).
func (d *DebianAdapter) NeedsNvimUpdater() bool {
	return true
}

func joinArgs(args []string) string {
	result := ""
	for i, a := range args {
		if i > 0 {
			result += " "
		}
		result += a
	}
	return result
}
