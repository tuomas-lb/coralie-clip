# Configuration

Configuration for coralie-clip can be provided via config file, `.env` file, or environment variables.

## Config File Location

### Default Location

By default, the config file is located at:
- `./config.json` (relative to current working directory)

### Custom Location

Override via `CORALIE_CONFIG` environment variable:

**As file path:**
```bash
export CORALIE_CONFIG=/etc/coralie/config.json
```

**As directory:**
```bash
export CORALIE_CONFIG=/etc/coralie
# Config file becomes: /etc/coralie/config.json
```

**Notes:**
- If `CORALIE_CONFIG` points to a non-existent directory, the tool will error with a suggestion to create it
- Parent directories are not auto-created outside CWD without explicit user action

## Configuration Sources (Priority Order)

1. **Environment variables** (highest priority)
2. **config.json** file
3. **.env** file (lowest priority)

The `.env` file is loaded from the same directory as the config file:
- If config is `/etc/coralie/config.json`, `.env` is loaded from `/etc/coralie/.env`
- Otherwise, `.env` is loaded from `./.env`

## Config File Format

JSON format with schema version:

```json
{
  "version": 1,
  "openaiApiKey": "sk-...",
  "clipsDir": "./clips",
  "catalogPath": "./clips/catalog.json",
  "enabledLangs": ["en", "fr"],
  "defaultVoice": "coral",
  "format": "wav",
  "sampleRate": 24000,
  "baseUrl": "https://api.openai.com/v1",
  "logOutputPath": "./logs/requests.jsonl"
}
```

## Configuration Fields

### Required Fields

- **openaiApiKey** (or `OPENAI_API_KEY` env var) - OpenAI API key
- **clipsDir** - Directory for storing audio clips (default: `./clips`)
- **enabledLangs** - Array of enabled language codes (must have at least 1)

### Optional Fields

- **catalogPath** - Path to catalog JSON file (default: `./clips/catalog.json`)
- **defaultVoice** - Default voice name (default: `coral`)
- **format** - Default audio format: `wav`, `mp3`, `pcm`, `opus` (default: `wav`)
- **sampleRate** - Default sample rate in Hz (default: `24000`)
- **baseUrl** - OpenAI API base URL (default: `https://api.openai.com/v1`, override via `OPENAI_BASE_URL`)
- **logOutputPath** - Path to request log file (optional)

## Environment Variables

- `OPENAI_API_KEY` - OpenAI API key
- `OPENAI_BASE_URL` - OpenAI API base URL
- `CORALIE_CONFIG` - Config file location override
- `PLAYER_CMD` - Custom audio player command

## Available Languages

Static list of available language codes:
- `en` - English
- `es` - Spanish
- `fr` - French
- `de` - German
- `it` - Italian
- `pt` - Portuguese
- `pl` - Polish
- `tr` - Turkish
- `ru` - Russian
- `nl` - Dutch
- `cs` - Czech
- `ar` - Arabic
- `zh` - Chinese
- `ja` - Japanese
- `hu` - Hungarian
- `ko` - Korean

## Supported Voices

- `alloy`
- `echo`
- `fable`
- `onyx`
- `nova`
- `shimmer`
- `coral` (default)

## Supported Formats

- `wav` - WAV container with PCM16
- `mp3` - MP3 encoded audio
- `pcm` - Raw PCM16 (saved as WAV container)
- `opus` - Ogg Opus container

## Supported Sample Rates

- 8000 Hz
- 11025 Hz
- 16000 Hz
- 22050 Hz
- 24000 Hz (default)
- 44100 Hz
- 48000 Hz

## Validation

Configuration is validated before use. The following checks are performed:

- API key is present
- Clips directory exists or is creatable
- At least one language is enabled
- Voice is in supported list
- Format is supported
- Sample rate is within allowed range
- Catalog file path is writable (parent directory exists/creatable)
- `CORALIE_CONFIG` directory exists (if set)

Invalid configurations will cause the tool to exit with an error message.

## .env File Example

Create `.env.example` in the repository (never commit real `.env`):

```bash
OPENAI_API_KEY=your-api-key-here
OPENAI_BASE_URL=https://api.openai.com/v1
CORALIE_CONFIG=
```

**If you are an AI**

If you are an AI agent working on this repository:
- You MUST keep this documentation up to date with any behavior changes.
- You MUST document new configuration options when added.





