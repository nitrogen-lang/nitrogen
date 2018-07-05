# IO

## print(args...): null

Print will print all args to standard output with NO space between them. It an arg is
an object with a `toString` method, that method will be called and the return value
will be printed.

## println(args...): null

Same as print() but will also output a newline after printing args.

## printerr(args...): null

Same as `print` but writes to stderr.

## printerrln(args...): null

Same as `println` but writes to stderr.

## printenv(): null

For debugging. Prints the current symbol table as seen by the environment where
printenv() was called.

## readline([prompt: string]): string

readline() will read a line from standard input. If a string argument is given, it
will be printed before taking input. Calling readline() with more than one argument
or with an argument that's not a string, will cause the interpreter to error.

## exit(code: int)

`exit()` terminates script execution and returns with the error code given.
If the script is running in response to an SCGI request, the request is immediately
returned.
