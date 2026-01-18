package helpers

import (
	"os"
	"strings"
	"testing"
)

func TestGetEditor(t *testing.T) {
	tests := []struct {
		name     string
		editor   string
		visual   string
		expected string
	}{
		{
			name:     "EDITOR set",
			editor:   "vim",
			visual:   "",
			expected: "vim",
		},
		{
			name:     "VISUAL set, no EDITOR",
			editor:   "",
			visual:   "code",
			expected: "code",
		},
		{
			name:     "EDITOR takes precedence over VISUAL",
			editor:   "nvim",
			visual:   "code",
			expected: "nvim",
		},
		{
			name:     "fallback to nano",
			editor:   "",
			visual:   "",
			expected: "nano",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origEditor := os.Getenv("EDITOR")
			origVisual := os.Getenv("VISUAL")
			defer func() {
				os.Setenv("EDITOR", origEditor)
				os.Setenv("VISUAL", origVisual)
			}()

			// Set test values
			os.Setenv("EDITOR", tt.editor)
			os.Setenv("VISUAL", tt.visual)

			result := GetEditor()
			if result != tt.expected {
				t.Errorf("GetEditor() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCreateTempFile(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		expectEmpty     bool
		expectInContent string
	}{
		{
			name:        "empty content creates empty file",
			content:     "",
			expectEmpty: true,
		},
		{
			name:            "existing content is preserved",
			content:         "My Entry Title\n\nSome body text",
			expectEmpty:     false,
			expectInContent: "My Entry Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := CreateTempFile(tt.content)
			if err != nil {
				t.Fatalf("CreateTempFile() error = %v", err)
			}
			defer os.Remove(path)

			// Verify file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Error("CreateTempFile() did not create file")
			}

			// Verify file has .md extension
			if !strings.HasSuffix(path, ".md") {
				t.Error("CreateTempFile() should create .md file")
			}

			// Read content
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}
			content := string(data)

			if tt.expectEmpty {
				if content != "" {
					t.Errorf("Expected empty file, got %q", content)
				}
			}

			if tt.expectInContent != "" {
				if !strings.Contains(content, tt.expectInContent) {
					t.Errorf("Expected %q in file content", tt.expectInContent)
				}
			}
		})
	}
}

func TestParseEntryFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		wantContent string
		wantEmpty   bool
	}{
		{
			name: "strips template comments",
			fileContent: `<!-- First line becomes the title -->
<!-- Use @tags for organization -->
My Entry Title

Some body text`,
			wantContent: "My Entry Title\n\nSome body text",
			wantEmpty:   false,
		},
		{
			name: "empty after stripping comments",
			fileContent: `<!-- First line becomes the title -->
<!-- Use @tags for organization -->

`,
			wantContent: "",
			wantEmpty:   true,
		},
		{
			name:        "content without comments",
			fileContent: "Title\n\nBody with @tag",
			wantContent: "Title\n\nBody with @tag",
			wantEmpty:   false,
		},
		{
			name: "preserves inline content",
			fileContent: `<!-- comment -->
Title
<!-- another comment -->

Body text
`,
			wantContent: "Title\n\nBody text",
			wantEmpty:   false,
		},
		{
			name:        "completely empty file",
			fileContent: "",
			wantContent: "",
			wantEmpty:   true,
		},
		{
			name:        "whitespace only",
			fileContent: "   \n\n\t  ",
			wantContent: "",
			wantEmpty:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file with test content
			tmpFile, err := os.CreateTemp("", "test-*.md")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.fileContent); err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}
			tmpFile.Close()

			// Test ParseEntryFile
			content, isEmpty, err := ParseEntryFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("ParseEntryFile() error = %v", err)
			}

			if content != tt.wantContent {
				t.Errorf("ParseEntryFile() content = %q, want %q", content, tt.wantContent)
			}

			if isEmpty != tt.wantEmpty {
				t.Errorf("ParseEntryFile() isEmpty = %v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}

func TestParseEntryFile_NonExistent(t *testing.T) {
	_, _, err := ParseEntryFile("/nonexistent/path/file.md")
	if err == nil {
		t.Error("ParseEntryFile() should return error for non-existent file")
	}
}
