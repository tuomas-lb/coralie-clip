// Package player handles audio playback on the local machine.
package player

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Player handles audio playback.
type Player struct {
	cmd string
}

// NewPlayer creates a new player instance.
func NewPlayer() *Player {
	// Check for custom player command
	if cmd := os.Getenv("PLAYER_CMD"); cmd != "" {
		return &Player{cmd: cmd}
	}

	// Detect platform-specific player
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "afplay"
	case "linux":
		// Try common Linux players
		for _, player := range []string{"paplay", "aplay", "mpg123", "ffplay"} {
			if _, err := exec.LookPath(player); err == nil {
				cmd = player
				break
			}
		}
	case "windows":
		cmd = "powershell"
	default:
		cmd = ""
	}

	return &Player{cmd: cmd}
}

// Play plays an audio file.
func (p *Player) Play(filePath string) error {
	if p.cmd == "" {
		return fmt.Errorf("no audio player found. Install a player or set PLAYER_CMD environment variable")
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("audio file not found: %s", filePath)
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve file path: %w", err)
	}

	// Build command
	var cmd *exec.Cmd
	switch p.cmd {
	case "afplay":
		cmd = exec.Command("afplay", absPath)
	case "paplay", "aplay", "mpg123", "ffplay":
		cmd = exec.Command(p.cmd, absPath)
	case "powershell":
		// Windows PowerShell command
		psCmd := fmt.Sprintf(`(New-Object Media.SoundPlayer "%s").PlaySync()`, absPath)
		cmd = exec.Command("powershell", "-Command", psCmd)
	default:
		// Try to run as-is (user-provided command)
		cmd = exec.Command(p.cmd, absPath)
	}

	// Run command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to play audio: %w", err)
	}

	return nil
}

// IsAvailable checks if a player is available.
func (p *Player) IsAvailable() bool {
	if p.cmd == "" {
		return false
	}
	_, err := exec.LookPath(p.cmd)
	return err == nil
}

