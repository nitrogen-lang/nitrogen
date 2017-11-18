# Including Scripts

Nitrogen like most scripting languages allows scripts to include other scripts. This allows the programmer to
separate parts of a larger application into manageable pieces. Include paths are relative to the script doing
the including. It is NOT relative to the current working directory.

## include(path: string[,once: bool]): T

`include` will "insert" the script at `path` as if it took the place of the original include call. The included
script is executed within its own block scope. This means, any functions or variables created in the script are
not available to the script that included the file. However, the included script has full access to any variable
in the current scope. An included script can return a variable from the global scope which will be returned by
the include call itself. This allows scripts to export functions or values to the calling script. If `once` is
given and is true, the included script will only be executed once. Additional calls will result in a noop.
Even if the same path is used in a separate include call with once set to false or not provided, the call will
still result in a noop. This state is carried across all scripts. So, using `include(path, true)` will always
include `path` only once.

`include` will return an error type if something went wrong when including or executing the script.

## require(path: string[,once: bool]): T

Same as `include` except an exception is thrown instead of returning an error. This will cause execution to completely
halt if the include fails.

## evalScript(path: script): T

This function will execute a script in an isolated environment. The script will have no access to the environment of the
calling script. `evalScript` will return the value of the called script or an error. The called script retains access to
`_ENV` and `_ARGV`.

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
let otherFile = include('another.ni')

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
let math = require("math.ni", true)

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
require('config.ni', true)

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
require('config.ni', true)

println(config)
```

This will raise an exception of "identifier not found" because `config` only exists in the included scripts environment. Not
the calling script.

### evalScript Example

second.ni

```
println(isDefined('config'))
```

main.ni

```
let config = {}
evalScript('math.ni')
```

This script will print "false" to standard output. In the second script, config is not defined since it's being called
with a clean, isolated environment.
