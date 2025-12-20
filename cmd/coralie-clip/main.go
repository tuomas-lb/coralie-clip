// Package main is the entry point for the coralie-clip CLI tool.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/coralie/coralie-clip/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "set":
		err = handleSetCommand(args)
	case "fetch":
		err = handleFetchCommand(args)
	case "lang":
		err = handleLangCommand(args)
	case "voice":
		err = handleVoiceCommand(args)
	case "format":
		err = handleFormatCommand(args)
	case "find":
		err = handleFindCommand(args)
	case "play":
		err = handlePlayCommand(args)
	case "rebuild-catalog":
		err = handleRebuildCatalogCommand(args)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `coralie-clip - Audio clip management tool

Usage:
  coralie-clip <command> [arguments]

Commands:
  set <key> <value>              Set a configuration value
  fetch "<text>" [options]       Generate audio clip from text
  lang enable <code|all>         Enable a language
  lang disable <code|all>        Disable a language
  lang list [all]                List enabled or all languages
  voice <name>                   Set default voice
  voice list                     List available voices
  format <wav|mp3|pcm|opus>      Set audio format
  find "<query>"                 Search for clips
  play <id>                      Play a clip by ID
  rebuild-catalog [--force]      Rebuild catalog from clips directory

Examples:
  coralie-clip set openai-apikey <key>
  coralie-clip fetch "Hello world" --lang en
  coralie-clip lang enable en
  coralie-clip find "hello"
  coralie-clip play abc123

`)
}

func handleSetCommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: coralie-clip set <key> <value>")
	}

	return cli.RunSetCommandStandalone(args[0], args[1])
}

func handleFetchCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: coralie-clip fetch \"<text>\" [--lang <code>] [--voice <name>] [--format <fmt>] [--sample-rate <hz>]")
	}

	app, err := cli.NewApp()
	if err != nil {
		return err
	}
	defer app.Close()

	// Parse text (first argument) - handle with or without quotes
	text := args[0]
	text = strings.Trim(text, "\"'")

	// Parse flags
	fs := flag.NewFlagSet("fetch", flag.ContinueOnError)
	lang := fs.String("lang", "", "Language code")
	voice := fs.String("voice", "", "Voice name")
	format := fs.String("format", "", "Audio format")
	sampleRateStr := fs.String("sample-rate", "", "Sample rate")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	sampleRate := 0
	if *sampleRateStr != "" {
		var err error
		sampleRate, err = strconv.Atoi(*sampleRateStr)
		if err != nil {
			return fmt.Errorf("invalid sample rate: %s", *sampleRateStr)
		}
	}

	return app.RunFetchCommand(text, *lang, *voice, *format, sampleRate)
}

func handleLangCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: coralie-clip lang <enable|disable|list> [code|all]")
	}

	app, err := cli.NewApp()
	if err != nil {
		return err
	}
	defer app.Close()

	action := args[0]
	code := ""
	if len(args) > 1 {
		code = args[1]
	} else if action == "list" {
		code = ""
	} else {
		return fmt.Errorf("usage: coralie-clip lang <enable|disable|list> [code|all]")
	}

	return app.RunLangCommand(action, code)
}

func handleVoiceCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: coralie-clip voice <list|<name>>")
	}

	app, err := cli.NewApp()
	if err != nil {
		return err
	}
	defer app.Close()

	action := args[0]
	if action == "list" {
		return app.RunVoiceCommand("list", "")
	}
	return app.RunVoiceCommand("set", action)
}

func handleFormatCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: coralie-clip format <wav|mp3|pcm|opus>")
	}

	app, err := cli.NewApp()
	if err != nil {
		return err
	}
	defer app.Close()

	return app.RunFormatCommand(args[0])
}

func handleFindCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: coralie-clip find \"<query>\"")
	}

	app, err := cli.NewApp()
	if err != nil {
		return err
	}
	defer app.Close()

	query := strings.Trim(args[0], "\"")
	return app.RunFindCommand(query)
}

func handlePlayCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: coralie-clip play <id>")
	}

	app, err := cli.NewApp()
	if err != nil {
		return err
	}
	defer app.Close()

	return app.RunPlayCommand(args[0])
}

func handleRebuildCatalogCommand(args []string) error {
	app, err := cli.NewApp()
	if err != nil {
		return err
	}
	defer app.Close()

	force := false
	for _, arg := range args {
		if arg == "--force" {
			force = true
		}
	}

	return app.RunRebuildCatalogCommand(force)
}

