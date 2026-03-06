# Testing Standards

## Commands
| Action | Command |
|--------|---------|
| Run all tests | `make test` or `go test ./...` |
| Run specific package | `go test ./internal/config/` |
| Run specific test | `go test ./internal/config/ -run TestDefaultConfig` |
| Verbose output | `go test -v ./...` |
| With coverage | `go test -cover ./...` |
| Coverage report | `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out` |
| Race detection | `go test -race ./...` |

## Test File Conventions
- Test files: `*_test.go` in the same package as the code under test
- Test functions: `TestXxx(t *testing.T)`
- Table-driven tests preferred for multiple cases
- Use `t.TempDir()` for temporary files
- Use `t.Run()` for subtests

## Test Structure (Table-Driven)
```go
func TestFoo(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {name: "valid input", input: "hello", expected: "HELLO"},
        {name: "empty input", input: "", wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Foo(tt.input)
            if (err != nil) != tt.wantErr {
                t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Coverage Requirements
- New code should include tests
- Test both happy path and error cases
- External API calls (OpenAI) should be tested with mocks/interfaces

## What to Test
| Layer | Test Focus |
|-------|-----------|
| `internal/config` | Load/save, defaults, validation, env var override |
| `internal/catalog` | CRUD, search, persistence, atomic writes |
| `internal/audio` | Format conversion, filename generation |
| `internal/player` | Player detection, availability check |
| `internal/openai` | Request building, response parsing (with mocks) |
| `cmd/` | CLI argument parsing, command routing |

## Environment in Tests
- Save and restore environment variables (`defer os.Setenv(...)`)
- Use `t.TempDir()` for file system tests
- Never depend on external services in unit tests
