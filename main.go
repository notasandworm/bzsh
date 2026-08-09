package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/notasandworm/bzsh/internal/installer"
	"github.com/notasandworm/bzsh/internal/ui"
)

// 󰅶 Busy Shell (bzsh) - Go Edition!
//
// Hey! Matt here. Welcome to the Go binary rewrite of bzsh.
// Built to make installation blazing fast, zero-dependency (embedded zsh configs!),
// and distro-aware (Debian & Arch).
//
// Usage:
//   bzsh setup        # Non-interactive default setup ("Yes to all")
//   bzsh setup -i     # Interactive setup mode (toggle features manually)
//   bzsh update       # Update configuration and binaries
//   bzsh uninstall    # Remove bzsh configuration and binary cleanly

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "setup":
		setupFlags := flag.NewFlagSet("setup", flag.ExitOnError)
		interactive := setupFlags.Bool("i", false, "Run in interactive mode to toggle individual features")
		interactiveLong := setupFlags.Bool("interactive", false, "Run in interactive mode to toggle individual features")
		_ = setupFlags.Bool("y", true, "Run non-interactively (default behavior)")
		_ = setupFlags.Bool("yes", true, "Run non-interactively (default behavior)")

		_ = setupFlags.Parse(os.Args[2:])
		isInteractive := *interactive || *interactiveLong

		if err := installer.RunSetup(isInteractive); err != nil {
			ui.PrintError(fmt.Sprintf("Setup failed: %v", err))
			os.Exit(1)
		}

	case "update":
		if err := installer.RunUpdate(); err != nil {
			ui.PrintError(fmt.Sprintf("Update failed: %v", err))
			os.Exit(1)
		}

	case "uninstall":
		uninstallFlags := flag.NewFlagSet("uninstall", flag.ExitOnError)
		autoYes := uninstallFlags.Bool("y", false, "Auto-approve uninstallation prompts")
		autoYesLong := uninstallFlags.Bool("yes", false, "Auto-approve uninstallation prompts")

		_ = uninstallFlags.Parse(os.Args[2:])
		isAutoYes := *autoYes || *autoYesLong

		if err := installer.RunUninstall(isAutoYes); err != nil {
			ui.PrintError(fmt.Sprintf("Uninstall failed: %v", err))
			os.Exit(1)
		}

	case "font", "install-font":
		if err := installer.RunFontInstall(); err != nil {
			ui.PrintError(fmt.Sprintf("Font installation failed: %v", err))
			os.Exit(1)
		}

	case "version", "-v", "--version", "-version":
		ui.PrintVersion()

	case "help", "-h", "--help":
		showHelp()

	default:
		ui.PrintError(fmt.Sprintf("Unknown command: '%s'", command))
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	ui.PrintTitle()
	fmt.Println("Usage: bzsh <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  setup        Run the installation process (defaults to auto-yes)")
	fmt.Println("  update       Refresh script files and update local binary")
	fmt.Println("  font         Download and install SymbolsNerdFontMono to ~/.local/share/fonts/")
	fmt.Println("  uninstall    Remove bzsh settings and configuration folders")
	fmt.Println("  version      Show bzsh binary version (-v, --version)")
	fmt.Println("  help         Show this help message")
	fmt.Println()
	fmt.Println("Options for setup:")
	fmt.Println("  -i, --interactive  Run interactively to toggle features manually")
	fmt.Println("  -y, --yes          Run non-interactively (default)")
	fmt.Println()
	fmt.Println("Options for uninstall:")
	fmt.Println("  -y, --yes          Auto-approve prompt choices")
	ui.PrintFooter()
}
