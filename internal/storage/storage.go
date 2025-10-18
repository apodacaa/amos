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
