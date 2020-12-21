# bebop
[![GoDoc](https://godoc.org/github.com/200sc/bebop?status.svg)](https://godoc.org/github.com/200sc/bebop)
[![Go Reference](https://pkg.go.dev/badge/github.com/200sc/bebop.svg)](https://pkg.go.dev/github.com/200sc/bebop)

bebop is a bebop parser written in Go, for generating Go code.

bebop can read .bop files and output .go files representing them:

```go
package main

import (
    "github.com/200sc/bebop"
)

func main() {
    f, _ := os.Open("mybebop.bop")
    defer f.Close()
    bopf, _ := bebop.ReadFile(f)
    out, _ := os.Create("mybebop.go")
    defer out.Close()
    settings := bebop.GenerateSettings{
        PackageName: "mybebop",
    }
    bopf.Generate(out, settings)
}
```

These files can then be used to encode and decode their represented records:

```go
package server

import (
    "github.com/my/project/mybebop"
    "net"
)

func sendFoo(cn net.Conn) (success bool) {
    fooReq := mybebop.FooRequest{
        Bar: "buzz",
    }
    if err := fooReq.EncodeBebop(cn); err != nil {
        // ...
    }
    fooResp := mybebop.FooResponse{}
    if err := fooResp.DecodeBebop(cn); err != nil {
        // ...
    }
    return fooResp.Success
}

```

## bebobc-go

At `main/bebobc-go` there is a cli utility to take .bop files and produce .go files from them, what is effectively a
little more sophisticated version of the first example in this document.

## bebopfmt

At `main/bebopfmt` there is a cli utility to format and rewrite bop files. It takes in a `-w` argument to rewrite files in place, and accepts both direct files and filepaths.

## Known limitations

The following is a list of known issues with the current version of the project, ordered by approximate priority for addressing them.

Tokenization and parsing errors do not currently report locations (line number, character position / column of the error's source).

Original bebop does not support one .bop file importing type definitions from another .bop file, and so neither does this, yet.

- This is nontrivial, and requires a lot of design toward the importing / packaging ecosystem.

We should be able to generate server code with handlers given all valid messages have op codes.

## Credit

Original bebop compiler (C#, Typescript, ...): https://github.com/RainwayApp/bebop

Valid testdata largely from above project.
