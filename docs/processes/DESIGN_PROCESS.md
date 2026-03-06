# Design Process

## Overview
This process applies when adding features that touch multiple packages or change external behavior.

## Steps

### 1. Understand the Problem
- Read existing docs: `docs/cli.md`, `docs/config.md`, `docs/catalog.md`
- Read relevant source in `internal/` and `cmd/`
- Identify affected packages

### 2. Create Feature Directory
```bash
mkdir -p docs/design/features/<feature-name>/security
```

### 3. Write Specification
- Copy `docs/templates/design/features/FEATURE_SPEC.md`
- Define scope, goals, non-goals
- List affected CLI commands and config options

### 4. Define Requirements
- Copy `docs/templates/design/features/REQUIREMENTS.md`
- Functional requirements (what it does)
- Non-functional requirements (performance, security, compatibility)

### 5. Write Design
- Copy `docs/templates/design/features/DESIGN.md`
- Package-level changes
- Data flow
- Error handling strategy
- Config changes

### 6. Security Analysis (if applicable)
- Copy threat analysis and security review templates
- Identify secrets, external calls, file I/O
- Document mitigations

### 7. Review
- Self-check against requirements
- Verify backward compatibility of config format
- Verify CLI interface consistency

## Checklist Before Implementation
- [ ] Feature spec written and reviewed
- [ ] Requirements defined
- [ ] Design covers all affected packages
- [ ] Security analysis complete (if needed)
- [ ] No breaking changes to CLI or config (or migration plan exists)
