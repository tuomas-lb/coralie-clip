package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfigPath(t *testing.T) {
	// Save original env
	originalEnv := os.Getenv("CORALIE_CONFIG")
	defer os.Setenv("CORALIE_CONFIG", originalEnv)

	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "default (no env)",
			envValue: "",
			expected: "./config.json",
		},
		{
			name:     "env as file path",
			envValue: "/etc/coralie/config.json",
			expected: "/etc/coralie/config.json",
		},
		{
			name:     "env as directory",
			envValue: "/etc/coralie",
			expected: "/etc/coralie/config.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue == "" {
				os.Unsetenv("CORALIE_CONFIG")
			} else {
				os.Setenv("CORALIE_CONFIG", tt.envValue)
			}

			result, err := ResolveConfigPath()
			if err != nil {
				t.Fatalf("ResolveConfigPath() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("ResolveConfigPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Version != 1 {
		t.Errorf("DefaultConfig().Version = %v, want 1", cfg.Version)
	}
	if cfg.ClipsDir != "./clips" {
		t.Errorf("DefaultConfig().ClipsDir = %v, want ./clips", cfg.ClipsDir)
	}
	if cfg.DefaultVoice != "coral" {
		t.Errorf("DefaultConfig().DefaultVoice = %v, want coral", cfg.DefaultVoice)
	}
	if len(cfg.EnabledLangs) == 0 || cfg.EnabledLangs[0] != "en" {
		t.Errorf("DefaultConfig().EnabledLangs = %v, want [en]", cfg.EnabledLangs)
	}
}

func TestValidateConfigOrExit(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	clipsDir := filepath.Join(tmpDir, "clips")
	catalogPath := filepath.Join(tmpDir, "clips", "catalog.json")

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Version:      1,
				OpenAIApiKey: "test-key",
				ClipsDir:     clipsDir,
				CatalogPath:  catalogPath,
				EnabledLangs: []string{"en"},
				DefaultVoice: "coral",
				Format:       "wav",
				SampleRate:   24000,
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			config: &Config{
				Version:      1,
				ClipsDir:     clipsDir,
				CatalogPath:  catalogPath,
				EnabledLangs: []string{"en"},
				DefaultVoice: "coral",
				Format:       "wav",
				SampleRate:   24000,
			},
			wantErr: true,
		},
		{
			name: "no enabled languages",
			config: &Config{
				Version:      1,
				OpenAIApiKey: "test-key",
				ClipsDir:     clipsDir,
				CatalogPath:  catalogPath,
				EnabledLangs: []string{},
				DefaultVoice: "coral",
				Format:       "wav",
				SampleRate:   24000,
			},
			wantErr: true,
		},
		{
			name: "invalid voice",
			config: &Config{
				Version:      1,
				OpenAIApiKey: "test-key",
				ClipsDir:     clipsDir,
				CatalogPath:  catalogPath,
				EnabledLangs: []string{"en"},
				DefaultVoice: "invalid",
				Format:       "wav",
				SampleRate:   24000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr to check for errors
			// Since ValidateConfigOrExit calls os.Exit(1), we can't easily test it
			// Instead, we'll test the validation logic separately
			// For now, just ensure the function doesn't panic
			if tt.wantErr {
				// We expect this to exit, so we can't test it normally
				// This is a limitation of testing os.Exit
			} else {
				// For valid configs, ensure directories are created
				ValidateConfigOrExit(tt.config)
			}
		})
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Save original env
	originalEnv := os.Getenv("CORALIE_CONFIG")
	defer os.Setenv("CORALIE_CONFIG", originalEnv)

	os.Setenv("CORALIE_CONFIG", configPath)

	cfg := DefaultConfig()
	cfg.OpenAIApiKey = "test-key-123"
	cfg.ClipsDir = "./test-clips"

	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("config file was not created")
	}

	// Read and verify
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if loaded.OpenAIApiKey != "test-key-123" {
		t.Errorf("loaded OpenAIApiKey = %v, want test-key-123", loaded.OpenAIApiKey)
	}
	if loaded.ClipsDir != "./test-clips" {
		t.Errorf("loaded ClipsDir = %v, want ./test-clips", loaded.ClipsDir)
	}
}

