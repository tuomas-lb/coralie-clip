# OpenAI Integration

coralie-clip uses OpenAI's TTS and STT APIs for audio generation and transcription.

## TTS (Text-to-Speech)

### Endpoint

`POST /audio/speech`

### Model

Currently uses `tts-1` model.

### Request Format

```json
{
  "model": "tts-1",
  "input": "Text to convert to speech",
  "voice": "coral",
  "response_format": "wav",
  "speed": 1.0
}
```

### Supported Voices

- `alloy`
- `echo`
- `fable`
- `onyx`
- `nova`
- `shimmer`
- `coral`

### Supported Response Formats

- `wav` - WAV container
- `mp3` - MP3 encoded
- `pcm` - Raw PCM16 (converted to WAV container on save)
- `opus` - Ogg Opus container

### Sample Rates

Supported sample rates: 8000, 11025, 16000, 22050, 24000, 44100, 48000 Hz

### Token Usage

**Note:** OpenAI TTS API does not return token usage information in the response. The catalog will log "token usage not provided by API" for TTS requests.

## STT (Speech-to-Text)

### Endpoint

`POST /audio/transcriptions`

### Model

Uses `whisper-1` model.

### Request Format

Multipart form data:
- `file` - Audio file
- `model` - Model name (default: `whisper-1`)
- `language` - Optional language code

### Supported Audio Formats

- WAV
- MP3
- Opus
- Other formats supported by Whisper

### Language Detection

- If `language` parameter is provided, it is used
- Otherwise, Whisper auto-detects language
- Language is returned in the response (if available)

### Token Usage

**Note:** OpenAI STT API does not return token usage information in the response. The catalog will log "token usage not provided by API" for STT requests.

## Error Handling

### Retry Logic

The client implements exponential backoff retry for:
- HTTP 429 (rate limit)
- HTTP 5xx (server errors)
- Network errors

Retry configuration:
- Maximum 3 retries
- Initial backoff: 1 second
- Exponential backoff: 2x per retry

### Error Messages

Errors are returned with actionable messages:
- API errors include status code and response body
- Network errors include underlying error details

## Request Logging

If `logOutputPath` is configured, requests are logged as JSON lines:

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "endpoint": "/audio/speech",
  "method": "POST",
  "model": "tts-1",
  "voice": "coral",
  "language": "en",
  "format": "wav",
  "sample_rate": 24000,
  "response_bytes": 12345,
  "request_id": "req_123",
  "token_usage": null,
  "latency": "1.234s",
  "status": 200
}
```

**Privacy:** Text payloads are not logged unless user opts in.

## Rate Limiting

- TTS requests: Subject to OpenAI rate limits
- STT requests (rebuild-catalog): Rate-limited to 500ms between requests

## Base URL

Default: `https://api.openai.com/v1`

Can be overridden via:
- `baseUrl` in config
- `OPENAI_BASE_URL` environment variable

This allows using OpenAI-compatible APIs.

## Limitations

1. **Token Usage:** Not provided by audio APIs, so catalog entries will have empty token usage
2. **Language Detection:** STT language detection may not be 100% accurate
3. **File Size:** Large audio files may timeout (60s timeout configured)
4. **Format Conversion:** PCM format is converted to WAV container on save

## If you are an AI

If you are an AI agent working on this repository:
- You MUST keep this documentation up to date with any behavior changes.
- You MUST document API changes when OpenAI updates their APIs.

