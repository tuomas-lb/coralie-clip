# Secure Coding Standards

## Secrets Management
| Rule | Detail |
|------|--------|
| Never log secrets | Use `logging.RedactSecret()` for API keys |
| Config file permissions | `config.json` written with `0600` (owner-only read/write) |
| Environment variable precedence | `OPENAI_API_KEY` env var overrides config file |
| No secrets in source | `.env` and `config.json` are in `.gitignore` |

## Input Validation
- Validate all CLI arguments before use
- Validate config values against allowed lists (voices, formats, languages, sample rates)
- Use `ValidateConfigOrExit()` to enforce config constraints
- Trim and sanitize user-provided text before sending to API

## File Operations
| Practice | Implementation |
|----------|---------------|
| Atomic writes | Catalog uses write-to-temp + rename pattern |
| Directory creation | `os.MkdirAll` with `0755` for directories |
| File permissions | `0644` for data files, `0600` for config with secrets |
| Path validation | Use `filepath.Join()` and `filepath.Abs()` -- never string concat |

## HTTP Client Security
- Set explicit timeouts (60s) on HTTP clients
- Use `context.Context` for request cancellation
- Retry with exponential backoff on 429/5xx only
- Do not retry on 4xx (except 429)
- Validate response status before reading body

## Error Handling
- Wrap errors with `fmt.Errorf("context: %w", err)` for stack context
- Never expose internal errors to stdout -- use stderr
- Log errors with structured logging, not raw prints
- Exit with non-zero status on fatal errors

## Dependencies
- Minimize external dependencies (currently only `github.com/joho/godotenv`)
- Vet new dependencies for maintenance status and security
- Use `go mod tidy` to remove unused dependencies
