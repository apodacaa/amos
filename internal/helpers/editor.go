package helpers

import (
	"os"
	"strings"
)

// entryTemplate is the template for new entries
const entryTemplate = ``

// GetEditor returns the editor command to use.
// Checks $EDITOR, then $VISUAL, then falls back to "nano".
func GetEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if visual := os.Getenv("VISUAL"); visual != "" {
		return visual
	}
	return "nano"
}

// CreateTempFile creates a temporary file for editing an entry.
// If existingContent is empty, uses the entry template.
// Returns the path to the temp file.
func CreateTempFile(existingContent string) (string, error) {
	// Create temp file with .md extension for editor syntax highlighting
	tmpFile, err := os.CreateTemp("", "amos-entry-*.md")
	if err != nil {
		return "", err
	}

	content := existingContent
	if content == "" {
		content = entryTemplate
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", err
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// ParseEntryFile reads a temp file and returns the content with template comments stripped.
// Returns (content, isEmpty, error) where isEmpty is true if the file is empty after stripping.
func ParseEntryFile(path string) (string, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false, err
	}

	content := string(data)

	// Strip HTML comment lines (template instructions)
	lines := strings.Split(content, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip HTML comment lines
		if strings.HasPrefix(trimmed, "<!--") && strings.HasSuffix(trimmed, "-->") {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	// Join lines back
	cleaned := strings.Join(cleanLines, "\n")

	// Trim leading/trailing whitespace
	cleaned = strings.TrimSpace(cleaned)

	isEmpty := cleaned == ""

	return cleaned, isEmpty, nil
}
