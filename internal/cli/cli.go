// Package cli implements the command-line interface.
package cli

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coralie/coralie-clip/internal/audio"
	"github.com/coralie/coralie-clip/internal/catalog"
	"github.com/coralie/coralie-clip/internal/config"
	"github.com/coralie/coralie-clip/internal/logging"
	"github.com/coralie/coralie-clip/internal/openai"
	"github.com/coralie/coralie-clip/internal/player"
)

// App represents the CLI application.
type App struct {
	cfg     *config.Config
	catalog *catalog.Catalog
	client  *openai.Client
	logger  *logging.Logger
	player  *player.Player
}

// NewApp creates a new CLI application instance.
func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	config.ValidateConfigOrExit(cfg)

	catalog, err := catalog.LoadCatalog(cfg.CatalogPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load catalog: %w", err)
	}

	client := openai.NewClient(cfg.OpenAIApiKey, cfg.BaseURL)

	var logger *logging.Logger
	if cfg.LogOutputPath != "" {
		logger, err = logging.NewLogger(cfg.LogOutputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
	}

	app := &App{
		cfg:     cfg,
		catalog: catalog,
		client:  client,
		logger:  logger,
		player:  player.NewPlayer(),
	}

	return app, nil
}

// Close cleans up resources.
func (a *App) Close() error {
	if a.logger != nil {
		return a.logger.Close()
	}
	return nil
}

// generateID generates a short, stable ID for a clip.
func generateID() string {
	b := make([]byte, 6)
	rand.Read(b)
	return strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b))[:10]
}

