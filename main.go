package main

import (
	"flag"
)

var prettyPrint bool

func main() {
	var pro string
	flag.StringVar(&pro, "pro", "", "PlayReady Object")
	flag.BoolVar(&prettyPrint, "pretty", true, "Pretty-print output")
	flag.Parse()

	if pro != "" {
		decodePlayReady(pro)
		return
	}
}
