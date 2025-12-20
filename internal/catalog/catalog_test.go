package catalog

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadAndSaveCatalog(t *testing.T) {
	tmpDir := t.TempDir()
	catalogPath := filepath.Join(tmpDir, "catalog.json")

	// Test loading non-existent catalog
	catalog, err := LoadCatalog(catalogPath)
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	if len(catalog.Entries) != 0 {
		t.Errorf("LoadCatalog() should return empty catalog for non-existent file")
	}

	// Add an entry
	entry := Entry{
		ID:         "test123",
		CreatedAt:  time.Now(),
		Lang:       "en",
		Text:       "Hello world",
		Voice:      "coral",
		Format:     "wav",
		SampleRate: 24000,
		FilePath:   "./clips/test123_en_coral_24000.wav",
	}
	catalog.AddEntry(entry)

	// Save catalog
	if err := SaveCatalog(catalog, catalogPath); err != nil {
		t.Fatalf("SaveCatalog() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(catalogPath); os.IsNotExist(err) {
		t.Fatalf("catalog file was not created")
	}

	// Load catalog again
	loaded, err := LoadCatalog(catalogPath)
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}

	if len(loaded.Entries) != 1 {
		t.Fatalf("LoadCatalog() expected 1 entry, got %d", len(loaded.Entries))
	}

	if loaded.Entries[0].ID != "test123" {
		t.Errorf("LoadCatalog() entry ID = %v, want test123", loaded.Entries[0].ID)
	}
}

func TestAddEntry(t *testing.T) {
	catalog := &Catalog{Entries: []Entry{}}

	entry1 := Entry{ID: "id1", Text: "Hello"}
	entry2 := Entry{ID: "id2", Text: "World"}

	catalog.AddEntry(entry1)
	catalog.AddEntry(entry2)

	if len(catalog.Entries) != 2 {
		t.Errorf("AddEntry() expected 2 entries, got %d", len(catalog.Entries))
	}

	// Update existing entry
	entry1Updated := Entry{ID: "id1", Text: "Hello Updated"}
	catalog.AddEntry(entry1Updated)

	if len(catalog.Entries) != 2 {
		t.Errorf("AddEntry() expected 2 entries after update, got %d", len(catalog.Entries))
	}

	if catalog.Entries[0].Text != "Hello Updated" {
		t.Errorf("AddEntry() entry text = %v, want Hello Updated", catalog.Entries[0].Text)
	}
}

func TestFindEntry(t *testing.T) {
	catalog := &Catalog{Entries: []Entry{}}
	entry := Entry{ID: "test123", Text: "Hello"}
	catalog.AddEntry(entry)

	found := catalog.FindEntry("test123")
	if found == nil {
		t.Fatalf("FindEntry() returned nil for existing entry")
	}
	if found.Text != "Hello" {
		t.Errorf("FindEntry() entry text = %v, want Hello", found.Text)
	}

	notFound := catalog.FindEntry("nonexistent")
	if notFound != nil {
		t.Errorf("FindEntry() should return nil for non-existent entry")
	}
}

func TestSearch(t *testing.T) {
	catalog := &Catalog{Entries: []Entry{}}
	catalog.AddEntry(Entry{ID: "1", Text: "Hello world", Transcription: ""})
	catalog.AddEntry(Entry{ID: "2", Text: "Goodbye", Transcription: "farewell message"})
	catalog.AddEntry(Entry{ID: "3", Text: "Test", Transcription: ""})

	results := catalog.Search("hello")
	if len(results) != 1 {
		t.Errorf("Search() expected 1 result for 'hello', got %d", len(results))
	}

	results = catalog.Search("world")
	if len(results) != 1 {
		t.Errorf("Search() expected 1 result for 'world', got %d", len(results))
	}

	results = catalog.Search("farewell")
	if len(results) != 1 {
		t.Errorf("Search() expected 1 result for 'farewell', got %d", len(results))
	}

	results = catalog.Search("HELLO")
	if len(results) != 1 {
		t.Errorf("Search() should be case-insensitive, expected 1 result for 'HELLO', got %d", len(results))
	}

	results = catalog.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("Search() expected 0 results for 'nonexistent', got %d", len(results))
	}
}

