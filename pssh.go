package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
	widevine "github.com/iyear/gowidevine"
)

func decodePSSH(pssh string) {
	var i = os.Stdin
	var err error
	if pssh != "-" {
		i, err = os.Open(pssh)
		if err != nil {
			log.Fatal("ERR_1 ", err.Error())
		}
	}

	defer i.Close()

	data, err := io.ReadAll(i)
	if err != nil {
		log.Fatal("ERR_2 ", err.Error())
	}

	bb, err := base64.StdEncoding.DecodeString(string(data))
	if err == nil {
		data = bb
	}

	o, err := widevine.NewPSSH(data)
	if err != nil {
		log.Fatal("ERR_3 ", err.Error())
	}

	var buf bytes.Buffer
	var j = json.NewEncoder(&buf)
	if err = j.Encode(o.Data()); err != nil {
		log.Fatal("ERR_4 ", err.Error())
	}

	var parsed = make(map[string]any)
	if err = json.NewDecoder(&buf).Decode(&parsed); err != nil {
		log.Fatal("ERR_5 ", err.Error())
	}

	var keyIDs = make([]string, 0)
	for _, k := range o.Data().KeyIds {
		u, err := uuid.FromBytes(k)
		if err == nil {
			keyIDs = append(keyIDs, u.String())
		}
	}

	var contentID = string(o.Data().ContentId)

	parsed["key_ids"] = keyIDs
	parsed["content_id"] = contentID

	var enc = json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	if prettyPrint {
		enc.SetIndent("", "  ")
	}

	if err = enc.Encode(parsed); err != nil {
		log.Fatal("ERR_6 ", err.Error())
	}
}
