// Package config handles configuration loading, validation, and management.
// It supports loading from config.json, .env files, and environment variables.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config represents the application configuration.
type Config struct {
	Version       int      `json:"version"`
	OpenAIApiKey  string   `json:"openaiApiKey,omitempty"`
	ClipsDir      string   `json:"clipsDir,omitempty"`
	CatalogPath   string   `json:"catalogPath,omitempty"`
	EnabledLangs  []string `json:"enabledLangs,omitempty"`
	DefaultVoice  string   `json:"defaultVoice,omitempty"`
	Format        string   `json:"format,omitempty"`
	SampleRate    int      `json:"sampleRate,omitempty"`
	BaseURL       string   `json:"baseUrl,omitempty"`
	LogOutputPath string   `json:"logOutputPath,omitempty"`
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Version:      1,
		ClipsDir:     "./clips",
		CatalogPath:  "./clips/catalog.json",
		EnabledLangs: []string{"en"},
		DefaultVoice: "coral",
		Format:       "wav",
		SampleRate:   24000,
		BaseURL:      "https://api.openai.com/v1",
	}
}

// AvailableLanguages returns the list of all available language codes.
func AvailableLanguages() []string {
	return []string{
		"en", "es", "fr", "de", "it", "pt", "pl", "tr", "ru", "nl",
		"cs", "ar", "zh", "ja", "hu", "ko",
	}
}

// SupportedVoices returns the list of supported OpenAI TTS voices.
func SupportedVoices() []string {
	return []string{"alloy", "echo", "fable", "onyx", "nova", "shimmer", "coral"}
}

// SupportedFormats returns the list of supported audio formats.
func SupportedFormats() []string {
	return []string{"wav", "mp3", "pcm", "opus"}
}

// SupportedSampleRates returns the list of supported sample rates.
func SupportedSampleRates() []int {
	return []int{8000, 11025, 16000, 22050, 24000, 44100, 48000}
}

// ResolveConfigPath determines the config file path based on CORALIE_CONFIG env var or default.
func ResolveConfigPath() (string, error) {
	envConfig := os.Getenv("CORALIE_CONFIG")
	if envConfig == "" {
		// Default: current working directory
		return "./config.json", nil
	}

	// If it ends with .json, treat as full file path
	if strings.HasSuffix(envConfig, ".json") {
		return envConfig, nil
	}

	// Otherwise, treat as directory
	configPath := filepath.Join(envConfig, "config.json")
	return configPath, nil
}

// ResolveEnvPath determines the .env file path based on config file location.
func ResolveEnvPath(configPath string) string {
	configDir := filepath.Dir(configPath)
	envPath := filepath.Join(configDir, ".env")
	return envPath
}

// LoadConfig loads configuration from config file, .env, and environment variables.
func LoadConfig() (*Config, error) {
	configPath, err := ResolveConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	// Start with defaults
	cfg := DefaultConfig()

	// Load .env file (from same directory as config)
	envPath := ResolveEnvPath(configPath)
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			return nil, fmt.Errorf("failed to load .env: %w", err)
		}
	}

	// Load environment variables
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		cfg.OpenAIApiKey = apiKey
	}
	if baseURL := os.Getenv("OPENAI_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	// Load config.json if it exists
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		var fileConfig Config
		if err := json.Unmarshal(data, &fileConfig); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}

		// Merge file config into defaults (file config takes precedence for non-empty values)
		if fileConfig.OpenAIApiKey != "" {
			cfg.OpenAIApiKey = fileConfig.OpenAIApiKey
		}
		if fileConfig.ClipsDir != "" {
			cfg.ClipsDir = fileConfig.ClipsDir
		}
		if fileConfig.CatalogPath != "" {
			cfg.CatalogPath = fileConfig.CatalogPath
		}
		if len(fileConfig.EnabledLangs) > 0 {
			cfg.EnabledLangs = fileConfig.EnabledLangs
		}
		if fileConfig.DefaultVoice != "" {
			cfg.DefaultVoice = fileConfig.DefaultVoice
		}
		if fileConfig.Format != "" {
			cfg.Format = fileConfig.Format
		}
		if fileConfig.SampleRate > 0 {
			cfg.SampleRate = fileConfig.SampleRate
		}
		if fileConfig.BaseURL != "" {
			cfg.BaseURL = fileConfig.BaseURL
		}
		if fileConfig.LogOutputPath != "" {
			cfg.LogOutputPath = fileConfig.LogOutputPath
		}
		cfg.Version = fileConfig.Version
	}

	// Environment variables override everything
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		cfg.OpenAIApiKey = apiKey
	}
	if baseURL := os.Getenv("OPENAI_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	return cfg, nil
}

