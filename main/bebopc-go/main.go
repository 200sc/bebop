package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/200sc/bebop"
)

var inputFile = flag.String("i", "", "the name of the file to compile")
var outputFile = flag.String("o", "", "the name of the output file to write")
var printVersion = flag.String("version", "", "print the version of the compiler")
var printHelp = flag.String("help", "", "print usage text")
var packageName = flag.String("package", "bebopgen", "specify the name of the package to generate")

const version = "bebopc-go v0.0.3"

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
	if *printHelp != "" {
		flag.Usage()
		return nil
	}
	if *printVersion != "" {
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
		return fmt.Errorf("failed to read input file: %w", err)
	}
	out, err := os.Create(*outputFile)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer out.Close()
	settings := bebop.GenerateSettings{
		PackageName: *packageName,
	}
	if err := bopf.Generate(out, settings); err != nil {
		return fmt.Errorf("failed to generate file: %w", err)
	}
	return nil
}
