# Implementation Process

## Overview
Step-by-step process for implementing features in coralie-clip.

## Steps

### 1. Set Up Task
```bash
cp task/templates/golang/IMPLEMENTATION.md task/active/<task-name>.md
```
Fill in the task file with scope, acceptance criteria, and checklist.

### 2. Create Branch
```bash
git checkout -b feature/<name>
```

### 3. Implement (Bottom-Up)
Order of implementation for this project:

| Order | Layer | Location | Notes |
|-------|-------|----------|-------|
| 1 | Data types | `internal/*/` | Structs, interfaces |
| 2 | Business logic | `internal/*/` | Core functions |
| 3 | Tests | `internal/*_test.go` | Test alongside code |
| 4 | CLI wiring | `cmd/coralie-clip/main.go` | Command routing |
| 5 | Config changes | `internal/config/config.go` | New fields, validation |
| 6 | Documentation | `docs/`, `README.md` | CLI docs, config docs |

### 4. Test
```bash
make test              # All tests
go test -v ./internal/<pkg>/  # Specific package
go test -race ./...    # Race detection
```

### 5. Build and Verify
```bash
make build
./bin/coralie-clip <new-command> ...
```

### 6. Update Documentation
- Update `docs/cli.md` if new commands added
- Update `docs/config.md` if new config options added
- Update `README.md` if user-facing behavior changed

### 7. Complete Task
```bash
mv task/active/<task-name>.md task/done/
```

## Commit Message Format
```
<type>: <short description>

<optional body>
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`

## Pre-Merge Checklist
- [ ] `make test` passes
- [ ] `make build` succeeds
- [ ] No new linting issues
- [ ] Documentation updated
- [ ] Task file moved to `task/done/`
