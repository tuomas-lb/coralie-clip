# coralie-clip

A CLI tool for creating and managing audio clips using OpenAI TTS and STT.

## Quick Start

1. **Build the tool:**
   ```bash
   make build
   ```

2. **Set your OpenAI API key:**
   ```bash
   ./bin/coralie-clip set openai-apikey <your-api-key>
   ```

3. **Enable a language:**
   ```bash
   ./bin/coralie-clip lang enable en
   ```

4. **Create your first clip:**
   ```bash
   ./bin/coralie-clip fetch "Hello, world!"
   ```

5. **Search for clips:**
   ```bash
   ./bin/coralie-clip find "hello"
   ```

6. **Play a clip:**
   ```bash
   ./bin/coralie-clip play <clip-id>
   ```

## Installation

```bash
make install
```

This installs the binary to `/usr/local/bin/coralie-clip` (or `$(PREFIX)/bin` if `PREFIX` is set).

## Features

- Generate audio clips from text using OpenAI TTS
- Searchable catalog of all clips
- Rebuild catalog from existing audio files using OpenAI STT
- Cross-platform audio playback
- Configurable languages, voices, formats, and sample rates
- Request logging with token usage tracking

## Documentation

See the [docs/](docs/) directory for complete documentation:

- [README.md](docs/README.md) - Project overview and quickstart
- [CLI Commands](docs/cli.md) - Complete command reference
- [Configuration](docs/config.md) - Configuration options and environment variables
- [Catalog](docs/catalog.md) - Catalog format and rebuild process
- [OpenAI Integration](docs/openai.md) - Provider notes and limitations
- [Troubleshooting](docs/troubleshooting.md) - Common issues and solutions

## Building

```bash
make build
```

The binary will be created in `./bin/coralie-clip`.

## Running from Source

```bash
make run ARGS='fetch "Hello"'
```

## Testing

```bash
make test
```

## Makefile Targets

- `make build` - Build the binary
- `make run ARGS="..."` - Run from source with arguments
- `make test` - Run unit tests
- `make install` - Install binary to system
- `make uninstall` - Remove installed binary
- `make clean` - Remove build artifacts

## Requirements

- Go 1.22 or later
- OpenAI API key
- Audio player (afplay on macOS, paplay/aplay on Linux)

## License

See LICENSE file for details.

