// Package audio handles audio file operations and format conversions.
package audio

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SaveAudio saves audio data to a file with the specified format.
func SaveAudio(data []byte, filePath string, format string) error {
	// Create parent directory if needed
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// For formats that need conversion, we'll handle them
	// For now, OpenAI TTS returns the format directly, so we can save as-is
	// except for PCM which needs WAV container
	if format == "pcm" {
		// Convert PCM to WAV
		wavData, err := PCMToWAV(data, 24000) // Default sample rate, should be configurable
		if err != nil {
			return fmt.Errorf("failed to convert PCM to WAV: %w", err)
		}
		return os.WriteFile(filePath, wavData, 0644)
	}

	// For other formats, save directly
	return os.WriteFile(filePath, data, 0644)
}

// PCMToWAV converts raw PCM data to WAV format.
func PCMToWAV(pcmData []byte, sampleRate int) ([]byte, error) {
	// WAV header structure
	// RIFF header
	header := make([]byte, 44)
	copy(header[0:4], "RIFF")
	
	fileSize := uint32(len(pcmData) + 36)
	binary.LittleEndian.PutUint32(header[4:8], fileSize)
	
	copy(header[8:12], "WAVE")
	
	// fmt chunk
	copy(header[12:16], "fmt ")
	binary.LittleEndian.PutUint32(header[16:20], 16) // fmt chunk size
	binary.LittleEndian.PutUint16(header[20:22], 1)   // audio format (PCM)
	binary.LittleEndian.PutUint16(header[22:24], 1)   // num channels (mono)
	binary.LittleEndian.PutUint32(header[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(header[28:32], uint32(sampleRate*2)) // byte rate
	binary.LittleEndian.PutUint16(header[32:34], 2)   // block align
	binary.LittleEndian.PutUint16(header[34:36], 16)  // bits per sample
	
	// data chunk
	copy(header[36:40], "data")
	binary.LittleEndian.PutUint32(header[40:44], uint32(len(pcmData)))
	
	// Combine header and PCM data
	wavData := append(header, pcmData...)
	return wavData, nil
}

// GetFileExtension returns the file extension for a given format.
func GetFileExtension(format string) string {
	format = strings.ToLower(format)
	switch format {
	case "wav":
		return "wav"
	case "mp3":
		return "mp3"
	case "pcm":
		return "wav" // PCM is saved as WAV
	case "opus":
		return "opus"
	default:
		return "wav"
	}
}

// GenerateFileName generates a deterministic filename for a clip.
func GenerateFileName(id, lang, voice string, sampleRate int, format string) string {
	ext := GetFileExtension(format)
	return fmt.Sprintf("%s_%s_%s_%d.%s", id, lang, voice, sampleRate, ext)
}

