package zsh

import (
	"embed"
	"fmt"
)

// Hey! This embeds all of our custom Zsh scripts directly into the compiled Go binary.
// No extra downloads or HTTP requests needed at runtime. Super clean!
//
//go:embed embedded/*.bzsh
var embeddedFiles embed.FS

// GetEmbeddedScript fetches a zsh script content from our embedded filesystem.
func GetEmbeddedScript(filename string) ([]byte, error) {
	path := fmt.Sprintf("embedded/%s", filename)
	data, err := embeddedFiles.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("oops, couldn't load embedded script '%s': %w", filename, err)
	}
	return data, nil
}
