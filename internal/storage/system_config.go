package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SystemConfig stores all application settings
// Stored in ~/.config/amos/settings.json
type SystemConfig struct {
	DataDir            string `json:"data_dir,omitempty"`             // Custom data directory path
	Theme              string `json:"theme,omitempty"`                // Theme name (e.g., "cyberpunk", "brutalist")
	HideCompletedTodos bool   `json:"hide_completed_todos,omitempty"` // Hide done todos in todo list
}

// DefaultSystemConfig returns the default configuration
func DefaultSystemConfig() SystemConfig {
	return SystemConfig{
		Theme: "cyberpunk", // Default to cyberpunk theme
	}
}

// GetSystemConfigPath returns ~/.config/amos/settings.json
func GetSystemConfigPath() (string, error) {
	configDir, err := os.UserConfigDir() // Returns ~/.config on Unix
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "amos", "settings.json"), nil
}

// LoadSystemConfig loads system config from ~/.config/amos/settings.json
// If theme is not set, tries to migrate from old config.json, otherwise uses default
func LoadSystemConfig() (SystemConfig, error) {
	path, err := GetSystemConfigPath()
	if err != nil {
		return DefaultSystemConfig(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// No system config yet - try to migrate theme from old config.json
			config := DefaultSystemConfig()
			migratedTheme := migrateThemeFromOldConfig()
			if migratedTheme != "" {
				config.Theme = migratedTheme
			}
			return config, nil
		}
		return DefaultSystemConfig(), err
	}

	var config SystemConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultSystemConfig(), err
	}

	// If theme is empty, try migration or use default
	if config.Theme == "" {
		migratedTheme := migrateThemeFromOldConfig()
		if migratedTheme != "" {
			config.Theme = migratedTheme
		} else {
			config.Theme = DefaultSystemConfig().Theme
		}
	}

	return config, nil
}

// migrateThemeFromOldConfig attempts to read theme from old-style <data-dir>/config.json
func migrateThemeFromOldConfig() string {
	// Get data directory (use default .amos if not configured)
	dir, err := GetAmosDir()
	if err != nil {
		return ""
	}

	oldConfigPath := filepath.Join(dir, "config.json")
	data, err := os.ReadFile(oldConfigPath)
	if err != nil {
		return "" // Old config doesn't exist or can't be read
	}

	var oldConfig struct {
		Theme string `json:"theme"`
	}
	if err := json.Unmarshal(data, &oldConfig); err != nil {
		return ""
	}

	return oldConfig.Theme
}

// SaveSystemConfig saves system config to ~/.config/amos/settings.json
func SaveSystemConfig(config SystemConfig) error {
	path, err := GetSystemConfigPath()
	if err != nil {
		return err
	}

	// Ensure ~/.config/amos directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
