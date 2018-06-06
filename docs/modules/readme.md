# Nitrogen Module Documentation

Modules are shared libraries that can be loaded into the Nitrogen interpreter at runtime. These are typically sets of functionality
that are still in active development and aren't part of the core library yet. Once they are stable and have a fixed API, they'll
be added to the core library.

- [file](file.md): File IO and management
- [filepath](filepath.md): Functions dealing with file paths
- [os](os.md): Interfacing with the OS
- [strings](strings.md): Functions to manipulate strings.

## modulesSupported(): bool

Returns if the platform supports dynamic binary modules.

## module(filename: string[, required: bool]): module|error|nil

`module` will attempt to import a binary shared library into the interpreter. The returned value depends on how the module
registers itself and if required is true. `module` will return nil if the module is found and imported successfully. An error
will be returned if required is false and the module isn't found. An exception is thrown if an error occurs and required is true
or if the module is found but fails importing regardless of `required`. A module object is returned if the module registered
such an object. A module is able to register global functions or an encapsulated module object. Consult the module documentation
for specifics.

If a module returns a Module object, functions or variables can be retrieved using arrow, index, or dot notation.

Example:

```
// Attempt to load module, but it's not required
let os = module('os.so')
if isError(os) {
    println('Failed loading module os: ', os) // Print error message
}

// All of these do the same thing
print(os->system('whoami')[0]) // Call the system function and print stdout
print(os.system('whoami')[0]) // Call the system function and print stdout
print(os["system"]('whoami')[0]) // Call the system function and print stdout
```
