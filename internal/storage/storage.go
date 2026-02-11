package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/apodacaa/amos/internal/models"
)

const (
	entriesFile = "entries.json"
	todosFile   = "todos.json"
)

// Package-level variable for data directory (set during initialization)
var dataDir string = ".amos" // Default value

// SetDataDir configures the data directory for this session
// Must be called before any other storage operations
func SetDataDir(dir string) error {
	// Expand tilde if present
	expanded, err := expandPath(dir)
	if err != nil {
		return err
	}

	dataDir = expanded
	return nil
}

// expandPath expands ~ to home directory and converts user-specified relative paths to absolute
func expandPath(path string) (string, error) {
	// Handle tilde expansion
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		if path == "~" {
			return home, nil
		}

		// Replace ~/ with home/
		if strings.HasPrefix(path, "~/") {
			return filepath.Join(home, path[2:]), nil
		}
	}

	// If already absolute, return as-is
	if filepath.IsAbs(path) {
		return path, nil
	}

	// Special case: ".amos" is the default and should remain relative
	// so GetAmosDir() can join it with home directory
	if path == ".amos" {
		return path, nil
	}

	// Convert other relative paths to absolute (relative to current working directory)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// InitializeDataDir sets up the data directory based on env var or system config
// Call this once at startup before any storage operations
func InitializeDataDir() error {
	var dir string

	// 1. Check environment variable first (highest priority)
	if envDir := os.Getenv("AMOS_DATA_DIR"); envDir != "" {
		dir = envDir
	} else {
		// 2. Check system config
		sysConfig, err := LoadSystemConfig()
		if err != nil {
			// If system config fails to load, use default
			dir = ".amos"
		} else if len(sysConfig.Workspaces) > 0 && sysConfig.ActiveWorkspace != "" {
			// 3. Check active workspace
			for _, ws := range sysConfig.Workspaces {
				if ws.Name == sysConfig.ActiveWorkspace {
					dir = ws.DataDir
					break
				}
			}
			if dir == "" {
				// Active workspace not found in list, use default
				dir = ".amos"
			}
		} else if sysConfig.DataDir != "" {
			dir = sysConfig.DataDir
		} else {
			// System config exists but DataDir is empty - use default
			dir = ".amos"
		}
	}

	// Set the data directory (with tilde expansion)
	if err := SetDataDir(dir); err != nil {
		return err
	}

	// Ensure directory exists
	return EnsureAmosDir()
}

// SwitchWorkspace switches to a new data directory and ensures it exists
func SwitchWorkspace(dir string) error {
	if err := SetDataDir(dir); err != nil {
		return err
	}
	return EnsureAmosDir()
}

// GetAmosDir returns the path to the configured data directory (expanded)
func GetAmosDir() (string, error) {
	// If dataDir is already an absolute path, return it as-is
	if filepath.IsAbs(dataDir) {
		return dataDir, nil
	}

	// Otherwise, join with home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, dataDir), nil
}

// EnsureAmosDir creates ~/.amos directory if it doesn't exist
func EnsureAmosDir() error {
	dir, err := GetAmosDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

// LoadEntries loads all entries from entries.json
func LoadEntries() ([]models.Entry, error) {
	dir, err := GetAmosDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, entriesFile)

	// If file doesn't exist, return empty slice
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []models.Entry{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []models.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// SaveEntries saves all entries to entries.json
func SaveEntries(entries []models.Entry) error {
	if err := EnsureAmosDir(); err != nil {
		return err
	}

	dir, err := GetAmosDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, entriesFile)

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// SaveEntry saves or updates a single entry in the entries list
func SaveEntry(entry models.Entry) error {
	entries, err := LoadEntries()
	if err != nil {
		return err
	}

	// Check if entry exists (by ID) and update, or append new
	found := false
	for i, e := range entries {
		if e.ID == entry.ID {
			entries[i] = entry
			found = true
			break
		}
	}

	if !found {
		entries = append(entries, entry)
	}

	return SaveEntries(entries)
}

// LoadTodos loads all todos from todos.json
func LoadTodos() ([]models.Todo, error) {
	dir, err := GetAmosDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, todosFile)

	// If file doesn't exist, return empty slice
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []models.Todo{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var todos []models.Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}

	// Migration: ensure UpdatedAt is populated (set to CreatedAt if zero)
	for i := range todos {
		if todos[i].UpdatedAt.IsZero() {
			todos[i].UpdatedAt = todos[i].CreatedAt
		}
	}

	return todos, nil
}

// SaveTodos saves all todos to todos.json
func SaveTodos(todos []models.Todo) error {
	if err := EnsureAmosDir(); err != nil {
		return err
	}

	dir, err := GetAmosDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, todosFile)

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// SaveTodo saves or updates a single todo in the todos list
func SaveTodo(todo models.Todo) error {
	todos, err := LoadTodos()
	if err != nil {
		return err
	}

	// Check if todo exists (by ID) and update, or append new
	found := false
	for i, t := range todos {
		if t.ID == todo.ID {
			todos[i] = todo
			found = true
			break
		}
	}

	if !found {
		todos = append(todos, todo)
	}

	return SaveTodos(todos)
}

// DeleteEntry removes an entry by ID from entries.json
func DeleteEntry(id string) error {
	entries, err := LoadEntries()
	if err != nil {
		return err
	}

	// Filter out the entry with matching ID
	filtered := make([]models.Entry, 0, len(entries))
	for _, e := range entries {
		if e.ID != id {
			filtered = append(filtered, e)
		}
	}

	return SaveEntries(filtered)
}

// DeleteTodo removes a todo by ID from todos.json
func DeleteTodo(id string) error {
	todos, err := LoadTodos()
	if err != nil {
		return err
	}

	// Filter out the todo with matching ID
	filtered := make([]models.Todo, 0, len(todos))
	for _, t := range todos {
		if t.ID != id {
			filtered = append(filtered, t)
		}
	}

	return SaveTodos(filtered)
}

// DeleteEntryCascade removes an entry and all its linked todos
func DeleteEntryCascade(entryID string) error {
	// Load both entries and todos
	entries, err := LoadEntries()
	if err != nil {
		return err
	}

	todos, err := LoadTodos()
	if err != nil {
		return err
	}

	// Find the entry to get its linked todo IDs
	var todoIDsToDelete []string
	for _, e := range entries {
		if e.ID == entryID {
			todoIDsToDelete = e.TodoIDs
			break
		}
	}

	// Remove the entry
	filteredEntries := make([]models.Entry, 0, len(entries))
	for _, e := range entries {
		if e.ID != entryID {
			filteredEntries = append(filteredEntries, e)
		}
	}

	// Remove linked todos
	filteredTodos := make([]models.Todo, 0, len(todos))
	for _, t := range todos {
		// Keep todo if it's not in the delete list
		shouldDelete := false
		for _, deleteID := range todoIDsToDelete {
			if t.ID == deleteID {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			filteredTodos = append(filteredTodos, t)
		}
	}

	// Save both files
	if err := SaveEntries(filteredEntries); err != nil {
		return err
	}

	return SaveTodos(filteredTodos)
}
