package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSystemConfig_SaveAndLoad(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Override UserConfigDir for testing
	originalConfigDir := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if originalConfigDir == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", originalConfigDir)
		}
	}()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Test save
	config := SystemConfig{
		DataDir: "~/my-custom-dir",
	}

	err := SaveSystemConfig(config)
	if err != nil {
		t.Fatalf("SaveSystemConfig failed: %v", err)
	}

	// Verify file exists
	path, err := GetSystemConfigPath()
	if err != nil {
		t.Fatalf("GetSystemConfigPath failed: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", path)
	}

	// Test load
	loaded, err := LoadSystemConfig()
	if err != nil {
		t.Fatalf("LoadSystemConfig failed: %v", err)
	}

	if loaded.DataDir != config.DataDir {
		t.Errorf("Loaded DataDir = %q, want %q", loaded.DataDir, config.DataDir)
	}
}

func TestSystemConfig_LoadMissing(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Override UserConfigDir for testing
	originalConfigDir := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if originalConfigDir == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", originalConfigDir)
		}
	}()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Load from non-existent file should return empty config
	config, err := LoadSystemConfig()
	if err != nil {
		t.Fatalf("LoadSystemConfig failed: %v", err)
	}

	if config.DataDir != "" {
		t.Errorf("Expected empty DataDir, got %q", config.DataDir)
	}
}

func TestSystemConfig_DirectoryCreation(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Override UserConfigDir for testing
	originalConfigDir := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if originalConfigDir == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", originalConfigDir)
		}
	}()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Verify directory doesn't exist yet
	amosDir := filepath.Join(tmpDir, "amos")
	if _, err := os.Stat(amosDir); !os.IsNotExist(err) {
		t.Fatalf("Test setup failed: amos directory already exists")
	}

	// Save config should create directory
	config := SystemConfig{DataDir: "~/test"}
	err := SaveSystemConfig(config)
	if err != nil {
		t.Fatalf("SaveSystemConfig failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(amosDir); os.IsNotExist(err) {
		t.Fatalf("SaveSystemConfig did not create directory")
	}
}
