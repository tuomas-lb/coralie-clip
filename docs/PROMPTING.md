# AI Prompting Conventions

## Context Loading
When starting work on this project, the AI agent should:
1. Read `CLAUDE.md` (auto-generated project policies)
2. Read `docs/design/system/ARCHITECTURE.md` for system overview
3. Read relevant `internal/` source files for the task at hand

## Task Prompting Pattern
```
Task: <what to do>
Context: <relevant files or packages>
Constraints: <what NOT to do>
Output: <expected deliverable>
```

## Rules for AI Agents
1. **Read before writing** -- always read existing code before modifying
2. **Run tests** -- execute `make test` after changes
3. **Run build** -- execute `make build` to verify compilation
4. **Follow conventions** -- match existing code style in the file being edited
5. **Update docs** -- update `docs/cli.md`, `docs/config.md`, `README.md` when changing behavior
6. **No new dependencies** without explicit approval
7. **Use `internal/`** -- never create a `pkg/` directory
8. **Wrap errors** -- always add context with `fmt.Errorf("...: %w", err)`
9. **Table-driven tests** -- use the pattern from existing test files
10. **Atomic operations** -- use temp file + rename for critical writes

## Common Tasks

### Add a New CLI Command
```
Read: cmd/coralie-clip/main.go, internal/cli/cli.go
Create: handler in main.go, RunXxxCommand in cli.go
Test: add tests in appropriate _test.go files
Update: docs/cli.md, README.md
Verify: make build && make test
```

### Add a Config Option
```
Read: internal/config/config.go
Modify: Config struct, DefaultConfig(), LoadConfig(), ValidateConfigOrExit()
Test: internal/config/config_test.go
Update: docs/config.md
Verify: make build && make test
```

### Fix a Bug
```
Read: relevant source file
Write: fix + test that reproduces the bug
Verify: make test (new test passes, existing tests still pass)
```
