package zsh

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractUserContentStandardMarkers(t *testing.T) {
	input := `export PRE_ENV="1"
# >>> bzsh initialize >>>
# bzsh config block here
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

# >>> bzsh initialize >>>
# new bzsh block
# <<< bzsh initialize <<<

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

# >>> bzsh initialize >>>
# new bzsh block
# <<< bzsh initialize <<<

export BACKUP_POST="ok"
`
	if updatedStr != expected {
		t.Errorf("unexpected restored zshrc content.\nExpected:\n%s\nGot:\n%s", expected, updatedStr)
	}
}
