// Package openai provides a client for OpenAI TTS and STT APIs.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// Client handles OpenAI API requests.
type Client struct {
	apiKey    string
	baseURL   string
	httpClient *http.Client
}

// NewClient creates a new OpenAI client.
func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey: apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// TTSRequest represents a text-to-speech request.
type TTSRequest struct {
	Model       string `json:"model"`
	Input       string `json:"input"`
	Voice       string `json:"voice"`
	ResponseFormat string `json:"response_format"`
	Speed       float64 `json:"speed,omitempty"`
}

// TTSResponse represents a text-to-speech response.
type TTSResponse struct {
	AudioData   []byte
	RequestID   string
	TokenUsage  *TokenUsage
	Latency     time.Duration
}

// STTRequest represents a speech-to-text request.
type STTRequest struct {
	File        io.Reader
	Filename    string
	Language    string
	Model       string
}

// STTResponse represents a speech-to-text response.
type STTResponse struct {
	Text        string
	Language    string
	RequestID   string
	TokenUsage  *TokenUsage
	Latency     time.Duration
}

// TokenUsage represents token usage information.
type TokenUsage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
	TotalTokens  int `json:"total_tokens,omitempty"`
}

// RequestLog represents a logged API request.
type RequestLog struct {
	Timestamp   time.Time `json:"timestamp"`
	Endpoint    string    `json:"endpoint"`
	Method      string    `json:"method"`
	Model       string    `json:"model,omitempty"`
	Voice       string    `json:"voice,omitempty"`
	Language    string    `json:"language,omitempty"`
	Format      string    `json:"format,omitempty"`
	SampleRate  int       `json:"sample_rate,omitempty"`
	ResponseBytes int     `json:"response_bytes"`
	RequestID   string    `json:"request_id,omitempty"`
	TokenUsage  *TokenUsage `json:"token_usage,omitempty"`
	Latency     string    `json:"latency"`
	Status      int       `json:"status"`
	Error       string    `json:"error,omitempty"`
}

// TTS generates audio from text using OpenAI TTS API.
func (c *Client) TTS(ctx context.Context, req TTSRequest) (*TTSResponse, error) {
	start := time.Now()

	url := fmt.Sprintf("%s/audio/speech", c.baseURL)

	reqBody := TTSRequest{
		Model:          req.Model,
		Input:          req.Input,
		Voice:          req.Voice,
		ResponseFormat: req.ResponseFormat,
		Speed:          1.0,
	}
	if req.Speed > 0 {
		reqBody.Speed = req.Speed
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	var lastErr error
	maxRetries := 3
	backoff := time.Second

	for i := 0; i < maxRetries; i++ {
		resp, lastErr = c.httpClient.Do(httpReq)
		if lastErr == nil {
			if resp.StatusCode == http.StatusOK || resp.StatusCode < 500 {
				break
			}
			// Retry on 429 or 5xx
			if resp.StatusCode == 429 || resp.StatusCode >= 500 {
				resp.Body.Close()
				if i < maxRetries-1 {
					time.Sleep(backoff)
					backoff *= 2
					continue
				}
			} else {
				// Non-retryable error
				break
			}
		} else {
			if i < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("request failed: %w", lastErr)
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	requestID := resp.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = resp.Header.Get("Request-Id")
	}

	// OpenAI TTS API doesn't return token usage in response headers
	// We'll log that it's not available
	tokenUsage := &TokenUsage{}

	return &TTSResponse{
		AudioData:  audioData,
		RequestID:  requestID,
		TokenUsage: tokenUsage,
		Latency:    latency,
	}, nil
}

// STT transcribes audio using OpenAI STT API.
func (c *Client) STT(ctx context.Context, req STTRequest) (*STTResponse, error) {
	start := time.Now()

	url := fmt.Sprintf("%s/audio/transcriptions", c.baseURL)

	// Create multipart form
	// We'll use a simple approach: read the file and create form data
	fileData, err := io.ReadAll(req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create multipart form
	var formBody bytes.Buffer
	writer := multipart.NewWriter(&formBody)

	// Add file
	part, err := writer.CreateFormFile("file", req.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	// Add model
	model := req.Model
	if model == "" {
		model = "whisper-1"
	}
	if err := writer.WriteField("model", model); err != nil {
		return nil, fmt.Errorf("failed to write model field: %w", err)
	}

	// Add language if provided
	if req.Language != "" {
		if err := writer.WriteField("language", req.Language); err != nil {
			return nil, fmt.Errorf("failed to write language field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, &formBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	var resp *http.Response
	var lastErr error
	maxRetries := 3
	backoff := time.Second

	for i := 0; i < maxRetries; i++ {
		resp, lastErr = c.httpClient.Do(httpReq)
		if lastErr == nil {
			if resp.StatusCode == http.StatusOK || resp.StatusCode < 500 {
				break
			}
			if resp.StatusCode == 429 || resp.StatusCode >= 500 {
				resp.Body.Close()
				if i < maxRetries-1 {
					time.Sleep(backoff)
					backoff *= 2
					continue
				}
			} else {
				break
			}
		} else {
			if i < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("request failed: %w", lastErr)
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var sttResp struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sttResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	requestID := resp.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = resp.Header.Get("Request-Id")
	}

	// Language detection: OpenAI Whisper doesn't return language in transcription endpoint
	// We'll use a heuristic or leave it empty
	detectedLang := req.Language
	if detectedLang == "" {
		// Could implement language detection here, but for now leave empty
		detectedLang = ""
	}

	tokenUsage := &TokenUsage{}

	return &STTResponse{
		Text:       sttResp.Text,
		Language:   detectedLang,
		RequestID:  requestID,
		TokenUsage: tokenUsage,
		Latency:    latency,
	}, nil
}

