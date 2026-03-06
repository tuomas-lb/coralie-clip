# Go Development Reference

## Prerequisites
- Go 1.22+
- OpenAI API key
- Audio player: `afplay` (macOS), `paplay`/`aplay` (Linux)

## Project Layout
```
coralie-cli/
  cmd/coralie-clip/main.go    -- Entry point, CLI dispatch
  internal/
    cli/cli.go                -- App struct, all command handlers
    config/config.go          -- Config types, load/save/validate
    catalog/catalog.go        -- Catalog CRUD, search, persistence
    openai/client.go          -- TTS and STT HTTP client
    audio/audio.go            -- PCM-to-WAV, filename generation
    player/player.go          -- Platform-aware audio playback
    logging/logging.go        -- JSON line logger
  docs/                       -- Documentation
  clips/                      -- Generated audio (gitignored)
  config.json                 -- User config (gitignored)
  .env                        -- Secrets (gitignored)
```

## Key Patterns

### App Struct (internal/cli)
Central orchestrator -- holds config, catalog, OpenAI client, logger, player.
```
NewApp() -> loads config, catalog, creates client
app.RunXxxCommand(...) -> executes command logic
app.Close() -> cleanup (close logger)
```

### Configuration Layering
1. `DefaultConfig()` sets baseline
2. `.env` loaded via godotenv
3. `config.json` merged (non-empty values override)
4. Env vars (`OPENAI_API_KEY`, `OPENAI_BASE_URL`) override everything

### Catalog Persistence
- JSON file at `clips/catalog.json`
- Atomic writes: write to `.tmp`, fsync, rename
- In-memory `[]Entry` with `AddEntry`, `FindEntry`, `Search`

### OpenAI Client
- HTTP client with 60s timeout
- Retry: up to 3 attempts with exponential backoff on 429/5xx
- TTS: POST `/audio/speech` (JSON body)
- STT: POST `/audio/transcriptions` (multipart form)

### Audio File Naming
Format: `<id>_<lang>_<voice>_<sampleRate>.<ext>`
Example: `7xnc52ubjm_en_alloy_24000.mp3`

## Adding a New Command
1. Add handler function in `cmd/coralie-clip/main.go`
2. Add case to `switch command` in `main()`
3. Implement `RunXxxCommand()` method on `App` in `internal/cli/cli.go`
4. Add tests
5. Update `docs/cli.md`
