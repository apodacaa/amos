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

// TestSaveAndLoadTodos tests saving and loading todos
func TestSaveAndLoadTodos(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	entryID := "test-entry-1"
	testTodos := []models.Todo{
		{
			ID:        "todo-1",
			Title:     "Buy groceries",
			Status:    "open",
			Tags:      []string{"personal", "shopping"},
			CreatedAt: time.Now(),
			EntryID:   &entryID,
		},
		{
			ID:        "todo-2",
			Title:     "Write tests",
			Status:    "done",
			Tags:      []string{"work"},
			CreatedAt: time.Now(),
			EntryID:   nil, // Standalone todo
		},
	}

	err := SaveTodos(testTodos)
	if err != nil {
		t.Fatalf("SaveTodos() failed: %v", err)
	}

	// Test loading todos
	loaded, err := LoadTodos()
	if err != nil {
		t.Fatalf("LoadTodos() failed: %v", err)
	}

	if len(loaded) != len(testTodos) {
		t.Errorf("LoadTodos() returned %d todos, want %d", len(loaded), len(testTodos))
	}

	// Verify first todo
	if loaded[0].ID != testTodos[0].ID {
		t.Errorf("LoadTodos() ID = %v, want %v", loaded[0].ID, testTodos[0].ID)
	}
	if loaded[0].Title != testTodos[0].Title {
		t.Errorf("LoadTodos() Title = %v, want %v", loaded[0].Title, testTodos[0].Title)
	}
	if loaded[0].Status != testTodos[0].Status {
		t.Errorf("LoadTodos() Status = %v, want %v", loaded[0].Status, testTodos[0].Status)
	}

	// Verify second todo (standalone - nil EntryID)
	if loaded[1].EntryID != nil {
		t.Errorf("LoadTodos() EntryID = %v, want nil", loaded[1].EntryID)
	}
}

// TestLoadTodosEmptyFile tests loading when todos.json doesn't exist
func TestLoadTodosEmptyFile(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Load todos when file doesn't exist
	todos, err := LoadTodos()
	if err != nil {
		t.Fatalf("LoadTodos() failed: %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("LoadTodos() returned %d todos, want 0 for non-existent file", len(todos))
	}
}

// TestSaveTodoNew tests saving a new todo
func TestSaveTodoNew(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	newTodo := models.Todo{
		ID:        "todo-new",
		Title:     "New task",
		Status:    "open",
		Tags:      []string{"test"},
		CreatedAt: time.Now(),
		EntryID:   nil,
	}

	err := SaveTodo(newTodo)
	if err != nil {
		t.Fatalf("SaveTodo() failed: %v", err)
	}

	// Load and verify
	todos, err := LoadTodos()
	if err != nil {
		t.Fatalf("LoadTodos() failed: %v", err)
	}

	if len(todos) != 1 {
		t.Fatalf("LoadTodos() returned %d todos, want 1", len(todos))
	}

	if todos[0].ID != newTodo.ID {
		t.Errorf("SaveTodo() ID = %v, want %v", todos[0].ID, newTodo.ID)
	}
}

// TestSaveTodoUpdate tests updating an existing todo
func TestSaveTodoUpdate(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Save initial todo
	initialTodo := models.Todo{
		ID:        "todo-update",
		Title:     "Original title",
		Status:    "open",
		Tags:      []string{"tag1"},
		CreatedAt: time.Now(),
		EntryID:   nil,
	}

	err := SaveTodo(initialTodo)
	if err != nil {
		t.Fatalf("SaveTodo() initial save failed: %v", err)
	}

	// Update the todo
	updatedTodo := initialTodo
	updatedTodo.Title = "Updated title"
	updatedTodo.Status = "done"
	updatedTodo.Tags = []string{"tag1", "tag2"}

	err = SaveTodo(updatedTodo)
	if err != nil {
		t.Fatalf("SaveTodo() update failed: %v", err)
	}

	// Load and verify only one todo exists with updated values
	todos, err := LoadTodos()
	if err != nil {
		t.Fatalf("LoadTodos() failed: %v", err)
	}

	if len(todos) != 1 {
		t.Errorf("LoadTodos() returned %d todos, want 1 after update", len(todos))
	}

	if todos[0].Title != "Updated title" {
		t.Errorf("SaveTodo() updated Title = %v, want 'Updated title'", todos[0].Title)
	}

	if todos[0].Status != "done" {
		t.Errorf("SaveTodo() updated Status = %v, want 'done'", todos[0].Status)
	}
}
