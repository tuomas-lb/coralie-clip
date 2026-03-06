# Go Policies

## Language & Tooling
- **Go version**: 1.22+ (per `go.mod`)
- **Module path**: `github.com/coralie/coralie-clip`
- **Binary**: `coralie-clip` (built to `./bin/`)
- **Dependencies**: minimal -- currently only `github.com/joho/godotenv`

## Package Layout Rules
| Package | Responsibility | Rules |
|---------|---------------|-------|
| `cmd/coralie-clip` | Entry point, argument parsing | No business logic; delegate to `internal/cli` |
| `internal/cli` | Command implementations | Orchestrates other packages via `App` struct |
| `internal/config` | Configuration | Load/save/validate; supports env, file, defaults |
| `internal/catalog` | Data persistence | JSON-based catalog; atomic writes |
| `internal/openai` | API client | HTTP only; no file I/O |
| `internal/audio` | Audio file operations | Format conversion, filename generation |
| `internal/player` | Audio playback | Platform detection, exec-based playback |
| `internal/logging` | Structured logging | JSON log lines, file-based |

## Coding Conventions
- All packages under `internal/` -- no public API surface
- Every `.go` file starts with `// Package <name> ...` doc comment
- Exported types and functions MUST have doc comments
- Error wrapping: `fmt.Errorf("context: %w", err)` always
- No naked returns
- No `init()` functions
- Use `context.Context` for operations that call external services

## Naming
| Item | Convention | Example |
|------|-----------|---------|
| Files | `snake_case.go` | `config.go`, `client.go` |
| Test files | `*_test.go` | `config_test.go` |
| Types | PascalCase | `Config`, `Entry`, `App` |
| Functions | PascalCase (exported), camelCase (unexported) | `LoadConfig`, `generateID` |
| Constants | PascalCase | `DefaultSampleRate` |
| Variables | camelCase | `catalogPath`, `httpClient` |

## Testing Rules
- Use standard `testing` package only (no testify, no gomock)
- Table-driven tests with descriptive `name` field
- `t.TempDir()` for file system tests
- Save/restore env vars in tests
- No external service calls in unit tests
- Test file in same package (white-box testing)

## Error Handling
- Functions return `error` as last return value
- Wrap with context before returning up
- `main()` calls `os.Exit(1)` on error -- nowhere else
- `ValidateConfigOrExit()` is the only function allowed to call `os.Exit` outside `main`

## Build Verification
```bash
make build && make test
```
Both must pass before any PR or merge.
