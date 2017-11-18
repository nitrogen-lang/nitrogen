# Nitrogen

Nitrogen is a dynamically typed, interpreted programming language written in Go. Nitrogen draws inspiration from Go, C, and several other languages.
It's meant to be a simple, easy to use language for making quick scripts and utilities.

## Building the Interpreter

Run `go get github.com/nitrogen-lang/nitrogen/cmd/nitrogen`. This will install the interpreter in your GOBIN path.

## Running the Interpreter

### Interactive Mode

Nitrogen can run in interactive mode much like other interpreted languages. Run Nitrogen with the `-i` flag to start the REPL.

### Scripts

Run Nitrogen like so: `nitrogen filename.ni`. The file extension for Nitrogen files is `.ni`.

### SCGI Server

Nitrogen can run as an SCGI server using multiple workers and the embedded interpreter for performance. Use the `-scgi`
flag to start the server. See the [SCGI docs](docs/scgi-server.md) for more details.

## Documentation

Documentation for the standard library and language is available in the [docs](docs) directory.

## Examples

Example programs can be found in the [examples](examples) directory as well as the [tests](tests) directory.

## Contributing

Issues and pull requests are welcome. Once I write a contributors guide, please read it ;) Until then, always have an issue open for anything
you want to work on so we can discuss. Especially for major design issues.

All code should be ran through `go fmt`. Any request where the files haven't been through gofmt will be denied until they're fixed. Anything
written in Nitrogen, use 4 space indent, keep lines relatively short, and use camelCase for function names. Once classes are implemented, PascalCase
will be used for those.

All contributions must be licensed under the 3-Clause BSD license or a more permissive license such as MIT, or CC0. Any other license will be
rejected.

## License

Both the language specification and this reference interpreter are released under the 3-Clause BSD License which can be found in the LICENSE file.

## Inspiration

- [Writing an Interpreter in Go](https://interpreterbook.com/)
