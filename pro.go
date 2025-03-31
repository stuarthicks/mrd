package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/go-restruct/restruct"
	"github.com/go-xmlfmt/xmlfmt"
	"golang.org/x/text/encoding/unicode"
)

type PlayReadyObject struct {
	Length      int                   `struct:"int32"`
	RecordCount int                   `struct:"int16"`
	Records     PlayReadyObjectRecord // FIXME: restruct doesn't seem to correctly parse slices
}

type PlayReadyObjectRecord struct {
	Type   int    `struct:"int16"`
	Length int    `struct:"int16,sizeof=Value"`
	Value  []byte `struct:"[]byte,sizefrom=Length,lsb"`
}

func decodePlayReadyObject(b []byte) error {
	var o PlayReadyObject

	if err := restruct.Unpack(b, binary.LittleEndian, &o); err != nil {
		return fmt.Errorf("failed to parse pro: %w", err)
	}

	var dec = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	prh, err := dec.Bytes(o.Records.Value)
	if err != nil {
		return fmt.Errorf("failed to decode UTF-16LE: %w", err)
	}

	fmt.Fprintln(os.Stderr, "PlayReady Object:")
	fmt.Fprintln(os.Stderr, "-----------------")

	if prettyPrint {
		fmt.Println(strings.TrimSpace(xmlfmt.FormatXML(string(prh), "", "  ")))
		return nil
	}

	fmt.Println(string(prh))
	return nil
}
