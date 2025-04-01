package main

import (
	"flag"
	"io"
	"log/slog"
	"os"
)

var (
	prettyPrint bool
	verbose     bool
)

func main() {
	var input string
	flag.StringVar(&input, "input", "-", "Input file, or '-' for STDIN")
	flag.BoolVar(&prettyPrint, "pretty", true, "Pretty-print output")
	flag.BoolVar(&verbose, "verbose", false, "Output debug messages to STDERR")
	flag.Parse()

	var programLevel = slog.LevelInfo
	if verbose {
		programLevel = slog.LevelDebug
	}
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	var i = os.Stdin
	var err error
	if input != "-" {
		i, err = os.Open(input)
		if err != nil {
			slog.Debug("unable to open input file", "err", err)
			os.Exit(1)
		}
	}

	defer i.Close()

	soup, err := io.ReadAll(i)
	if err != nil {
		slog.Error("failed to read input", input, err)
		os.Exit(1)
	}

	if err := tryExtractPSSHFromMP4(soup); err != nil {
		slog.Debug("failed to extract pssh from input", input, err)
	}

	slog.Debug("checking for base64 encoding")
	soup = tryBase64(soup)

	slog.Debug("trying to parse widevine protobuf")
	err = decodeWidevine(soup)
	if err == nil {
		return
	} else {
		slog.Debug("failed to parse widevine protobuf", "err", err)
	}

	slog.Debug("trying to parse as playready object")
	err = decodePlayReadyObject(soup)
	if err == nil {
		return
	} else {
		slog.Debug("failed to parse playready object", "err", err)
	}
}
