# Go Build Reference

## Build Commands
| Command | Description |
|---------|-------------|
| `make build` | Build binary to `./bin/coralie-clip` |
| `make clean` | Remove `./bin/` directory |
| `make install` | Build and install to `/usr/local/bin/` (or `$(PREFIX)/bin`) |
| `make uninstall` | Remove installed binary |

## Build Details
```makefile
BINARY_NAME := coralie-clip
CMD_DIR     := ./cmd/coralie-clip
BIN_DIR     := ./bin
GO          ?= go
```

Build command: `$(GO) build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)`

## Environment Variables
| Variable | Default | Purpose |
|----------|---------|---------|
| `GO` | `go` | Go binary to use |
| `PREFIX` | `/usr/local` | Install prefix |

## Run from Source
```bash
make run ARGS='fetch "Hello"'
# Equivalent to: go run ./cmd/coralie-clip -- fetch "Hello"
```

## Release Builds
- Git tags trigger GitHub Actions Docker builds
- Image: `ghcr.io/lastbotinc/coralie-cli:<tag>`
- Tags: `vX.Y.Z`, `X.Y`, `latest`

| Command | Action |
|---------|--------|
| `make version` | Show current version tag |
| `make bump-patch` | Tag + push vX.Y.Z+1 |
| `make bump-minor` | Tag + push vX.Y+1.0 |
| `make bump-major` | Tag + push vX+1.0.0 |
