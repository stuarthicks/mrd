package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Eyevinn/mp4ff/mp4"
)

func tryExtractPSSHFromMP4(b []byte) error {
	var buf = bytes.NewReader(b)
	m, err := mp4.DecodeFile(buf)
	if err != nil {
		slog.Debug("unable to parse input as mp4", "err", err.Error())
		return errors.New("unable to parse input as mp4")
	}

	for _, pssh := range m.Moov.Psshs {
		switch strings.ToLower(pssh.SystemID.String()) {
		case mp4.UUIDPlayReady:
			slog.Debug("detected playready", "system_id", pssh.SystemID.String(), "version", int(pssh.Version))
			if err := decodePlayReadyObject(pssh.Data); err != nil {
				slog.Debug("failed to decode playready pssh", "err", err)
			}
			fmt.Fprintln(os.Stderr, "")
		case mp4.UUIDWidevine:
			slog.Debug("detected widevine", "system_id", pssh.SystemID.String(), "version", int(pssh.Version))
			var buf bytes.Buffer
			if err := pssh.Encode(&buf); err != nil {
				slog.Debug("failed to extract widevine pssh data", "err", err)
				continue
			}
			if err := decodeWidevine(buf.Bytes()); err != nil {
				slog.Debug("failed to decode widevine pssh", "err", err)
				continue
			}
			fmt.Fprintln(os.Stderr, "")
		default:
			slog.Debug("skipping unknown pssh", "system_id", pssh.SystemID.String(), "version", int(pssh.Version))
			continue
		}

	}

	os.Exit(0)
	return nil
}
