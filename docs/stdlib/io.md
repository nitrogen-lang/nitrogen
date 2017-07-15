# IO

## print(args...)

Print will print all args to standard output with NO space between them.

## println(args...)

Same as print() but will also output a newline after printing args.

## printenv()

For debugging. Prints the current symbol table as seen by the environment where
printenv() was called.

## readline([prompt: string]): string

readline() will read a line from standard input. If a string argument is given, it
will be printed before taking input. Calling readline() with more than one argument
or with an argument that's not a string, will cause the interpreter to error.
