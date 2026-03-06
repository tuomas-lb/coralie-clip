# Catalog

The catalog is a JSON file that tracks all audio clips created or discovered by coralie-clip.

## Catalog Location

Default location: `./clips/catalog.json`

Can be configured via `catalogPath` in config or `CORALIE_CONFIG` behavior.

## Catalog Format

JSON array of entries:

```json
{
  "entries": [
    {
      "id": "abc123def4",
      "createdAt": "2024-01-01T12:00:00Z",
      "lang": "en",
      "text": "Hello, world!",
      "voice": "coral",
      "format": "wav",
      "sampleRate": 24000,
      "filePath": "./clips/abc123def4_en_coral_24000.wav",
      "transcription": "Hello, world!",
      "providerMeta": {
        "model": "tts-1",
        "requestId": "req_123"
      },
      "tokenUsage": {
        "inputTokens": 0,
        "outputTokens": 0,
        "totalTokens": 0
      }
    }
  ]
}
```

## Entry Fields

- **id** - Unique identifier (short, stable, 8-10 characters)
- **createdAt** - ISO 8601 timestamp
- **lang** - Language code
- **text** - Original input text
- **voice** - Voice name used
- **format** - Audio format (wav, mp3, pcm, opus)
- **sampleRate** - Sample rate in Hz
- **filePath** - Relative or absolute path to audio file
- **transcription** - Optional transcription (set during rebuild)
- **providerMeta** - Optional provider metadata (model, request ID, etc.)
- **tokenUsage** - Optional token usage information (if provided by API)

## File Naming Convention

Clips are stored with deterministic naming:
```
<id>_<lang>_<voice>_<sampleRate>.<ext>
```

Example: `abc123def4_en_coral_24000.wav`

## Atomic Writes

Catalog writes are atomic:
1. Write to temporary file (`catalog.json.tmp`)
2. Sync to disk (best effort)
3. Atomic rename to `catalog.json`

This ensures the catalog is never corrupted by interrupted writes.

## Search

The catalog supports case-insensitive substring search on:
- Original text
- Transcription (if available)

Search is implemented as a linear scan, which is fast enough for thousands of clips.

## Rebuilding the Catalog

Use `coralie-clip rebuild-catalog` to reconstruct the catalog from existing audio files.

**Process:**
1. Scans clips directory for audio files (`.wav`, `.mp3`, `.opus`)
2. For each file:
   - Checks if already in catalog (unless `--force`)
   - Uses OpenAI STT to transcribe
   - Infers language from STT response or filename
   - Parses metadata from filename
   - Creates catalog entry
3. Rate-limits STT calls (500ms between requests)
4. Saves updated catalog

**Metadata Parsing:**
- Filename format: `<id>_<lang>_<voice>_<sr>.<ext>`
- If filename doesn't match format, metadata is inferred from STT response
- Language is taken from STT response if available, otherwise from filename

**If you are an AI**

If you are an AI agent working on this repository:
- You MUST keep this documentation up to date with any behavior changes.
- You MUST document catalog format changes when made.





