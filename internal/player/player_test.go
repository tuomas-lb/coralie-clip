package player

import (
	"os"
	"testing"
)

func TestNewPlayer(t *testing.T) {
	// Save original env
	originalEnv := os.Getenv("PLAYER_CMD")
	defer os.Setenv("PLAYER_CMD", originalEnv)

	// Test with custom player
	os.Setenv("PLAYER_CMD", "custom-player")
	p := NewPlayer()
	if p.cmd != "custom-player" {
		t.Errorf("NewPlayer() with PLAYER_CMD = %v, want custom-player", p.cmd)
	}

	// Test without custom player (platform detection)
	os.Unsetenv("PLAYER_CMD")
	p = NewPlayer()
	// Should detect platform-specific player or be empty
	_ = p
}

func TestPlayerIsAvailable(t *testing.T) {
	p := NewPlayer()
	// Just test that it doesn't panic
	_ = p.IsAvailable()
}

