package distro

import (
	"os"
	"os/exec"
	"strings"
)

// FeaturePackage specifies distro-specific package names and binaries.
type FeaturePackage struct {
	LogicalName  string // e.g. "bat", "fd-find", "eza", "zoxide", "zsh"
	PackageName  string // distro package name for installation
	BinaryName   string // executable command name to check in PATH
}

// Adapter interface abstracts away all distro-specific package managers and quirks!
type Adapter interface {
	Name() string
	CheckMissingPackages(features []string) ([]FeaturePackage, error)
	InstallPackages(pkgs []string) error
	GetBatAlias() string
	NeedsFdSymlink() bool
	GetFdTargetBinary() string
}

// CommandExists checks if a binary exists in the user's PATH.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// DetectDistro figures out what Linux distro we're running on!
// Matt's note: We read /etc/os-release first, but if that fails, we check for apt-get or pacman binaries.
func DetectDistro() Adapter {
	osReleaseData, err := os.ReadFile("/etc/os-release")
	if err == nil {
		content := strings.ToLower(string(osReleaseData))
		if strings.Contains(content, "arch") || strings.Contains(content, "manjaro") || strings.Contains(content, "endeavouros") {
			return NewArchAdapter()
		}
		if strings.Contains(content, "debian") || strings.Contains(content, "ubuntu") || strings.Contains(content, "pop") || strings.Contains(content, "mint") {
			return NewDebianAdapter()
		}
	}

	// Fallback check based on package manager binaries
	if CommandExists("pacman") {
		return NewArchAdapter()
	}
	if CommandExists("apt-get") {
		return NewDebianAdapter()
	}

	// Default to Debian if we really can't tell (most common Linux target)
	return NewDebianAdapter()
}
