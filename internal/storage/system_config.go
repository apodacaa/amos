package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SystemConfig stores system-level settings (like data directory path)
// Stored in ~/.config/amos/settings.json
type SystemConfig struct {
	DataDir string `json:"data_dir"` // Custom data directory path
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
func LoadSystemConfig() (SystemConfig, error) {
	path, err := GetSystemConfigPath()
	if err != nil {
		return SystemConfig{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// No system config yet - return empty (will use default)
			return SystemConfig{}, nil
		}
		return SystemConfig{}, err
	}

	var config SystemConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return SystemConfig{}, err
	}

	return config, nil
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
