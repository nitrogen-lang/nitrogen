# Nitrogen Module Documentation

Modules are shared libraries that can be loaded into the Nitrogen interpreter at runtime.
These are typically sets of functionality that are still in active development and aren't
part of the core library yet. Once they are stable and have a fixed API, they'll be
considered for the core library.

- [file](file.md): File IO and management
- [filepath](filepath.md): Functions dealing with file paths
- [os](os.md): Interfacing with the OS

## Support

Binary shared object modules are only supported in Linux and macOS. This is a limitation of the underlying Go runtime
and there is currently no expectation to support other platforms.

## Importing

Modules can be imported in two ways. The first is when the interpreter starts but before a script is executed. Use the `-M`
flag to specify import search directories (the working directory is added by default). Then use the `-al` flag to pre-load
specific modules. Scripts still need to use the `import` statement to retrieve any module object created by the module.
Pre-loading modules can be used to add extra global objects or provide some other extra functionality before any script is executed.

Importing modules uses the same mechanism as importing other files. Please see the [import docs](../stdlib/global/imports.md).

Example:

```
// Attempt to load module
import 'os.so'

// All of these do the same thing
print(os.system('whoami')[0]) // Call the system function and print stdout
print(os["system"]('whoami')[0]) // Call the system function and print stdout
```

## Writing Modules

A module can register global functions, create a Module object to encapsulate functionality, or even both. If a module registers
a Module object, that object will be bound to the identifier of the `import` statement. Registered global functions are available
immediately after import.
