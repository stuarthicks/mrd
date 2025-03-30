package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

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

func decodePlayReadyObject(pro string) {
	var o PlayReadyObject

	var i = os.Stdin
	var err error
	if pro != "-" {
		i, err = os.Open(pro)
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

	if err = restruct.Unpack(data, binary.LittleEndian, &o); err != nil {
		log.Fatal("ERR_3 ", err.Error())
	}

	var dec = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	prh, err := dec.Bytes(o.Records.Value)
	if err != nil {
		log.Fatal("ERR_4 ", err.Error())
	}

	if prettyPrint {
		fmt.Println(xmlfmt.FormatXML(string(prh), "", "  "))
		return
	}

	fmt.Println(string(prh))
}
