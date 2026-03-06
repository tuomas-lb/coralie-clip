# Project Setup

## Prerequisites
| Requirement | Version | Check Command |
|------------|---------|---------------|
| Go | 1.22+ | `go version` |
| Git | any | `git --version` |
| Make | any | `make --version` |
| Audio player | any | `which afplay` (macOS) or `which paplay` (Linux) |

## Quick Setup
```bash
# Clone
git clone <repo-url>
cd coralie-cli

# Build
make build

# Configure
./bin/coralie-clip set openai-apikey <your-key>
./bin/coralie-clip lang enable en

# Verify
make test
./bin/coralie-clip fetch "Hello, world!"
```

## Configuration
Configuration is loaded from multiple sources (highest priority first):

| Source | Location | Notes |
|--------|----------|-------|
| Environment vars | `OPENAI_API_KEY`, `OPENAI_BASE_URL` | Highest priority |
| Config file | `./config.json` (or `$CORALIE_CONFIG`) | Set via `coralie-clip set` |
| .env file | `./.env` | Loaded via godotenv |
| Defaults | `internal/config/DefaultConfig()` | Lowest priority |

## Environment Variables
| Variable | Purpose |
|----------|---------|
| `OPENAI_API_KEY` | OpenAI API key (overrides config) |
| `OPENAI_BASE_URL` | Custom API base URL |
| `CORALIE_CONFIG` | Path to config file or directory |
| `PLAYER_CMD` | Custom audio player command |

## Generate CLAUDE.md
```bash
bash scripts/setup-stacks.sh golang
```
This concatenates generic policies + Go policies into `CLAUDE.md` for AI agent context.

## Install System-Wide
```bash
make install               # Installs to /usr/local/bin/
make install PREFIX=~/bin   # Custom prefix
```
