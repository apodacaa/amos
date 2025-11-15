package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/apodacaa/amos/internal/models"
)

const (
	amosDir     = ".amos"
	entriesFile = "entries.json"
	todosFile   = "todos.json"
	configFile  = "config.json"
)

// GetAmosDir returns the path to ~/.amos directory
func GetAmosDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, amosDir), nil
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

// LoadConfig loads user configuration from config.json
func LoadConfig() (models.Config, error) {
	dir, err := GetAmosDir()
	if err != nil {
		return models.DefaultConfig(), err
	}

	path := filepath.Join(dir, configFile)

	// If file doesn't exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return models.DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return models.DefaultConfig(), err
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return models.DefaultConfig(), err
	}

	return config, nil
}

// SaveConfig saves user configuration to config.json
func SaveConfig(config models.Config) error {
	if err := EnsureAmosDir(); err != nil {
		return err
	}

	dir, err := GetAmosDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, configFile)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
