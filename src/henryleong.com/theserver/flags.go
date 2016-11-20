package main

import (
	"flag"
)

type Flags struct {
	Port int
}

var flags *Flags = &Flags{}

func init() {
	flag.IntVar(&flags.Port, "port", 7000, "HTTP Server port")

	flag.Parse()
}
