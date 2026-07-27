package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ANSI Color codes for that sweet, sweet terminal aesthetic! 🎨
const (
	ColorReset   = "\033[0m"
	ColorBlue    = "\033[1;34m"
	ColorCyan    = "\033[1;36m"
	ColorGreen   = "\033[1;32m"
	ColorYellow  = "\033[1;33m"
	ColorRed     = "\033[1;31m"
	ColorMagenta = "\033[1;35m"
)

// PrintTitle outputs our signature Busy Shell (bzsh) header!
func PrintTitle() {
	fmt.Println()
	fmt.Printf("%s󰅶 Busy Shell (bzsh) CLI%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s=========================%s\n\n", ColorBlue, ColorReset)
}

// PrintStep displays a step in progress.
func PrintStep(msg string) {
	fmt.Printf("%s➜%s %s\n", ColorCyan, ColorReset, msg)
}

// PrintOK displays a successful step!
func PrintOK(msg string) {
	fmt.Printf("%s✔%s %s\n", ColorGreen, ColorReset, msg)
}

// PrintWarn displays a warning.
func PrintWarn(msg string) {
	fmt.Printf("%s⚠%s %s\n", ColorYellow, ColorReset, msg)
}

// PrintError displays an error message.
func PrintError(msg string) {
	fmt.Printf("%s✘%s %s\n", ColorRed, ColorReset, msg)
}

// PrintShellChangePrompt shows the user how to switch their default shell to Zsh using chsh / cat /etc/shells!
func PrintShellChangePrompt() {
	fmt.Println()
	fmt.Printf("%s➜ To make Zsh your default shell, check valid shells with:%s\n", ColorYellow, ColorReset)
	fmt.Printf("    cat /etc/shells   (or chsh -l)\n")
	fmt.Printf("%s  Then set Zsh as your default shell:%s\n", ColorYellow, ColorReset)
	fmt.Printf("    chsh -s $(which zsh)\n")
}

// PrintFooter shows our warm closing thank-you note.
func PrintFooter() {
	fmt.Println()
	fmt.Printf("%s♥ Thank you for using Busy Shell!%s\n", ColorGreen, ColorReset)
	fmt.Printf("%s✉ For queries, feedback, or issues, email notasandworm@gmail.com%s\n\n", ColorCyan, ColorReset)
}

// AskYesNo prompts the user for a Y/n question.
// If autoYes is true, it automatically returns true without blocking!
func AskYesNo(prompt string, defaultYes bool, autoYes bool) bool {
	// Auto-approve mode: fast, silent, no fuss!
	if autoYes {
		defaultStr := "Y"
		if !defaultYes {
			defaultStr = "N"
		}
		fmt.Printf("%s➜%s %s [%s] (auto-selected)\n", ColorCyan, ColorReset, prompt, defaultStr)
		return defaultYes
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		if defaultYes {
			fmt.Printf("%s?%s %s [Y/n]: ", ColorMagenta, ColorReset, prompt)
		} else {
			fmt.Printf("%s?%s %s [y/N]: ", ColorMagenta, ColorReset, prompt)
		}

		input, err := reader.ReadString('\n')
		if err != nil {
			return defaultYes
		}

		trimmed := strings.ToLower(strings.TrimSpace(input))
		if trimmed == "" {
			return defaultYes
		}
		if trimmed == "y" || trimmed == "yes" {
			return true
		}
		if trimmed == "n" || trimmed == "no" {
			return false
		}

		fmt.Println("Please answer 'y' or 'n'.")
	}
}
