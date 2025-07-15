package main

import "encoding/base64"

func tryBase64(i []byte) []byte {
	// Some browsers add this prefix when you copy the response body from the devtools
	const prefix = "data:application/octet-stream;base64,"
	if len(i) >= len(prefix) && string(i[:len(prefix)]) == prefix {
		i = i[len(prefix):]
	}

	o, err := base64.StdEncoding.DecodeString(string(i))
	if err != nil {
		return i
	}
	return []byte(o)
}
