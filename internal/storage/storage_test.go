package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apodacaa/amos/internal/models"
)

func TestSaveAndLoadEntries(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test saving entries
	testEntries := []models.Entry{
		{
			ID:        "test-1",
			Title:     "Test Entry 1",
			Body:      "Test body 1",
			Tags:      []string{"tag1", "tag2"},
			Timestamp: time.Now(),
		},
		{
			ID:        "test-2",
			Title:     "Test Entry 2",
			Body:      "Test body 2",
			Tags:      []string{"tag3"},
			Timestamp: time.Now(),
		},
	}

	err := SaveEntries(testEntries)
	if err != nil {
		t.Fatalf("SaveEntries() failed: %v", err)
	}

	// Test loading entries
	loaded, err := LoadEntries()
	if err != nil {
		t.Fatalf("LoadEntries() failed: %v", err)
	}

	if len(loaded) != len(testEntries) {
		t.Errorf("LoadEntries() returned %d entries, want %d", len(loaded), len(testEntries))
	}

	// Verify first entry
	if loaded[0].ID != testEntries[0].ID {
		t.Errorf("LoadEntries() ID = %v, want %v", loaded[0].ID, testEntries[0].ID)
	}
	if loaded[0].Title != testEntries[0].Title {
		t.Errorf("LoadEntries() Title = %v, want %v", loaded[0].Title, testEntries[0].Title)
	}
}

func TestLoadEntriesEmptyFile(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Load from non-existent file should return empty slice
	entries, err := LoadEntries()
	if err != nil {
		t.Fatalf("LoadEntries() failed: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("LoadEntries() returned %d entries, want 0", len(entries))
	}
}

func TestSaveEntryNew(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Save new entry
	entry := models.Entry{
		ID:        "new-entry",
		Title:     "New Entry",
		Body:      "New body",
		Tags:      []string{"new"},
		Timestamp: time.Now(),
	}

	err := SaveEntry(entry)
	if err != nil {
		t.Fatalf("SaveEntry() failed: %v", err)
	}

	// Load and verify
	entries, err := LoadEntries()
	if err != nil {
		t.Fatalf("LoadEntries() failed: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	if entries[0].ID != entry.ID {
		t.Errorf("Entry ID = %v, want %v", entries[0].ID, entry.ID)
	}
}

func TestSaveEntryUpdate(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Save initial entry
	entry := models.Entry{
		ID:        "update-test",
		Title:     "Original Title",
		Body:      "Original body",
		Tags:      []string{"original"},
		Timestamp: time.Now(),
	}

	err := SaveEntry(entry)
	if err != nil {
		t.Fatalf("SaveEntry() failed: %v", err)
	}

	// Update the entry
	entry.Title = "Updated Title"
	entry.Body = "Updated body"
	entry.Tags = []string{"updated"}

	err = SaveEntry(entry)
	if err != nil {
		t.Fatalf("SaveEntry() update failed: %v", err)
	}

	// Load and verify update
	entries, err := LoadEntries()
	if err != nil {
		t.Fatalf("LoadEntries() failed: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry after update, got %d", len(entries))
	}

	if entries[0].Title != "Updated Title" {
		t.Errorf("Entry Title = %v, want 'Updated Title'", entries[0].Title)
	}

	if entries[0].Body != "Updated body" {
		t.Errorf("Entry Body = %v, want 'Updated body'", entries[0].Body)
	}
}

func TestGetAmosDir(t *testing.T) {
	dir, err := GetAmosDir()
	if err != nil {
		t.Fatalf("GetAmosDir() failed: %v", err)
	}

	if dir == "" {
		t.Error("GetAmosDir() returned empty string")
	}

	if !filepath.IsAbs(dir) {
		t.Errorf("GetAmosDir() = %v, want absolute path", dir)
	}
}

func TestEnsureAmosDir(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	err := EnsureAmosDir()
	if err != nil {
		t.Fatalf("EnsureAmosDir() failed: %v", err)
	}

	// Verify directory was created
	dir, _ := GetAmosDir()
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}
}
