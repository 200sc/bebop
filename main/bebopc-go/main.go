package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/200sc/bebop"
)

var (
	inputFile             = flag.String("i", "", "the name of the file to compile")
	outputFile            = flag.String("o", "", "the name of the output file to write")
	printVersion          = flag.Bool("version", false, "print the version of the compiler")
	printHelp             = flag.Bool("help", false, "print usage text")
	packageName           = flag.String("package", "bebopgen", "specify the name of the package to generate")
	generateUnsafeMethods = flag.Bool("generate-unsafe", false, "whether unchecked additional methods should be generated")
	shareStringMemory     = flag.Bool("share-string-memory", false, "whether strings read in unmarshalling should share memory with the original byte slice")
	combinedImports       = flag.Bool("combined-imports", false, "whether imported files should be combined and generated as one, or to separate files")
	generateTags          = flag.Bool("generate-tags", false, "whether field tags found in comments should be parsed and generated")
	privateDefinitions    = flag.Bool("private-definitions", false, "whether generated code should be private to the generated package")
	format                = flag.String("formatter", "", "run the given formatter as '$formatter -w $o' after generation")
)

const version = "bebopc-go " + bebop.Version

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
		flag.Usage()
		return fmt.Errorf("please provide an input file (-i)")
	}
	if *outputFile == "" {
		flag.Usage()
		return fmt.Errorf("please provide an output file (-o)")
	}
	f, err := os.Open(*inputFile)
	if err != nil {
		flag.Usage()
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer f.Close()
	bopf, warnings, err := bebop.ReadFile(f)
	if err != nil {
		filename := filepath.Base(*inputFile)
		return fmt.Errorf("parsing input %s failed: %w", filename, err)
	}
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "warning: %v\n", w)
	}
	out, err := os.Create(*outputFile)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer out.Close()
	importMode := bebop.ImportGenerationModeSeparate
	if *combinedImports {
		importMode = bebop.ImportGenerationModeCombined
	}
	settings := bebop.GenerateSettings{
		PackageName:           *packageName,
		GenerateUnsafeMethods: *generateUnsafeMethods,
		SharedMemoryStrings:   *shareStringMemory,
		ImportGenerationMode:  importMode,
		GenerateFieldTags:     *generateTags,
		PrivateDefinitions:    *privateDefinitions,
	}
	if err := bopf.Generate(out, settings); err != nil {
		return fmt.Errorf("failed to generate file: %w", err)
	}
	if *format != "" {
		if err := exec.Command(*format, "-w", *outputFile).Run(); err != nil {
			return fmt.Errorf("failed to run formatter: %w", err)
		}
	}
	return nil
}
