# Including Scripts

Nitrogen like most scripting languages allows scripts to include other scripts. This allows the programmer to
separate parts of a larger application into manageable pieces.

## import "module/path"[ as var]

- `import "os.so"` - Import the "os.so" module from the search path, bound to the variable `os`
- `import "os.so" as os2` - Import the "os.so" module from the search path, bound to the variable `os2`

`import` will do a couple different things depending on whether the imported file is a Nitrogen script,
compiled Nitrogen script, or a [shared dynamically linked plugin](../../modules).

In the case of a script, compiled or not, an import statement will execute the imported script as
if it was part of the current environment. The included script is executed within its own block scope.
This means, any functions or variables created in the script are not available to the script that
included the file. However, the included script has full access to any variable in the current scope.
An included script can return a variable from the global scope which will be returned by
the include call itself. This allows scripts to export functions or values to the calling script.

## modulesSupported(): bool

Returns if the platform and build supports dynamic binary modules.

## Module path resolution

Path resolution is fairly simple. If an import path begins with a period '.' or forward slash '/' then
it's treated as an absolute or relative path. Relative paths are relative to the script file itself not
the working directory.

If a path doesn't begin with a period of slash, the module will be searched for in the module search paths.
The module search paths are available at runtime with the `_SEARCH_PATHS` variable. Paths can be added
by using the `-M` flag on the interpreter binary. Essentially, each path is joined with the given import path
until a valid file is found. If a file is found, it will be imported according to its type (script, compiled
script, or shared library). The working directory is always added as the first search path. Any other search
path needs to be added at execution time.

Each path is tried with the following extensions in order: ["", ".nib", ".ni", ".so"]. The first simply meaning
the path is checked by itself in case the path includes the extension. Leaving off the extension allows the interpreter
to include a file with the same basename. For example, a compiled `.nib` file can be loaded instead of a `.ni` thereby
removing the need to compile the code before execution. Note, the interpreter does not yet check source mod time or hashes.
This means, if a source file is changed, the corresponding `.nib` file must be compiled manually. The interpreter
won't attempt to create an up-to-date version. If the file is not compiled, the older `.nib` file will continue to be
loaded instead of the `.ni` source file. This is being worked on.

Directories can also be imported similar to Javascript. The interpreter looks for a file named `mod.ni` in the directory
and if found loads that. The `mod.ni` file is responsible for exporting everything the modules needs for its public API.

## Examples

### Simple

second.ni:

```
func otherFile() {
    println("Hello from ", _FILE)
}

return otherFile
```

main.ni:

```
import './second.ni' as otherFile

func main() {
    otherFile()
}

println("Calling main() from ", _FILE)
main()
```

Executing `main.ni` will print two lines, the "Calling main..." string and the "Hello from ..." string.
Notice that the included script returned a function which is saved in the main script to a variable.
That variable is then called like any other function.

### Module Emulation

Modules can be emulated by returning a map:

math.ni:

```
func add(a, b) { a + b }
func sub(a, b) { a - b }
func mul(a, b) { a * b }
func div(a, b) { a / b }

return {
    "add": add,
    "sub": sub,
    "mul": mul,
    "div": div,
}
```

main.ni:

```
import "./math.ni"

println(math.add(2, 3))
println(math.sub(2, 3))
println(math.mul(2, 3))
println(math.div(6, 3))
```

Here, math.ni returns a map that contains several functions. These functions are effectively "exported" by the script.
These functions can then be called from the main script.
