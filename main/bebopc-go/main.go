package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/200sc/bebop"
)

var inputFile = flag.String("i", "", "the name of the file to compile")
var outputFile = flag.String("o", "", "the name of the output file to write")
var printVersion = flag.Bool("version", false, "print the version of the compiler")
var printHelp = flag.Bool("help", false, "print usage text")
var packageName = flag.String("package", "bebopgen", "specify the name of the package to generate")
var generateUnsafeMethods = flag.Bool("generate-unsafe", false, "whether unchecked additional methods should be generated")

const version = "bebopc-go v0.0.7"

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	flag.Parse()
	if *printHelp {
		flag.Usage()
		return nil
	}
	if *printVersion {
		fmt.Println(version)
		return nil
	}
	if *inputFile == "" {
		return fmt.Errorf("please provide an input file (-i)")
	}
	if *outputFile == "" {
		return fmt.Errorf("please provide an output file (-o)")
	}
	f, err := os.Open(*inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer f.Close()
	bopf, err := bebop.ReadFile(f)
	if err != nil {
		filename := filepath.Base(*inputFile)
		return fmt.Errorf("parsing input failed: %s%w", filename, err)
	}
	out, err := os.Create(*outputFile)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer out.Close()
	settings := bebop.GenerateSettings{
		PackageName:           *packageName,
		GenerateUnsafeMethods: *generateUnsafeMethods,
	}
	if err := bopf.Generate(out, settings); err != nil {
		return fmt.Errorf("failed to generate file: %w", err)
	}
	return nil
}
