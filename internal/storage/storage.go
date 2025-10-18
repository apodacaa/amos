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
