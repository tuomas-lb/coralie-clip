# coralie-clip

A CLI tool for creating and managing audio clips using OpenAI TTS and STT.

## Quick Start

1. **Install the tool:**
   ```bash
   make install
   ```

2. **Set your OpenAI API key:**
   ```bash
   coralie-clip set openai-apikey <your-api-key>
   ```

3. **Enable a language:**
   ```bash
   coralie-clip lang enable en
   ```

4. **Create your first clip:**
   ```bash
   coralie-clip fetch "Hello, world!"
   ```

5. **Search for clips:**
   ```bash
   coralie-clip find "hello"
   ```

6. **Play a clip:**
   ```bash
   coralie-clip play <clip-id>
   ```

## Features

- Generate audio clips from text using OpenAI TTS
- Searchable catalog of all clips
- Rebuild catalog from existing audio files using OpenAI STT
- Cross-platform audio playback
- Configurable languages, voices, formats, and sample rates
- Request logging with token usage tracking

## Documentation

- [CLI Commands](cli.md) - Complete command reference
- [Configuration](config.md) - Configuration options and environment variables
- [Catalog](catalog.md) - Catalog format and rebuild process
- [OpenAI Integration](openai.md) - Provider notes and limitations
- [Troubleshooting](troubleshooting.md) - Common issues and solutions

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

## If you are an AI

If you are an AI agent working on this repository:
- You MUST keep this documentation up to date with any behavior changes.
- You MUST update the relevant documentation files when adding features or changing behavior.

