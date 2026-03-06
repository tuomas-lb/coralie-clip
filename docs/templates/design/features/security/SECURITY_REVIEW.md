# Security Review: [Feature Name]

## Review Date
YYYY-MM-DD

## Reviewer
Name / AI Agent

## Checklist

### Input Validation
- [ ] All user inputs validated
- [ ] No SQL/command injection vectors
- [ ] File paths sanitized with `filepath.Join()`

### Secrets
- [ ] No secrets logged or printed
- [ ] Secrets stored with appropriate file permissions (0600)
- [ ] Secrets not committed to version control

### Network
- [ ] TLS used for all external calls
- [ ] Timeouts set on HTTP clients
- [ ] No sensitive data in URLs or query parameters

### File System
- [ ] Appropriate file permissions set
- [ ] Atomic writes for critical data
- [ ] No path traversal vulnerabilities

### Process Execution
- [ ] No user-controlled command injection
- [ ] Arguments properly escaped
- [ ] File paths resolved to absolute before exec

### Error Handling
- [ ] Errors don't leak sensitive information
- [ ] All error paths handled
- [ ] Resources cleaned up on error

## Findings

| ID | Severity | Description | Status |
|----|----------|-------------|--------|
| S-1 | Low/Med/High | Finding | Open/Fixed |

## Conclusion
Summary of review outcome. Approved / Needs Changes / Blocked.
