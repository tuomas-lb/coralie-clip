# Go Stack Reference -- coralie-clip

## Quick Reference
| Item | Value |
|------|-------|
| Go version | 1.22+ |
| Module | `github.com/coralie/coralie-clip` |
| Binary | `coralie-clip` |
| Entry point | `cmd/coralie-clip/main.go` |
| Dependencies | `github.com/joho/godotenv` v1.5.1 |

## Stack Documents
| Document | Content |
|----------|---------|
| [BUILD.md](BUILD.md) | Build targets, binary output, installation |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Dev workflow, code organization, patterns |
| [TESTING.md](TESTING.md) | Test commands, patterns, coverage |
| [SECURITY.md](SECURITY.md) | Security practices for Go code |

## Package Map
```
cmd/coralie-clip/    --> CLI entry point
internal/cli/        --> Command orchestration
internal/config/     --> Configuration management
internal/catalog/    --> Clip catalog (JSON persistence)
internal/openai/     --> OpenAI API client (TTS + STT)
internal/audio/      --> Audio file I/O and conversion
internal/player/     --> Cross-platform audio playback
internal/logging/    --> Structured JSON logging
```
