# Nitrogen Module Documentation

Modules are shared libraries that can be loaded into the Nitrogen interpreter at runtime.
These are typically sets of functionality that are still in active development and aren't
part of the core library yet. Once they are stable and have a fixed API, they'll be
considered for the core library.

- [file](file.md): File IO and management
- [filepath](filepath.md): Functions dealing with file paths
- [os](os.md): Interfacing with the OS
- [strings](strings.md): Functions to manipulate strings.

## Imports

Importing modules uses the same mechanism as importing other files. Please see the [import docs](../stdlib/imports.md).

## modulesSupported(): bool

Returns if the platform and build supports dynamic binary modules.

Example:

```
// Attempt to load module, but it's not required
let os = import('os.so', false)
if isError(os) {
    println('Failed loading module os: ', os) // Print error message
}

// All of these do the same thing
print(os->system('whoami')[0]) // Call the system function and print stdout
print(os.system('whoami')[0]) // Call the system function and print stdout
print(os["system"]('whoami')[0]) // Call the system function and print stdout
```
