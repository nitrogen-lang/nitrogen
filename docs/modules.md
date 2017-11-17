# Language Modules

The Nitrogen interpreter supports dynamically linking to shared Go plugins to extend the functionality of the interpreter
by providing additional userland functions. These modules are only supported on Linux.

## Using Modules

All modules must end in `.so` and exist in a single directory. When starting Nitrogen, use the `-modules` flag to tell Nitrogen
where to load modules. All files in the directory that end in `.so` will be loaded.
