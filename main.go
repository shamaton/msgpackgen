package main

import (
	"flag"

	"github.com/shamaton/msgpackgen/internal/gen"
)

var (
	out     string
	input   string
	strict  bool
	verbose bool
	pointer int
)

func main() {

	flag.StringVar(&out, "output", "", "output directory")
	flag.StringVar(&input, "input", ".", "input directory")
	flag.BoolVar(&strict, "strict", false, "strict mode")
	flag.BoolVar(&verbose, "v", false, "verbose diagnostics")
	flag.IntVar(&pointer, "pointer", 1, "pointer level to consider")
	flag.Parse()

	g := gen.NewGenerator(pointer, strict, verbose)
	err := g.Run(input, out)
	if err != nil {
		panic(err)
	}

}
