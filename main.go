package main

import (
	"flag"
	"log"

	"github.com/shamaton/msgpackgen/internal/generator"
)

var (
	input    string
	output   string
	filename string
	strict   bool
	verbose  bool
	pointer  int
)

const (
	defaultFileName     = "resolver.msgpackgen.go"
	defaultPointerLevel = 1
)

//go:generate go run github.com/shamaton/msgpackgen -v

func main() {

	flag.StringVar(&input, "i", ".", "input directory")
	flag.StringVar(&output, "o", input, "output directory")
	flag.StringVar(&filename, "g", defaultFileName, "generated file name")
	flag.IntVar(&pointer, "p", defaultPointerLevel, "pointer level to consider")
	flag.BoolVar(&strict, "s", false, "strict mode")
	flag.BoolVar(&verbose, "v", false, "verbose diagnostics")
	flag.Parse()

	err := generator.Run(input, output, filename, pointer, strict, verbose)
	if err != nil {
		log.Fatal(err)
	}

}
