package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.StaleDays != 30 {
		t.Errorf("expected StaleDays=30, got %d", config.StaleDays)
	}

	if len(config.Protected) == 0 {
		t.Error("expected non-empty Protected slice")
	}
}

func TestLoadConfig_NoFile(t *testing.T) {
	// Should return default config when file doesn't exist
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.StaleDays != 30 {
		t.Errorf("expected default StaleDays=30, got %d", config.StaleDays)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Save a config to a temporary location
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	testConfig := &Config{
		StaleDays: 60,
		Protected: []string{"main", "develop", "staging"},
	}

	if err := SaveConfig(testConfig); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Load it back
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.StaleDays != 60 {
		t.Errorf("expected StaleDays=60, got %d", loaded.StaleDays)
	}

	if len(loaded.Protected) != 3 {
		t.Errorf("expected 3 protected patterns, got %d", len(loaded.Protected))
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create invalid YAML file
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	configPath := filepath.Join(tmpHome, ".branch-clean.yaml")
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("failed to write invalid config: %v", err)
	}

	_, err := LoadConfig()
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadConfig_MissingValues(t *testing.T) {
	// Create config with missing values
	origHome := os.Getenv("HOME")
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	configPath := filepath.Join(tmpHome, ".branch-clean.yaml")
	// Empty config file
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Should apply defaults
	if config.StaleDays != 30 {
		t.Errorf("expected default StaleDays=30, got %d", config.StaleDays)
	}

	if len(config.Protected) == 0 {
		t.Error("expected default protected patterns")
	}
}
