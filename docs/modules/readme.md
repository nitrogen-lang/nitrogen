# Nitrogen Module Documentation

Modules are shared libraries that can be loaded into the Nitrogen interpreter at runtime. These are typically sets of functionality
that are still in active development and aren't part of the core library yet. Once they are stable and have a fixed API, they'll
be added to the core library.

- [Files](files.md): File IO and management
- [OS](os.md): Interfacing with the OS
- [Strings](strings.md): Functions to manipulate strings.

## module(filename: string): error

`module` will attempt to import a binary shared library into the interpreter. If the import fails, `module` returns an error object.
Otherwise it returns nil.