// RunSetCommandStandalone handles the `set` command without requiring a fully validated config.
func RunSetCommandStandalone(key, value string) error {
	// Load config without validation
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Update config
	switch key {
	case "openai-apikey":
		cfg.OpenAIApiKey = value
		fmt.Println("openai-apikey set (redacted)")
	case "clips-dir":
		cfg.ClipsDir = value
		fmt.Printf("clips-dir set to %s\n", value)
	case "sample-rate":
		var rate int
		if _, err := fmt.Sscanf(value, "%d", &rate); err != nil {
			return fmt.Errorf("invalid sample rate: %s", value)
		}
		cfg.SampleRate = rate
		fmt.Printf("sample-rate set to %d\n", rate)
	case "format":
		cfg.Format = value
		fmt.Printf("format set to %s\n", value)
	case "catalog-path":
		cfg.CatalogPath = value
		fmt.Printf("catalog-path set to %s\n", value)
	case "default-voice":
		cfg.DefaultVoice = value
		fmt.Printf("default-voice set to %s\n", value)
	case "base-url":
		cfg.BaseURL = value
		fmt.Printf("base-url set to %s\n", value)
	case "log-output-path":
		cfg.LogOutputPath = value
		fmt.Printf("log-output-path set to %s\n", value)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	// Validate after update
	config.ValidateConfigOrExit(cfg)

	// Save config
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// RunSetCommand handles the `set` command.
func (a *App) RunSetCommand(key, value string) error {
	return RunSetCommandStandalone(key, value)
}

// RunFetchCommand handles the `fetch` command.
func (a *App) RunFetchCommand(text, lang, voice, format string, sampleRate int) error {
	// Determine language
	if lang == "" {
		if len(a.cfg.EnabledLangs) == 0 {
			return fmt.Errorf("no languages enabled. Use 'coralie-clip lang enable <code>' to enable a language")
		}
		if len(a.cfg.EnabledLangs) > 1 {
			return fmt.Errorf("multiple languages enabled. Please specify --lang <code>")
		}
		lang = a.cfg.EnabledLangs[0]
	} else {
		// Validate language is enabled
		enabled := false
		for _, l := range a.cfg.EnabledLangs {
			if l == lang {
				enabled = true
				break
			}
		}
		if !enabled {
			return fmt.Errorf("language '%s' is not enabled. Use 'coralie-clip lang enable %s'", lang, lang)
		}
	}

	// Use defaults if not provided
	if voice == "" {
		voice = a.cfg.DefaultVoice
	}
	if format == "" {
		format = a.cfg.Format
	}
	if sampleRate == 0 {
		sampleRate = a.cfg.SampleRate
	}

	// Generate ID
	id := generateID()

	// Map format to OpenAI response format
	openaiFormat := format
	if format == "pcm" {
		openaiFormat = "pcm"
	} else if format == "opus" {
		openaiFormat = "opus"
	} else if format == "wav" {
		openaiFormat = "wav"
	} else if format == "mp3" {
		openaiFormat = "mp3"
	}

	// Call OpenAI TTS
	ctx := context.Background()
	ttsReq := openai.TTSRequest{
		Model:          "tts-1",
		Input:          text,
		Voice:          voice,
		ResponseFormat: openaiFormat,
	}

	start := time.Now()
	resp, err := a.client.TTS(ctx, ttsReq)
	if err != nil {
		return fmt.Errorf("TTS request failed: %w", err)
	}

	// Log request
	if a.logger != nil {
		logEntry := openai.RequestLog{
			Timestamp:     start,
			Endpoint:      "/audio/speech",
			Method:        "POST",
			Model:         "tts-1",
			Voice:         voice,
			Language:      lang,
			Format:        format,
			SampleRate:    sampleRate,
			ResponseBytes: len(resp.AudioData),
			RequestID:     resp.RequestID,
			TokenUsage:    resp.TokenUsage,
			Latency:       logging.FormatDuration(resp.Latency),
			Status:        200,
		}
		a.logger.LogRequest(logEntry)
	}

	// Generate filename
	fileName := audio.GenerateFileName(id, lang, voice, sampleRate, format)
	filePath := filepath.Join(a.cfg.ClipsDir, fileName)

	// Save audio file
	if err := audio.SaveAudio(resp.AudioData, filePath, format); err != nil {
		return fmt.Errorf("failed to save audio: %w", err)
	}

	// Add to catalog
	entry := catalog.Entry{
		ID:            id,
		CreatedAt:     time.Now(),
		Lang:          lang,
		Text:          text,
		Voice:         voice,
		Format:        format,
		SampleRate:    sampleRate,
		FilePath:      filePath,
		ProviderMeta: map[string]interface{}{
			"model":     "tts-1",
			"requestId": resp.RequestID,
		},
		TokenUsage: map[string]interface{}{
			"inputTokens":  resp.TokenUsage.InputTokens,
			"outputTokens": resp.TokenUsage.OutputTokens,
			"totalTokens":  resp.TokenUsage.TotalTokens,
		},
	}

	a.catalog.AddEntry(entry)

	// Save catalog
	if err := catalog.SaveCatalog(a.catalog, a.cfg.CatalogPath); err != nil {
		return fmt.Errorf("failed to save catalog: %w", err)
	}

	fmt.Printf("Created clip: %s\n", id)
	fmt.Printf("File: %s\n", filePath)

	return nil
}

// RunLangCommand handles the `lang` subcommands.
func (a *App) RunLangCommand(action, code string) error {
	available := config.AvailableLanguages()
	availableMap := make(map[string]bool)
	for _, lang := range available {
		availableMap[lang] = true
	}

	switch action {
	case "enable":
		if code == "all" {
			a.cfg.EnabledLangs = available
		} else {
			if !availableMap[code] {
				return fmt.Errorf("language '%s' is not available", code)
			}
			// Add if not already enabled
			found := false
			for _, lang := range a.cfg.EnabledLangs {
				if lang == code {
					found = true
					break
				}
			}
			if !found {
				a.cfg.EnabledLangs = append(a.cfg.EnabledLangs, code)
			}
		}
	case "disable":
		if code == "all" {
			a.cfg.EnabledLangs = []string{}
		} else {
			// Remove from enabled list
			var newLangs []string
			for _, lang := range a.cfg.EnabledLangs {
				if lang != code {
					newLangs = append(newLangs, lang)
				}
			}
			a.cfg.EnabledLangs = newLangs
		}
	case "list":
		if code == "all" {
			// List all available languages
			fmt.Println("Available languages:")
			for _, lang := range available {
				enabled := false
				for _, e := range a.cfg.EnabledLangs {
					if e == lang {
						enabled = true
						break
					}
				}
				status := "disabled"
				if enabled {
					status = "enabled"
				}
				fmt.Printf("  %s (%s)\n", lang, status)
			}
		} else {
			// List enabled languages
			if len(a.cfg.EnabledLangs) == 0 {
				fmt.Println("No languages enabled.")
			} else {
				fmt.Println("Enabled languages:")
				for _, lang := range a.cfg.EnabledLangs {
					fmt.Printf("  %s\n", lang)
				}
			}
		}
	default:
		return fmt.Errorf("unknown lang action: %s", action)
	}

	// Save config if changed
	if action == "enable" || action == "disable" {
		config.ValidateConfigOrExit(a.cfg)
		if err := config.SaveConfig(a.cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
	}

	return nil
}

// RunVoiceCommand handles the `voice` command.
func (a *App) RunVoiceCommand(action, name string) error {
	switch action {
	case "list":
		fmt.Println("Available voices:")
		for _, voice := range config.SupportedVoices() {
			if voice == a.cfg.DefaultVoice {
				fmt.Printf("  %s (default)\n", voice)
			} else {
				fmt.Printf("  %s\n", voice)
			}
		}
	case "set":
		// Validate voice
		valid := false
		for _, v := range config.SupportedVoices() {
			if v == name {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("voice '%s' is not supported", name)
		}
		a.cfg.DefaultVoice = name
		config.ValidateConfigOrExit(a.cfg)
		if err := config.SaveConfig(a.cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Printf("Default voice set to %s\n", name)
	default:
		return fmt.Errorf("unknown voice action: %s (use 'list' or set a voice name)", action)
	}
	return nil
}

// RunFormatCommand handles the `format` command.
func (a *App) RunFormatCommand(format string) error {
	// Validate format
	valid := false
	for _, f := range config.SupportedFormats() {
		if f == format {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("format '%s' is not supported. Supported formats: %s", format, strings.Join(config.SupportedFormats(), ", "))
	}

	a.cfg.Format = format
	config.ValidateConfigOrExit(a.cfg)
	if err := config.SaveConfig(a.cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	fmt.Printf("Format set to %s\n", format)
	return nil
}

// RunFindCommand handles the `find` command.
func (a *App) RunFindCommand(query string) error {
	results := a.catalog.Search(query)
	if len(results) == 0 {
		fmt.Printf("No clips found matching '%s'\n", query)
		return nil
	}

	fmt.Printf("Found %d clip(s):\n", len(results))
	for _, entry := range results {
		preview := entry.Text
		if len(preview) > 50 {
			preview = preview[:50] + "..."
		}
		fmt.Printf("  %s [%s] %s - %s\n", entry.ID, entry.Lang, preview, entry.FilePath)
	}
	return nil
}

// RunPlayCommand handles the `play` command.
func (a *App) RunPlayCommand(id string) error {
	entry := a.catalog.FindEntry(id)
	if entry == nil {
		return fmt.Errorf("clip not found: %s", id)
	}

	if err := a.player.Play(entry.FilePath); err != nil {
		return fmt.Errorf("failed to play audio: %w", err)
	}

	return nil
}

// RunRebuildCatalogCommand handles the `rebuild-catalog` command.
func (a *App) RunRebuildCatalogCommand(force bool) error {
	// Check if catalog exists and has entries
	if !force && len(a.catalog.Entries) > 0 {
		return fmt.Errorf("catalog already exists. Use --force to overwrite")
	}

	// Scan clips directory
	clipsDir := a.cfg.ClipsDir
	files, err := os.ReadDir(clipsDir)
	if err != nil {
		return fmt.Errorf("failed to read clips directory: %w", err)
	}

	// Create new catalog
	newCatalog := &catalog.Catalog{Entries: []catalog.Entry{}}

	// Process each audio file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		if ext != ".wav" && ext != ".mp3" && ext != ".opus" {
			continue
		}

		filePath := filepath.Join(clipsDir, fileName)

		// Check if already in catalog
		existing := a.catalog.GetEntryByFilePath(filePath)
		if existing != nil && !force {
			// Keep existing entry
			newCatalog.AddEntry(*existing)
			continue
		}

		// Use STT to transcribe
		ctx := context.Background()
		audioFile, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to open %s: %v\n", filePath, err)
			continue
		}

		sttReq := openai.STTRequest{
			File:     audioFile,
			Filename: fileName,
		}

		resp, err := a.client.STT(ctx, sttReq)
		audioFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to transcribe %s: %v\n", filePath, err)
			continue
		}

		// Parse filename to extract metadata
		// Format: <id>_<lang>_<voice>_<sr>.<ext>
		parts := strings.Split(strings.TrimSuffix(fileName, ext), "_")
		id := parts[0]
		lang := ""
		voice := ""
		sampleRate := 24000

		if len(parts) >= 2 {
			lang = parts[1]
		}
		if len(parts) >= 3 {
			voice = parts[2]
		}
		if len(parts) >= 4 {
			fmt.Sscanf(parts[3], "%d", &sampleRate)
		}

		// Use detected language from STT if available
		if resp.Language != "" {
			lang = resp.Language
		}

		// Create entry
		entry := catalog.Entry{
			ID:            id,
			CreatedAt:     time.Now(),
			Lang:          lang,
			Text:          resp.Text,
			Transcription: resp.Text,
			Voice:         voice,
			Format:        strings.TrimPrefix(ext, "."),
			SampleRate:    sampleRate,
			FilePath:      filePath,
			ProviderMeta: map[string]interface{}{
				"model":     "whisper-1",
				"requestId": resp.RequestID,
			},
		}

		newCatalog.AddEntry(entry)

		// Rate limit: sleep between requests
		time.Sleep(500 * time.Millisecond)
	}

	// Save new catalog
	if err := catalog.SaveCatalog(newCatalog, a.cfg.CatalogPath); err != nil {
		return fmt.Errorf("failed to save catalog: %w", err)
	}

	fmt.Printf("Rebuilt catalog with %d entries\n", len(newCatalog.Entries))
	return nil
}

