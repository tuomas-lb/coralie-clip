# Generic Policies

## Project: coralie-clip
A CLI tool for creating and managing audio clips using OpenAI TTS and STT.

## Build & Run
| Action | Command |
|--------|---------|
| Build | `make build` |
| Run from source | `make run ARGS='<command> <args>'` |
| Test | `make test` |
| Install | `make install` |
| Clean | `make clean` |

## Project Structure
```
cmd/coralie-clip/main.go   -- Entry point, CLI routing
internal/
  cli/cli.go               -- App struct, command implementations
  config/config.go          -- Config load/save/validate
  catalog/catalog.go        -- Clip catalog CRUD and search
  openai/client.go          -- OpenAI TTS/STT API client
  audio/audio.go            -- Audio file ops, format conversion
  player/player.go          -- Cross-platform audio playback
  logging/logging.go        -- Structured request logging
```

## Mandatory Rules
1. Run `make test` before considering any task complete
2. Run `make build` to verify compilation
3. Never commit secrets (`config.json`, `.env`, API keys)
4. Update documentation when changing CLI commands or config options
5. Use `internal/` for all packages -- nothing goes in `pkg/`
6. Wrap errors with context: `fmt.Errorf("action: %w", err)`
7. Write to stderr for errors, stdout for normal output
8. Use table-driven tests with `t.Run()` subtests

## Configuration Priority (highest to lowest)
1. Environment variables (`OPENAI_API_KEY`, `OPENAI_BASE_URL`)
2. Config file (`config.json`)
3. `.env` file
4. Defaults (`internal/config/DefaultConfig()`)

## Version Management
| Action | Command |
|--------|---------|
| Show version | `make version` |
| Bump patch | `make bump-patch` |
| Bump minor | `make bump-minor` |
| Bump major | `make bump-major` |

## Deployment
- **Staging**: auto-deploys on push to `main` (Docker image `latest`)
- **Production**: deploy via semver tag (`make bump-*`)
- Docker image: `ghcr.io/lastbotinc/coralie-cli:<tag>`
