// Package catalog manages the clip catalog storage and search.
package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Entry represents a catalog entry for an audio clip.
type Entry struct {
	ID            string                 `json:"id"`
	CreatedAt     time.Time              `json:"createdAt"`
	Lang          string                 `json:"lang"`
	Text          string                 `json:"text"`
	Voice         string                 `json:"voice"`
	Format        string                 `json:"format"`
	SampleRate    int                    `json:"sampleRate"`
	FilePath      string                 `json:"filePath"`
	Transcription string                 `json:"transcription,omitempty"`
	ProviderMeta  map[string]interface{} `json:"providerMeta,omitempty"`
	TokenUsage    map[string]interface{} `json:"tokenUsage,omitempty"`
}

// Catalog represents the collection of clip entries.
type Catalog struct {
	Entries []Entry `json:"entries"`
}

// LoadCatalog loads the catalog from a file.
func LoadCatalog(path string) (*Catalog, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Catalog{Entries: []Entry{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read catalog: %w", err)
	}

	var catalog Catalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse catalog: %w", err)
	}

	return &catalog, nil
}

// SaveCatalog saves the catalog to a file atomically.
func SaveCatalog(catalog *Catalog, path string) error {
	// Create parent directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create catalog directory: %w", err)
	}

	// Write to temporary file
	tmpPath := path + ".tmp"
	data, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal catalog: %w", err)
	}

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write catalog: %w", err)
	}

	// Sync to disk (best effort)
	if f, err := os.OpenFile(tmpPath, os.O_RDWR, 0644); err == nil {
		f.Sync()
		f.Close()
	}

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename catalog: %w", err)
	}

	return nil
}

// AddEntry adds or updates an entry in the catalog.
func (c *Catalog) AddEntry(entry Entry) {
	// Check if entry with same ID exists
	for i, e := range c.Entries {
		if e.ID == entry.ID {
			c.Entries[i] = entry
			return
		}
	}
	// Add new entry
	c.Entries = append(c.Entries, entry)
}

// FindEntry finds an entry by ID.
func (c *Catalog) FindEntry(id string) *Entry {
	for _, e := range c.Entries {
		if e.ID == id {
			return &e
		}
	}
	return nil
}

// Search searches for entries matching the query (case-insensitive substring match).
func (c *Catalog) Search(query string) []Entry {
	query = strings.ToLower(query)
	var results []Entry

	for _, e := range c.Entries {
		textMatch := strings.Contains(strings.ToLower(e.Text), query)
		transcriptionMatch := e.Transcription != "" && strings.Contains(strings.ToLower(e.Transcription), query)
		
		if textMatch || transcriptionMatch {
			results = append(results, e)
		}
	}

	return results
}

// GetEntryByFilePath finds an entry by file path.
func (c *Catalog) GetEntryByFilePath(filePath string) *Entry {
	for _, e := range c.Entries {
		if e.FilePath == filePath {
			return &e
		}
	}
	return nil
}

