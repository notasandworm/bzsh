package zsh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractUserContentStandardMarkers(t *testing.T) {
	input := `export PRE_ENV="1"
# WARN: >>> bzsh initialize >>>
# bzsh config block here
# WARN: <<< bzsh initialize <<<
alias post_alias="echo hello"
`
	before, after := ExtractUserContent(input)
	if before != "export PRE_ENV=\"1\"\n" {
		t.Errorf("expected before to be %q, got %q", "export PRE_ENV=\"1\"\n", before)
	}
	if after != "\nalias post_alias=\"echo hello\"\n" {
		t.Errorf("expected after to be %q, got %q", "\nalias post_alias=\"echo hello\"\n", after)
	}
}

func TestExtractUserContentLegacyBzshMarkers(t *testing.T) {
	input := `export PRE_ENV="1"
# >>> bzsh initialize >>>
# old bzsh block
# <<< bzsh initialize <<<
alias post_alias="echo hello"
`
	before, after := ExtractUserContent(input)
	if before != "export PRE_ENV=\"1\"\n" {
		t.Errorf("expected before to be %q, got %q", "export PRE_ENV=\"1\"\n", before)
	}
	if after != "\nalias post_alias=\"echo hello\"\n" {
		t.Errorf("expected after to be %q, got %q", "\nalias post_alias=\"echo hello\"\n", after)
	}
}

func TestExtractUserContentLegacyMarkers(t *testing.T) {
	input := `export LEGACY_PRE="1"
# >>> bzhrc initialize >>>
# legacy block
# <<< bzhrc initialize <<<
export LEGACY_POST="2"
`
	before, after := ExtractUserContent(input)
	if before != "export LEGACY_PRE=\"1\"\n" {
		t.Errorf("expected before to be %q, got %q", "export LEGACY_PRE=\"1\"\n", before)
	}
	if after != "\nexport LEGACY_POST=\"2\"\n" {
		t.Errorf("expected after to be %q, got %q", "\nexport LEGACY_POST=\"2\"\n", after)
	}
}

func TestUpdateZshrcPreservesUserContent(t *testing.T) {
	tmpDir := t.TempDir()
	zshrcPath := filepath.Join(tmpDir, ".zshrc")

	initialContent := `export MY_VAR="true"
# >>> bzsh initialize >>>
old block
# <<< bzsh initialize <<<
alias my_alias="ls -la"
`
	if err := os.WriteFile(zshrcPath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("failed to write initial zshrc: %v", err)
	}

	paths := &ConfigPaths{
		ZshrcFile: zshrcPath,
		ConfigDir: filepath.Join(tmpDir, ".config", "bzsh"),
		BinDir:    filepath.Join(tmpDir, ".local", "bin"),
	}

	newBlock := StartMarker + "\n# new bzsh block\n" + EndMarker

	if err := UpdateZshrc(paths, newBlock); err != nil {
		t.Fatalf("UpdateZshrc failed: %v", err)
	}

	updatedData, err := os.ReadFile(zshrcPath)
	if err != nil {
		t.Fatalf("failed to read updated zshrc: %v", err)
	}
	updatedStr := string(updatedData)

	expected := `export MY_VAR="true"

# WARN: >>> bzsh initialize >>>
# new bzsh block
# WARN: <<< bzsh initialize <<<

alias my_alias="ls -la"
`
	if updatedStr != expected {
		t.Errorf("unexpected updated zshrc content.\nExpected:\n%s\nGot:\n%s", expected, updatedStr)
	}
}

