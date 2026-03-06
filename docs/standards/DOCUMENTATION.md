# Documentation Standards

## Principles
- Concise over verbose -- prefer tables, bullets, pseudocode
- Keep docs next to the code they describe when possible
- Every public function and type MUST have a Go doc comment
- Update docs in the same PR that changes behavior

## File Conventions
| Type | Location | Format |
|------|----------|--------|
| Package docs | `// Package X ...` in each `.go` file | Go doc comment |
| Architecture | `docs/design/system/` | Markdown |
| Feature specs | `docs/design/features/<name>/` | Markdown from template |
| API reference | Go doc comments | godoc-compatible |
| CLI reference | `docs/cli.md` | Markdown |
| Config reference | `docs/config.md` | Markdown |

## Go Doc Comments
- Every exported type, function, method, and constant must have a doc comment
- Start with the name of the thing being documented
- Use complete sentences

```go
// Config represents the application configuration.
// It supports loading from config.json, .env files, and environment variables.
type Config struct { ... }
```

## Markdown Standards
- Use ATX-style headers (`#`, `##`, etc.)
- One blank line before and after headers
- Use fenced code blocks with language identifier
- Tables for structured data
- Keep line length reasonable (no hard wrap required)

## What to Document
- **Always**: Public API, CLI commands, config options, architecture decisions
- **Never**: Implementation details that are obvious from code, auto-generated content
