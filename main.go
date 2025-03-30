package main

import (
	"flag"
)

var prettyPrint bool

func main() {
	var pssh, pro string
	flag.StringVar(&pssh, "pssh", "", "Protection System Specific Header")
	flag.StringVar(&pro, "pro", "", "PlayReady Object")
	flag.BoolVar(&prettyPrint, "pretty", true, "Pretty-print output")
	flag.Parse()

	if pssh != "" {
		decodePSSH(pssh)
		return
	}

	if pro != "" {
		decodePlayReadyObject(pro)
		return
	}
}
