# CLI Commands Reference

Complete reference for all coralie-clip commands.

## Global Configuration Pipeline

Every command that reads configuration follows this pipeline:
1. Resolve config file location (see [Configuration](config.md))
2. Load `.env` file (if present)
3. Read environment variables
4. Load `config.json` (if exists)
5. Merge with CLI overrides
6. Validate configuration
7. Execute command logic

## Commands

### `set <key> <value>`

Set a configuration value.

**Examples:**
```bash
coralie-clip set openai-apikey <apikey>
coralie-clip set clips-dir ./clips
coralie-clip set sample-rate 24000
coralie-clip set format wav
coralie-clip set default-voice coral
coralie-clip set catalog-path ./clips/catalog.json
coralie-clip set base-url https://api.openai.com/v1
coralie-clip set log-output-path ./logs/requests.jsonl
```

**Notes:**
- Secrets (like API keys) are redacted in output
- Configuration is validated before saving
- Invalid configurations will not be saved

### `fetch "<text>" [options]`

Generate an audio clip from text using OpenAI TTS.

**Options:**
- `--lang <code>` - Language code (required if multiple languages enabled)
- `--voice <name>` - Voice name (default: from config)
- `--format <fmt>` - Audio format: wav, mp3, pcm, opus (default: from config)
- `--sample-rate <hz>` - Sample rate in Hz (default: from config)

**Examples:**
```bash
coralie-clip fetch "Hello, world!"
coralie-clip fetch "Bonjour" --lang fr
coralie-clip fetch "Test" --voice nova --format mp3
```

**Language Selection:**
- If `--lang` is provided, it must be enabled
- If `--lang` is omitted:
  - If exactly 1 language is enabled, it is used automatically
  - If multiple languages are enabled, `--lang` is required
  - If no languages are enabled, the command fails

### `lang enable <code|all>`

Enable a language for TTS generation.

**Examples:**
```bash
coralie-clip lang enable en
coralie-clip lang enable all
```

### `lang disable <code|all>`

Disable a language.

**Examples:**
```bash
coralie-clip lang disable fr
coralie-clip lang disable all
```

### `lang list [all]`

List enabled languages or all available languages.

**Examples:**
```bash
coralie-clip lang list        # List enabled languages
coralie-clip lang list all    # List all available languages with status
```

### `voice <name>`

Set the default voice.

**Examples:**
```bash
coralie-clip voice coral
coralie-clip voice nova
```

### `voice list`

List all available voices.

**Example:**
```bash
coralie-clip voice list
```

### `format <wav|mp3|pcm|opus>`

Set the default audio format.

**Examples:**
```bash
coralie-clip format wav
coralie-clip format mp3
```

**Supported Formats:**
- `wav` - WAV container with PCM16
- `mp3` - MP3 encoded audio
- `pcm` - Raw PCM16 (saved as WAV container)
- `opus` - Ogg Opus container

### `find "<query>"`

Search for clips by text content.

**Examples:**
```bash
coralie-clip find "hello"
coralie-clip find "world"
```

**Search Behavior:**
- Case-insensitive substring matching
- Searches both original text and transcription
- Returns: ID, language, text preview, and file path

### `play <id>`

Play an audio clip by ID.

**Example:**
```bash
coralie-clip play abc123def4
```

**Platform Support:**
- macOS: Uses `afplay`
- Linux: Uses `paplay`, `aplay`, `mpg123`, or `ffplay` (first available)
- Windows: Uses PowerShell
- Custom: Set `PLAYER_CMD` environment variable

**Troubleshooting:**
If no player is found, install a player or set `PLAYER_CMD`:
```bash
export PLAYER_CMD=/path/to/player
```

### `rebuild-catalog [--force]`

Rebuild the catalog by scanning the clips directory and using OpenAI STT.

**Options:**
- `--force` - Overwrite existing catalog entries

**Example:**
```bash
coralie-clip rebuild-catalog
coralie-clip rebuild-catalog --force
```

**Process:**
1. Scans clips directory for audio files
2. Uses OpenAI STT to transcribe each file
3. Infers language from STT response or filename
4. Reconstructs catalog entries with metadata
5. Rate-limits STT calls (500ms between requests)

**Notes:**
- Existing entries are preserved unless `--force` is used
- Files already in catalog are skipped unless `--force` is used
- Only processes `.wav`, `.mp3`, and `.opus` files

## Makefile Targets

See [README.md](README.md) for build and installation instructions.

**If you are an AI**

If you are an AI agent working on this repository:
- You MUST keep this documentation up to date with any behavior changes.
- You MUST update command examples when CLI changes are made.

