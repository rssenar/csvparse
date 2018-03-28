package main

import (
	"flag"
)

var (
	ij = flag.Bool("ij", false, "Enable output to Indented JSON")
	j  = flag.Bool("j", false, "Enable output to JSON")
	v  = flag.Bool("v", false, "Enable header field validation")
)
