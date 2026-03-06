# System Architecture -- coralie-clip

## Overview
coralie-clip is a CLI tool for creating and managing audio clips via OpenAI TTS/STT APIs. It follows a layered architecture with all packages under `internal/`.

## Architecture Diagram
```
┌─────────────────────────────────────────┐
│  cmd/coralie-clip/main.go               │
│  CLI entry point, argument parsing,     │
│  command dispatch                       │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│  internal/cli                            │
│  App struct -- orchestrates all packages │
│  Command handlers (Fetch, Find, Play...) │
└──┬───────┬───────┬───────┬───────┬──────┘
   │       │       │       │       │
┌──▼──┐ ┌──▼──┐ ┌──▼──┐ ┌──▼──┐ ┌──▼──┐
│config│ │cata-│ │open-│ │audio│ │play-│
│     │ │log  │ │ai   │ │     │ │er   │
└─────┘ └─────┘ └──┬──┘ └─────┘ └──┬──┘
                    │               │
              ┌─────▼─────┐   ┌────▼────┐
              │ OpenAI API│   │ OS exec │
              │ (HTTPS)   │   │ (afplay)│
              └───────────┘   └─────────┘
```

## Package Responsibilities

| Package | Depends On | Responsibility |
|---------|-----------|----------------|
| `cmd/coralie-clip` | `internal/cli` | Parse CLI args, dispatch to App |
| `internal/cli` | config, catalog, openai, audio, player, logging | Orchestrate commands |
| `internal/config` | godotenv | Load/save/validate config from multiple sources |
| `internal/catalog` | (stdlib) | JSON catalog with CRUD, search, atomic persistence |
| `internal/openai` | (stdlib) | HTTP client for TTS and STT endpoints |
| `internal/audio` | (stdlib) | WAV header construction, filename generation |
| `internal/player` | (stdlib) | Detect and exec platform audio player |
| `internal/logging` | (stdlib) | Append JSON log lines to file |

## Data Flow: `fetch` Command
```
1. main.go parses args → handleFetchCommand()
2. cli.NewApp() loads config + catalog + creates client
3. app.RunFetchCommand(text, lang, voice, format, sampleRate)
4. openai.Client.TTS() → POST /audio/speech → audio bytes
5. audio.SaveAudio() → write file to clips/
6. catalog.AddEntry() + catalog.SaveCatalog() → update catalog.json
7. Print clip ID and path
```

## Data Flow: `rebuild-catalog` Command
```
1. Scan clips directory for audio files
2. For each file: openai.Client.STT() → transcribe
3. Parse filename for metadata (id, lang, voice, sample rate)
4. Build new catalog entries
5. Save catalog atomically
```

## External Dependencies
| Dependency | Protocol | Purpose |
|-----------|----------|---------|
| OpenAI API | HTTPS | TTS (text-to-speech) and STT (speech-to-text) |
| OS audio player | exec | Playback via afplay/paplay/aplay |

## Storage
| Data | Format | Location |
|------|--------|----------|
| Config | JSON | `./config.json` (or `$CORALIE_CONFIG`) |
| Secrets | dotenv | `./.env` |
| Clips | wav/mp3/opus | `./clips/` |
| Catalog | JSON | `./clips/catalog.json` |
| Request log | JSON lines | Configurable path |
