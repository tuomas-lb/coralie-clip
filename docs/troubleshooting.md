# Troubleshooting

Common issues and solutions for coralie-clip.

## Configuration Issues

### "OpenAI API key is required"

**Problem:** API key is not configured.

**Solution:**
```bash
coralie-clip set openai-apikey <your-key>
# Or set OPENAI_API_KEY environment variable
export OPENAI_API_KEY=<your-key>
```

### "Config directory does not exist"

**Problem:** `CORALIE_CONFIG` points to a non-existent directory.

**Solution:**
```bash
mkdir -p /path/to/config/dir
# Or unset CORALIE_CONFIG to use default location
unset CORALIE_CONFIG
```

### "At least one language must be enabled"

**Problem:** No languages are enabled in configuration.

**Solution:**
```bash
coralie-clip lang enable en
# Or enable multiple languages
coralie-clip lang enable all
```

### "Language 'X' is not enabled"

**Problem:** Trying to use a language that isn't enabled.

**Solution:**
```bash
coralie-clip lang enable <code>
```

## Audio Playback Issues

### "No audio player found"

**Problem:** No audio player is available on the system.

**Solution:**

**macOS:**
- `afplay` should be available by default

**Linux:**
- Install a player:
  ```bash
  # Ubuntu/Debian
  sudo apt-get install pulseaudio-utils  # for paplay
  sudo apt-get install alsa-utils        # for aplay
  sudo apt-get install mpg123            # for mpg123
  ```

**Custom Player:**
```bash
export PLAYER_CMD=/path/to/your/player
```

### "Failed to play audio"

**Problem:** Player command failed.

**Solution:**
- Check that the audio file exists
- Verify player command works manually
- Check file permissions
- Try setting `PLAYER_CMD` to a different player

## API Issues

### "TTS request failed: API error (status 401)"

**Problem:** Invalid API key.

**Solution:**
- Verify API key is correct
- Check API key has not expired
- Ensure API key has TTS/STT permissions

### "TTS request failed: API error (status 429)"

**Problem:** Rate limit exceeded.

**Solution:**
- Wait before retrying
- The client automatically retries with backoff
- Consider upgrading your OpenAI plan

### "TTS request failed: API error (status 500)"

**Problem:** OpenAI API server error.

**Solution:**
- The client automatically retries
- Wait and try again
- Check OpenAI status page

## File Issues

### "Cannot create clips directory"

**Problem:** Insufficient permissions or invalid path.

**Solution:**
- Check directory permissions
- Ensure parent directory exists
- Use absolute path if relative path fails

### "Audio file not found"

**Problem:** Clip file was moved or deleted.

**Solution:**
- Rebuild catalog: `coralie-clip rebuild-catalog`
- Or manually remove entry from catalog

## Catalog Issues

### "Catalog already exists. Use --force to overwrite"

**Problem:** Trying to rebuild catalog that already has entries.

**Solution:**
```bash
coralie-clip rebuild-catalog --force
```

### "Failed to read catalog"

**Problem:** Catalog file is corrupted or invalid JSON.

**Solution:**
- Backup current catalog
- Try rebuilding: `coralie-clip rebuild-catalog --force`
- Or manually fix JSON syntax

## Build Issues

### "command not found: coralie-clip"

**Problem:** Binary not in PATH.

**Solution:**
```bash
# Check if installed
which coralie-clip

# If not in PATH, add it:
# For zsh:
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# For bash:
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### "go: command not found"

**Problem:** Go is not installed.

**Solution:**
- Install Go 1.22 or later
- Verify: `go version`

## Performance Issues

### Rebuild catalog is slow

**Problem:** Rebuild processes many files with rate limiting.

**Solution:**
- This is expected behavior (500ms between STT calls)
- Use `--force` only when necessary
- Consider processing in batches

## Getting Help

1. Check error messages - they often include actionable suggestions
2. Review configuration: `coralie-clip lang list all`
3. Check logs if `logOutputPath` is configured
4. Verify API key and permissions

## If you are an AI

If you are an AI agent working on this repository:
- You MUST keep this documentation up to date with any behavior changes.
- You MUST add new troubleshooting entries when common issues are discovered.





