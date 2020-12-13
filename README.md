# bebop

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

## Known limitations

The following is a list of known issues with the current version of the project, ordered by approximate priority for addressing them.


Tokenization and parsing errors do not currently report locations (line number, character position / column of the error's source).

Original bebop does not support one .bop file importing type definitions from another .bop file, and so neither does this, yet.

Original bebop requires semicolons after field definitions, and so do we. It seems practical for the language parser to
treat newlines as semicolons (as Go does).

Invalid recursive type structures can be defined and generated. If they are impossible to size, the Go compiler will reject the generated code.
A future change to this project will detect these conditions at generation time.


## Credit

Original bebop compiler (C#, Typescript, ...): https://github.com/RainwayApp/bebop

Valid testdata largely from above project.
