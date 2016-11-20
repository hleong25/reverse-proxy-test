package main

import (
	"flag"
)

type Flags struct {
	Plugin bool
	Port   int
}

var flags *Flags = &Flags{}

func init() {
	flag.BoolVar(&flags.Plugin, "plugin", false, "Start as plugin")
	flag.IntVar(&flags.Port, "port", 7000, "HTTP Server port")

	flag.Parse()
}
