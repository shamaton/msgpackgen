package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/shamaton/msgpackgen/internal/generator"
)

var (
	inputDir  = flag.String("input-dir", ".", "input directory. input-file cannot be used at the same time")
	inputFile = flag.String("input-file", "", "input a specific file. input-dir cannot be used at the same time")
	outputDir = flag.String("output-dir", ".", "output directory")
	filename  = flag.String("output-file", defaultFileName, "name of generated file")
	pointer   = flag.Int("pointer", defaultPointerLevel, "pointer level to consider")
	useGopath = flag.Bool("use-gopath", false, "use GOPATH instead of go.mod")
	dryRun    = flag.Bool("dry-run", false, "dry run mode")
	strict    = flag.Bool("strict", false, "strict mode")
	verbose   = flag.Bool("v", false, "verbose diagnostics")
)

const (
	defaultFileName     = "resolver.msgpackgen.go"
	defaultPointerLevel = 1
)

func main() {
	flag.Parse()
	err := generate(
		*inputDir,
		*inputFile,
		*outputDir,
		*filename,
		*pointer,
		*useGopath,
		*dryRun,
		*strict,
		*verbose,
		os.Stdout,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func generate(iDir, iFile, oDir, oFile string, p int, ugo, dry, s, v bool, w io.Writer) error {
	return generator.Run(iDir, iFile, oDir, oFile, p, ugo, dry, s, v, w)
}
