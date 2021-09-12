package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/200sc/bebop"
)

var writeInPlace = flag.Bool("w", false, "rewrite the file in place instead of printing to stdout")
var printVersion = flag.Bool("version", false, "print the version of bebopfmt")
var printHelp = flag.Bool("help", false, "print usage text")

const version = "bebopfmt " + bebop.Version

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
	args := os.Args[1:]
	actualArgs := make([]string, 0)
	for _, arg := range args {
		if arg == "-w" || arg == "--w" {
			continue
		}
		actualArgs = append(actualArgs, arg)
	}
	if len(actualArgs) == 0 {
		flag.Usage()
		return fmt.Errorf("\tbobfmt (-w) myfile.bop")
	}

	for _, path := range actualArgs {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Failed to open path: %w", err)
		}
		finfo, err := f.Stat()
		if err != nil {
			f.Close()
			return fmt.Errorf("Failed to stat path: %w", err)
		}
		var files = []string{path}
		if finfo.IsDir() {
			fileInfos, err := f.Readdir(0)
			if err != nil {
				f.Close()
				return fmt.Errorf("Failed to read directory: %w", err)
			}
			files = make([]string, len(fileInfos))
			for i, info := range fileInfos {
				files[i] = filepath.Join(path, info.Name())
			}
		}
		f.Close()
		for _, fpath := range files {
			if err := formatFile(fpath); err != nil {
				return err
			}
		}
	}
	return nil
}

func formatFile(path string) error {
	// Todo: we read the file more than once, its wasteful
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to open path: %w", err)
	}
	_, _, err = bebop.ReadFile(f)
	if err != nil {
		f.Close()
		return fmt.Errorf("Failed to read file: %w", err)
	}
	f.Close()
	f, err = os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to open path: %w", err)
	}

	out := bytes.NewBuffer([]byte{})
	bebop.Format(f, out)

	f.Close()

	if *writeInPlace {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("Failed to open path to rewrite: %w", err)
		}
		f.Write(out.Bytes())
		f.Close()
	} else {
		fmt.Println(string(out.Bytes()))
	}
	return nil
}
