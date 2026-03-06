# Go Security Reference

## Secrets
| Secret | Storage | Access |
|--------|---------|--------|
| OpenAI API key | `config.json` (0600), `.env`, or env var | `cfg.OpenAIApiKey` |
| Base URL | `config.json` or env var | `cfg.BaseURL` |

- `config.json` written with `os.WriteFile(path, data, 0600)`
- Both `config.json` and `.env` are in `.gitignore`
- `logging.RedactSecret()` available for safe logging

## HTTP Security
- Bearer token auth via `Authorization` header
- 60-second client timeout prevents hanging
- TLS enforced by default (OpenAI base URL is HTTPS)
- Retry logic excludes 4xx (client errors) except 429

## File System Security
| Operation | Safeguard |
|-----------|-----------|
| Config save | `0600` permissions |
| Catalog save | Atomic write (temp + rename) |
| Clips save | `0644` permissions |
| Directory creation | `os.MkdirAll` with `0755` |
| Path construction | `filepath.Join()` / `filepath.Abs()` |

## Input Validation
- CLI commands validated via switch statement (unknown = error + exit)
- Voice names validated against `SupportedVoices()` allowlist
- Audio formats validated against `SupportedFormats()` allowlist
- Language codes validated against `AvailableLanguages()` allowlist
- Sample rates validated against `SupportedSampleRates()` allowlist

## Process Execution
- `internal/player` uses `exec.Command` to invoke audio players
- Player command is either platform-detected or from `PLAYER_CMD` env var
- File paths resolved to absolute before exec to prevent path confusion

## Dependency Security
- Single external dependency (`godotenv`) -- small attack surface
- Use `go mod verify` to check integrity
- Monitor for CVEs in dependencies
