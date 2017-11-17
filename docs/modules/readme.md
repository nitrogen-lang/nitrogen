# Nitrogen Module Documentation

Modules are shared libraries that can be loaded into the Nitrogen interpreter at runtime. Modules should be in a single
folder and the nitrogen binary needs to be given the `-modules` flag with the path to the modules directory. Any file with
the extension `.so` is loaded as a module.

- [Files](files.md): File IO and management
- [OS](os.md): Interfacing with the OS
- [Strings](strings.md): Functions to manipulate strings.