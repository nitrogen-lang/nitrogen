# Packages

Nitrogen has the concept of packages to divide programs into smaller, more maintainable units.
There are two types of packages, standard packages and binary modules. Standard packages are
written in Nitrogen while modules are written in Go as a shared library.

## Path resolution

Path resolution is fairly simple. If an import path begins with a forward slash '/' or period '.' then
it's treated as an absolute or relative path. Relative paths are relative to the script file not
the working directory.

If a path doesn't begin with a period of slash, the module will be searched for in the module search paths.
The module search paths are available at runtime as the `_SEARCH_PATHS` variable. Paths can be added
by using the `-M` flag on the interpreter binary. Essentially, each path is joined with the given import path
until a valid file is found. If a file is found, it will be imported according to its type (script, compiled
script, or shared library). The working directory is always added as the first search path. Any other search
path needs to be added at execution time.

Each path is tried with the following extensions in order: ["", ".nib", ".ni", ".so"].
The first simply meaning the path is checked by itself in case the path includes the extension.
Leaving off the extension allows the interpreter to include a file with the same basename.
For example, a compiled `.nib` file can be loaded instead of a `.ni` thereby removing the need to
compile the code before execution. If a `.nib` file is loaded, the corresponding source `.ni`
file is checked for modification time. If the source file is newer than the time recorded
in the nib, the file will be recompiled and the new version will be saved for later loads.
It's highly recommended to never use a file extension except when wanting to load a binary
module that happens to share the same basename as a Nitrogen package.

Directories can also be imported. The interpreter will look for a file named `mod.ni` in the directory
and if found loads that. The `mod.ni` file is responsible for exporting everything the module
needs for its public API.

## Exports

Nothing is exported by a package by default. To export values, definition statements
can be tagged with the `export` keyword. Only exported values are accessible
outside the module.

## Examples

### Simple

second.ni:

```
export fn hello() {
    println("Hello from ", _FILE)
}
```

main.ni:

```
import './second.ni' as otherFile

fn main() {
    otherFile.hello()
}

println("Calling main() from ", _FILE)
main()
```

Executing `main.ni` will print two lines, the "Calling main..." string and the "Hello from ..." string.
Notice that the included script exports a function which can be used in the main script.
That function is then called like any other function.

# Binary Modules

Modules are shared libraries that can be loaded into the Nitrogen interpreter at runtime.
These are typically sets of functionality that are still in active development and aren't
part of the core library yet. Once they are stable and have a fixed API, they'll be
considered for the core library. They can also be third party modules that are implemented
directly in Go instead of Nitrogen.

## Support

Binary shared object modules are only supported on Linux and macOS. This is a limitation
of the underlying Go runtime and there is currently no expectation to support other platforms.

## Importing

Modules are imported the same way as Nitrogen packages.

## Writing Modules

A module can register global functions, create a Module object to encapsulate functionality, or both.
If a module registers a Module object, that object will be bound to the identifier of the `import`
statement. Registered global functions are available immediately after import.
