# Go Testing Reference

## Commands
| Action | Command |
|--------|---------|
| All tests | `make test` |
| Verbose | `go test -v ./...` |
| Single package | `go test ./internal/catalog/` |
| Single test | `go test ./internal/catalog/ -run TestSearch` |
| Race detector | `go test -race ./...` |
| Coverage | `go test -cover ./...` |
| Coverage HTML | `go test -coverprofile=c.out ./... && go tool cover -html=c.out` |

## Existing Test Files
| File | Tests |
|------|-------|
| `internal/config/config_test.go` | ResolveConfigPath, DefaultConfig, ValidateConfigOrExit, SaveAndLoadConfig |
| `internal/catalog/catalog_test.go` | LoadAndSaveCatalog, AddEntry, FindEntry, Search |
| `internal/player/player_test.go` | NewPlayer, PlayerIsAvailable |

## Test Patterns Used

### Table-Driven Tests
All test files use table-driven patterns with `t.Run()` subtests.

### Temp Directories
```go
tmpDir := t.TempDir()
configPath := filepath.Join(tmpDir, "config.json")
```

### Environment Variable Isolation
```go
originalEnv := os.Getenv("CORALIE_CONFIG")
defer os.Setenv("CORALIE_CONFIG", originalEnv)
os.Setenv("CORALIE_CONFIG", configPath)
```

## Testing Gaps (Opportunities)
| Package | Missing Coverage |
|---------|-----------------|
| `internal/openai` | No tests -- needs HTTP mock |
| `internal/audio` | No tests for PCMToWAV, GenerateFileName |
| `internal/logging` | No tests |
| `internal/cli` | No integration tests |
| `cmd/coralie-clip` | No CLI argument tests |

## Writing a New Test
1. Create `*_test.go` in the same package
2. Use table-driven format (see `docs/standards/TESTING.md`)
3. For HTTP-dependent tests, use `httptest.NewServer`
4. Run `go test -v ./internal/<pkg>/` to verify
