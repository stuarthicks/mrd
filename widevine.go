package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	widevine "github.com/iyear/gowidevine"
	wvpb "github.com/iyear/gowidevine/widevinepb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func decodeWidevine(b []byte) error {
	slog.Debug("trying to parse widevine protobuf")

	indent := ""
	if prettyPrint {
		indent = "  "
	}

	signedMsg := &wvpb.SignedMessage{}
	if err := proto.Unmarshal(tryBase64(b), signedMsg); err != nil {
		slog.Debug("failed to unmarshal signed widevinelicense message", "err", err.Error())
	} else {
		out, err := protojson.MarshalOptions{
			EmitUnpopulated: false,
			UseEnumNumbers:  false,
			Indent:          indent,
		}.Marshal(signedMsg)
		if err != nil {
			slog.Debug("failed to marshal signed widevine license message as json", "err", err.Error())
		} else {
			fmt.Fprintln(os.Stderr, "Widevine License Response:")
			fmt.Fprintln(os.Stderr, "----------------")
			fmt.Fprintln(os.Stdout, string(out))
			fmt.Println()

			licenseMsg := &wvpb.License{}
			if err = proto.Unmarshal(signedMsg.Msg, licenseMsg); err != nil {
				return fmt.Errorf("failed to unmarshal license message: %w", err)
			} else {
				out, err := protojson.MarshalOptions{
					EmitUnpopulated: false,
					UseEnumNumbers:  false,
					Indent:          indent,
				}.Marshal(licenseMsg)
				if err != nil {
					slog.Debug("failed to marshal license message as json", "err", err.Error())
				} else {
					fmt.Fprintln(os.Stderr, "Widevine License Message:")
					fmt.Fprintln(os.Stderr, "----------------")
					fmt.Fprintln(os.Stdout, string(out))

					return nil
				}
			}
		}
	}

	slog.Debug("trying to parse widevine protobuf")
	o, err := widevine.NewPSSH(b)
	if err != nil {
		return fmt.Errorf("failed to parse widevine pssh: %w", err)
	}

	var buf bytes.Buffer
	var j = json.NewEncoder(&buf)
	if err = j.Encode(o.Data()); err != nil {
		return fmt.Errorf("failed to encode widevine proto to json: %w", err)
	}

	var parsed = make(map[string]any)
	if err = json.NewDecoder(&buf).Decode(&parsed); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	var keyIDs = make([]string, 0)
	for _, k := range o.Data().KeyIds {
		u, err := uuid.FromBytes(k)
		if err == nil {
			keyIDs = append(keyIDs, u.String())
		}
	}
	parsed["key_ids"] = keyIDs

	var contentID = string(o.Data().ContentId)
	contentID = string(tryBase64([]byte(contentID)))

	var contentIDMap = make(map[string]any)
	err = json.NewDecoder(bytes.NewBufferString(contentID)).Decode(&contentIDMap)
	if err == nil {
		parsed["content_id"] = contentIDMap
	} else {
		parsed["content_id"] = contentID
	}

	scheme, ok := parsed["protection_scheme"].(float64)
	if ok {
		parsed["protection_scheme"] = protectionSchemeName(scheme)
	}

	var enc = json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	if prettyPrint {
		enc.SetIndent("", "  ")
	}

	fmt.Fprintln(os.Stderr, "Widevine Object:")
	fmt.Fprintln(os.Stderr, "----------------")

	if err = enc.Encode(parsed); err != nil {
		return fmt.Errorf("failed to write json to STDERR: %w", err)
	}

	return nil
}

// 'cenc' (AES-CTR) = 0x63656E63,
// 'cbc1' (AES-CBC) = 0x63626331,
// 'cens' (AES-CTR pattern encryption) = 0x63656E73,
// 'cbcs' (AES-CBC pattern encryption) = 0x63626373.
func protectionSchemeName(n float64) string {
	switch n {
	case 0x63656E63:
		return "cenc"
	case 0x63626331:
		return "cbc1"
	case 0x63656E73:
		return "cens"
	case 0x63626373:
		return "cbcs"
	}
	return "unknown"
}
