package fuzz

import (
	"bytes"

	"github.com/200sc/bebop"
)

//go:generate go-fuzz-build

func Fuzz(data []byte) int {
	bopf, err := bebop.ReadFile(bytes.NewReader(data))
	if err != nil {
		return 0
	}
	var w = new(bytes.Buffer)
	err = bopf.Generate(w, bebop.GenerateSettings{
		PackageName: "fuzz",
	})
	if err != nil {
		return 0
	}
	return 1
}
