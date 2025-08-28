# runtime.ni

Runtime information and utilities.

To use: `import 'std/runtime'`

## osName: string

Name of the operating system (darwin, linux, freebsd, windows).

## osArch: string

The system architecture type (amd64, 386).

## dis(func: function): null

`dis` will print the bytecode and other compilation data for a function. `fn`
must be a function.

## disMember(cl: class, method: string): null

`dis_member` will print the bytecode and other compilation data for the named
method of the given class.

## debugVal(arg: T): T

`debugVal()` will print its argument to standard out and return the argument unchanged.
This is useful for debugging to check values without having to create a new variable.
