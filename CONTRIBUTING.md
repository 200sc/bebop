# Contributing

Please raise issues for discussion before opening PRs (or ideally before starting to write code changes) for large or breaking changes,
just to save time if its not the direction we want to take the project.

For small changes (anywhere from typos to non-breaking bug fixes) going straight to a PR is fine, but not required.

## Testing

This library is organized so `go test -v` will test tokenization, parsing, code generation, and the validity of the output generated code. Output code is validated via some integrated tests, and by relying on `go test` itself checking that output can be compiled.

If an example file fails tokenization, it can be added to `testTokenizeFiles` to track that failure and prevent regressions. Likewise for parsing, `TestReadFile` notes expected parsing output. Adding a variant that just ensures a file -can- be parsed, not ensuring a particular AST, is a welcome augmentation.

For generated code, add the example file to `generated/base` and add it's name, minus `.bop`, to `genTestFiles`.
