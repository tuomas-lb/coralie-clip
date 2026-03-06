# Implementation Task: [Title]

## Status
Not Started | In Progress | Done

## Date
YYYY-MM-DD

## Summary
What is being implemented and why.

## Design Reference
Link to design doc if applicable: `docs/design/features/<name>/`

## Scope
| Package | Changes |
|---------|---------|
| `internal/<pkg>` | Description |
| `cmd/coralie-clip` | Description |

## Acceptance Criteria
- [ ] Feature works as specified
- [ ] `make test` passes
- [ ] `make build` succeeds
- [ ] Tests added for new code
- [ ] Documentation updated

## Implementation Checklist
- [ ] Create/modify types and interfaces
- [ ] Implement business logic
- [ ] Write tests (table-driven, with subtests)
- [ ] Wire into CLI (`cmd/coralie-clip/main.go`)
- [ ] Update config if needed (`internal/config/config.go`)
- [ ] Update `docs/cli.md`
- [ ] Update `docs/config.md` (if config changes)
- [ ] Update `README.md` (if user-facing changes)
- [ ] Run `make build && make test`

## Test Plan
| Test | Package | Covers |
|------|---------|--------|
| TestXxx | `internal/<pkg>` | Happy path |
| TestXxx_Error | `internal/<pkg>` | Error case |

## Notes
