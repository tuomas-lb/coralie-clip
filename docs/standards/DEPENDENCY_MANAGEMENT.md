# Dependency Management

## Current Dependencies
| Module | Version | Purpose |
|--------|---------|---------|
| `github.com/joho/godotenv` | v1.5.1 | Load `.env` files |

## Go Version
- **Minimum**: Go 1.22 (set in `go.mod`)

## Commands
| Action | Command |
|--------|---------|
| Add dependency | `go get <module>@<version>` |
| Update dependency | `go get -u <module>` |
| Update all | `go get -u ./...` |
| Remove unused | `go mod tidy` |
| Verify checksums | `go mod verify` |
| List dependencies | `go list -m all` |
| Check for updates | `go list -m -u all` |

## Rules
1. **Minimize dependencies** -- prefer stdlib over third-party when quality is comparable
2. **Pin versions** -- always use explicit versions, never `latest`
3. **Verify before adding** -- check maintenance status, license, security history
4. **Run `go mod tidy`** after any dependency change
5. **Commit `go.sum`** alongside `go.mod` changes
6. **No `vendor/`** -- `vendor/` directory is gitignored; use module proxy

## Adding a New Dependency Checklist
- [ ] Is there a stdlib alternative?
- [ ] Is the module actively maintained?
- [ ] Compatible license?
- [ ] Reasonable dependency tree (check transitive deps)?
- [ ] Run `go mod tidy` after adding
- [ ] Run `go test ./...` to verify compatibility
