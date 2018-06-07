# Including Scripts

Nitrogen like most scripting languages allows scripts to include other scripts. This allows the programmer to
separate parts of a larger application into manageable pieces.

## import(path: string[,required: bool = true]): T

`import` will do a couple different things depending on whether the imported file is a Nitrogen script,
compiled Nitrogen script, or a shared dynamically linked plugin.

In the case of a script, compiled or not, an import statement will execute the imported script as
if it was part of the current environment. The included script is executed within its own block scope.
This means, any functions or variables created in the script are not available to the script that
included the file. However, the included script has full access to any variable in the current scope.
An included script can return a variable from the global scope which will be returned by
the include call itself. This allows scripts to export functions or values to the calling script.

The second argument determines if an error is returned or an exception is thrown if something fails
during import. If required is true, which is the default, then an exception is thrown. With no try/catch
this can be used to ensure the script always has its needed dependencies. If required is false, an error
object is returned instead which can then be dealt with accordingly. Using error objects is more performant
then using a try/catch.

## Module path resolution

Path resolution is fairly simple. If an import path begins with a period '.' or forward slash '/' then
it's treated as an absolute or relative path. Relative paths are relative to the script file itself not
the working directory.

If a path doesn't begin with a period of slash, the module will be searched for in the module search paths.
The module search paths are available at runtime with the `_SEARCH_PATHS` variable. Paths can be added
by using the `-M` flag on the interpreter binary. Essentially, each path is joined with given import path
until a valid file is found. If a file is found, it will be imported according to its type (script, compiled
script, or shared library). The working directory is always added as the first search path. Any other search
path needs to be added at execution time.

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
let otherFile = import('./another.ni')

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

Modules can be emulated by returning a hashmap and using the arrow syntax:

math.ni:

```
func add(a, b) { a + b}
func sub(a, b) { a - b}
func mul(a, b) { a * b}
func div(a, b) { a / b}

return {
    "add": add,
    "sub": sub,
    "mul": mul,
    "div": div,
}
```

main.ni:

```
let math = import("./math.ni")

println(math->add(2, 3))
println(math->sub(2, 3))
println(math->mul(2, 3))
println(math->div(6, 3))
```

Here, math.ni returns a hashmap that contains several functions. These functions are effectively "exported" by the script.
These functions can then be called from the main script. Here the arrow syntax is used to make it look a but nicer.

### Configuration Example

Here's a typical example of a script calling another script that will configure the application before running.

config.ni

```
config['thing1'] = 'val1'
config['thing2'] = 'val2'
```

main.ni

```
let config = {}
import('./config.ni')

println(config)
```

Note that the config map is created in the main script before including the config file. Remember, included scripts
run in their own enclosed scope. They have access to the calling script's environment at the location of the include,
but nothing created in the second script will be available to the first script unless something is returned from the
second and assigned to a variable in the first.

The following wouldn't work:

config.ni

```
let config = {}
config['thing1'] = 'val1'
config['thing2'] = 'val2'
```

main.ni

```
import('./config.ni')

println(config)
```

This will raise an exception of "identifier not found" because `config` only exists in the included scripts environment. Not
the calling script.
