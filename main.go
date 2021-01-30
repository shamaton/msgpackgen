package main

import (
	"flag"
	"log"

	"github.com/shamaton/msgpackgen/internal/generator"
)

var (
	input    = flag.String("i", ".", "input directory")
	output   = flag.String("o", ".", "output directory")
	filename = flag.String("g", defaultFileName, "generated file name")
	pointer  = flag.Int("p", defaultPointerLevel, "pointer level to consider")
	strict   = flag.Bool("s", false, "strict mode")
	verbose  = flag.Bool("v", false, "verbose diagnostics")
)

const (
	defaultFileName     = "resolver.msgpackgen.go"
	defaultPointerLevel = 1
)

func main() {

	flag.Parse()

	err := generator.Run(*input, *output, *filename, *pointer, *strict, *verbose)
	if err != nil {
		log.Fatal(err)
	}

}
