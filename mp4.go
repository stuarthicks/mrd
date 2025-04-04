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
	var box mp4.Box
	var err error

	m, err := mp4.DecodeFile(bytes.NewReader(b))
	if err == nil {
		for _, pssh := range m.Moov.Psshs {
			inspectPSSHBox(pssh)
		}
		os.Exit(0)
		return nil
	}
	slog.Debug("unable to parse input as mp4", "err", err.Error())

	b = tryBase64(b)

	box, err = mp4.DecodeBox(0, bytes.NewReader(b))
	if err != nil {
		slog.Debug("unable to parse input as mp4 box", "err", err.Error())
		return errors.New("unable to parse input as mp4 or mp4 box")
	}

	if psshBox, ok := box.(*mp4.PsshBox); ok {
		inspectPSSHBox(psshBox)
	} else {
		slog.Debug("box is not a PSSH box")
	}

	os.Exit(0)
	return nil
}

func inspectPSSHBox(pssh *mp4.PsshBox) {
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
			return
		}
		if err := decodeWidevine(buf.Bytes()); err != nil {
			slog.Debug("failed to decode widevine pssh", "err", err)
			return
		}
		fmt.Fprintln(os.Stderr, "")
	default:
		slog.Debug("skipping unknown pssh", "system_id", pssh.SystemID.String(), "version", int(pssh.Version))
		return
	}
}
