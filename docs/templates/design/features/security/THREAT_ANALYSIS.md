# Threat Analysis: [Feature Name]

## Assets
| Asset | Sensitivity | Location |
|-------|------------|----------|
| OpenAI API key | High | config.json, .env, env var |
| Audio clips | Low | clips/ directory |
| Catalog | Low | clips/catalog.json |

## Threat Model

### T1: [Threat Name]
- **Attack vector**: How the threat is exploited
- **Impact**: What damage could occur
- **Likelihood**: Low / Medium / High
- **Mitigation**: How it is addressed
- **Status**: Mitigated / Accepted / Open

### T2: [Threat Name]
- **Attack vector**:
- **Impact**:
- **Likelihood**:
- **Mitigation**:
- **Status**:

## Data Flow Security
```
User Input → Validation → API Call (HTTPS) → File Write
     ↑ sanitized         ↑ auth header       ↑ permissions
```

## Trust Boundaries
| Boundary | Between | Controls |
|----------|---------|----------|
| CLI input | User → Application | Input validation |
| API call | Application → OpenAI | TLS, Bearer auth |
| File system | Application → Disk | File permissions |
| Exec | Application → OS player | Path resolution |