// SaveConfig saves the configuration to the config file.
func SaveConfig(cfg *Config) error {
	configPath, err := ResolveConfigPath()
	if err != nil {
		return fmt.Errorf("failed to resolve config path: %w", err)
	}

	// Check if parent directory exists
	configDir := filepath.Dir(configPath)
	if configDir != "." && configDir != "" {
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			return fmt.Errorf("config directory does not exist: %s\nRun: mkdir -p %s", configDir, configDir)
		}
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ValidateConfigOrExit validates the configuration and exits with error message if invalid.
func ValidateConfigOrExit(cfg *Config) {
	var errors []string

	// Check API key
	if cfg.OpenAIApiKey == "" {
		errors = append(errors, "OpenAI API key is required (set via config, .env, or OPENAI_API_KEY env var)")
	}

	// Check clips directory
	if cfg.ClipsDir == "" {
		errors = append(errors, "clips directory is required")
	} else {
		// Check if directory exists or can be created
		if err := os.MkdirAll(cfg.ClipsDir, 0755); err != nil {
			errors = append(errors, fmt.Sprintf("cannot create clips directory %s: %v", cfg.ClipsDir, err))
		}
	}

	// Check enabled languages
	if len(cfg.EnabledLangs) == 0 {
		errors = append(errors, "at least one language must be enabled (use 'coralie-clip lang enable <code>')")
	} else {
		// Validate each enabled language
		available := AvailableLanguages()
		availableMap := make(map[string]bool)
		for _, lang := range available {
			availableMap[lang] = true
		}
		for _, lang := range cfg.EnabledLangs {
			if !availableMap[lang] {
				errors = append(errors, fmt.Sprintf("language '%s' is not available", lang))
			}
		}
	}

	// Check voice
	if cfg.DefaultVoice == "" {
		errors = append(errors, "default voice is required")
	} else {
		validVoice := false
		for _, v := range SupportedVoices() {
			if v == cfg.DefaultVoice {
				validVoice = true
				break
			}
		}
		if !validVoice {
			errors = append(errors, fmt.Sprintf("voice '%s' is not supported", cfg.DefaultVoice))
		}
	}

	// Check format
	if cfg.Format == "" {
		errors = append(errors, "format is required")
	} else {
		validFormat := false
		for _, f := range SupportedFormats() {
			if f == cfg.Format {
				validFormat = true
				break
			}
		}
		if !validFormat {
			errors = append(errors, fmt.Sprintf("format '%s' is not supported", cfg.Format))
		}
	}

	// Check sample rate
	if cfg.SampleRate <= 0 {
		errors = append(errors, "sample rate must be positive")
	} else {
		validRate := false
		for _, sr := range SupportedSampleRates() {
			if sr == cfg.SampleRate {
				validRate = true
				break
			}
		}
		if !validRate {
			errors = append(errors, fmt.Sprintf("sample rate %d is not supported", cfg.SampleRate))
		}
	}

	// Check catalog path
	if cfg.CatalogPath == "" {
		errors = append(errors, "catalog path is required")
	} else {
		catalogDir := filepath.Dir(cfg.CatalogPath)
		if catalogDir != "." && catalogDir != "" {
			if err := os.MkdirAll(catalogDir, 0755); err != nil {
				errors = append(errors, fmt.Sprintf("cannot create catalog directory %s: %v", catalogDir, err))
			}
		}
	}

	// Check CORALIE_CONFIG if set
	if envConfig := os.Getenv("CORALIE_CONFIG"); envConfig != "" {
		configPath, _ := ResolveConfigPath()
		configDir := filepath.Dir(configPath)
		if configDir != "." && configDir != "" {
			if _, err := os.Stat(configDir); os.IsNotExist(err) {
				errors = append(errors, fmt.Sprintf("CORALIE_CONFIG directory does not exist: %s\nRun: mkdir -p %s", configDir, configDir))
			}
		}
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Configuration errors:\n")
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", err)
		}
		os.Exit(1)
	}
}

