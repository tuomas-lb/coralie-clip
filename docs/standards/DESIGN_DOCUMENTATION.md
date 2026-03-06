# Design Documentation Standards

## When to Write a Design Doc
- New CLI command or subcommand
- New external integration (API provider, storage backend)
- Architectural change affecting multiple packages
- Security-sensitive feature
- Breaking change to config format or CLI interface

## Design Doc Structure
Use the templates in `docs/templates/design/features/`:

1. **FEATURE_SPEC.md** -- What and why
2. **REQUIREMENTS.md** -- Functional and non-functional requirements
3. **DESIGN.md** -- Technical design and implementation plan
4. **security/THREAT_ANALYSIS.md** -- Threat model (if security-relevant)
5. **security/SECURITY_REVIEW.md** -- Review checklist (if security-relevant)

## Where to Store
```
docs/design/features/<feature-name>/
  FEATURE_SPEC.md
  REQUIREMENTS.md
  DESIGN.md
  security/          # optional
    THREAT_ANALYSIS.md
    SECURITY_REVIEW.md
```

## Review Process
1. Draft design doc from templates
2. Self-review against requirements checklist
3. Peer review (if applicable)
4. Update design doc if implementation diverges

## Lightweight vs Full Design
| Scope | Documents Required |
|-------|-------------------|
| Bug fix | None |
| Small feature (< 1 file) | None |
| Medium feature (2-5 files) | FEATURE_SPEC + DESIGN |
| Large feature (> 5 files) | All documents |
| Security-sensitive | All documents + security/ |
