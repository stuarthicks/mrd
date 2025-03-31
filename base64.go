package main

import "encoding/base64"

func tryBase64(i []byte) []byte {
	o, err := base64.StdEncoding.DecodeString(string(i))
	if err != nil {
		return i
	}
	return []byte(o)
}
