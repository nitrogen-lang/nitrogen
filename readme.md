# Nitrogen

Nitrogen is a dynamically typed, interpreted programming language written in Go. Nitrogen draws inspiration from Go, C, and several other languages.
It's meant to be a simple, easy to use language for making quick scripts and utilities.

## Why Yet Another Language

Nitrogen is not meant to fill a space that's lacking in the programming language. It's more a personal project to design and implement a fully
functional programming language from scratch. My goal is to make Nitrogen functional enough to be used for small scripts and automation.

## Building the Interpreter

Run `go get github.com/nitrogen-lang/nitrogen/cmd/nitrogen`. This will install the interpreter in your GOBIN path.

## Running the Interpreter

### Interactive Mode

Nitrogen can run in interactive mode much like other interpreted languages. Run Nitrogen with the `-i` flag to start the REPL.

### Scripts

Nitrogen can execute a single file as a script (multiple file support is coming soon). Run Nitrogen like so: `nitrogen filename.ni`.
The file extention for Nitrogen files is `.ni`.

## Documentation

Documentation for the standard library and language is available in the [docs](docs) directory.

## Examples

Example programs can be found in the [examples](examples) directory as well as the [tests](tests) directory.

## Contributing

Issues and pull requests are welcome. Once I write a contributors guide, please read it ;) Until then, always have an issue open for anything
you want to work on so we can discuss. Especially for major design issues.

All code should be ran through `go fmt`. Any request where the files haven't been through gofmt will be denied until they're fixed. Anything
written in Nitrogen, use 4 space indent, keep lines relativly short, and use camelCase for function names. Once classes are implemented, PascalCase
will be used for those.

All contributions must be licensed under the 3-Clause BSD license or a more permissive license such as MIT, or CC0. Any other license will be
rejected.

## License

Both the language specification and this reference interpreter are released under the 3-Clause BSD License which can be found in the LICENSE file.

## Inspiration

- [Writing an Interpreter in Go](https://interpreterbook.com/)