func TestUpdateZshrcRestoresFromBackup(t *testing.T) {
	tmpDir := t.TempDir()
	zshrcPath := filepath.Join(tmpDir, ".zshrc")
	backupPath := filepath.Join(tmpDir, ".zshrc.bzsh-backup.20260806030000")

	backupContent := `export BACKUP_PRE="ok"
# >>> bzsh initialize >>>
old bzsh block
# <<< bzsh initialize <<<
export BACKUP_POST="ok"
`
	if err := os.WriteFile(backupPath, []byte(backupContent), 0644); err != nil {
		t.Fatalf("failed to write backup file: %v", err)
	}

	// Active .zshrc does not exist
	paths := &ConfigPaths{
		ZshrcFile: zshrcPath,
		ConfigDir: filepath.Join(tmpDir, ".config", "bzsh"),
		BinDir:    filepath.Join(tmpDir, ".local", "bin"),
	}

	newBlock := StartMarker + "\n# new bzsh block\n" + EndMarker

	if err := UpdateZshrc(paths, newBlock); err != nil {
		t.Fatalf("UpdateZshrc failed: %v", err)
	}

	updatedData, err := os.ReadFile(zshrcPath)
	if err != nil {
		t.Fatalf("failed to read updated zshrc: %v", err)
	}
	updatedStr := string(updatedData)

	expected := `export BACKUP_PRE="ok"

# WARN: >>> bzsh initialize >>>
# new bzsh block
# WARN: <<< bzsh initialize <<<

export BACKUP_POST="ok"
`
	if updatedStr != expected {
		t.Errorf("unexpected restored zshrc content.\nExpected:\n%s\nGot:\n%s", expected, updatedStr)
	}
}

func TestEmbeddedScriptsMatchRoot(t *testing.T) {
	scripts := []string{"prompt.bzsh", "autocomplete.bzsh", "nvim-update.bzsh"}

	for _, script := range scripts {
		embeddedData, err := GetEmbeddedScript(script)
		if err != nil {
			t.Fatalf("failed to get embedded script %s: %v", script, err)
		}

		rootPath := filepath.Join("..", "..", script)
		rootData, err := os.ReadFile(rootPath)
		if err != nil {
			// Skip if test is running outside source repo context
			if os.IsNotExist(err) {
				t.Logf("Root file %s not found; skipping sync check", rootPath)
				continue
			}
			t.Fatalf("failed reading root file %s: %v", rootPath, err)
		}

		if string(embeddedData) != string(rootData) {
			t.Errorf("Script desynchronization error: embedded/%s does not match root %s!\nIf you edited %s in the root folder, remember to sync internal/zsh/embedded/%s and run `go build`!", script, script, script, script)
		}
	}
}

func TestGenerateConfigBlockNerdFonts(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &ConfigPaths{
		ZshrcFile: filepath.Join(tmpDir, ".zshrc"),
		ConfigDir: filepath.Join(tmpDir, ".config", "bzsh"),
		BinDir:    filepath.Join(tmpDir, ".local", "bin"),
	}

	// Case 1: NerdFonts = true
	optsEnabled := ConfigOptions{
		NerdFonts:  true,
		EzaAliases: true,
	}
	blockEnabled, err := GenerateConfigBlock(optsEnabled, paths)
	if err != nil {
		t.Fatalf("GenerateConfigBlock failed: %v", err)
	}
	if !strings.Contains(blockEnabled, "export BZSH_NERD_FONTS=1") {
		t.Errorf("expected BZSH_NERD_FONTS=1 in config block")
	}
	if !strings.Contains(blockEnabled, "--icons=always") {
		t.Errorf("expected eza --icons=always when NerdFonts=true")
	}
	if !strings.Contains(blockEnabled, "alias lls=\"eza --long --icons=always --colour=always --sort=modified\"") {
		t.Errorf("expected lls alias with --long and --icons=always in config block")
	}

	// Case 2: NerdFonts = false
	optsDisabled := ConfigOptions{
		NerdFonts:  false,
		EzaAliases: true,
	}
	blockDisabled, err := GenerateConfigBlock(optsDisabled, paths)
	if err != nil {
		t.Fatalf("GenerateConfigBlock failed: %v", err)
	}
	if !strings.Contains(blockDisabled, "export BZSH_NERD_FONTS=0") {
		t.Errorf("expected BZSH_NERD_FONTS=0 in config block")
	}
	if !strings.Contains(blockDisabled, "--icons=never") {
		t.Errorf("expected eza --icons=never when NerdFonts=false")
	}
	if !strings.Contains(blockDisabled, "alias lls=\"eza --long --icons=never --colour=always --sort=modified\"") {
		t.Errorf("expected lls alias with --long and --icons=never in config block")
	}

}



