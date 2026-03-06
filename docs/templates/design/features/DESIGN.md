# Design: [Feature Name]

## Overview
Brief technical overview of the solution.

## Package Changes

### `internal/<package>`
```go
// New or modified types/functions
type NewType struct { ... }
func NewFunction(args) (result, error) { ... }
```

### `cmd/coralie-clip`
- New command handler: `handleXxxCommand()`
- New switch case in `main()`

### `internal/config`
- New config fields (if any)
- Validation changes

## Data Flow
```
1. Step 1
2. Step 2
3. Step 3
```

## Error Handling
| Error Case | Handling |
|-----------|----------|
| Invalid input | Return descriptive error |
| API failure | Retry with backoff (3 attempts) |
| File I/O error | Wrap and return |

## Config Changes
| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `newField` | string | `""` | What it does |

## Migration
- Is this backward compatible? Yes/No
- Migration steps if breaking

## Alternatives Considered
| Alternative | Pros | Cons | Decision |
|------------|------|------|----------|
| Option A | ... | ... | Chosen / Rejected |
