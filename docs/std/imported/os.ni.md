# os.ni

The os module returns a Module object. All documented functions are part of this returned object.

To use: `import 'std/os'`

## system(cmd: string[, args: array]): array|error|nil

`system` will execute `cmd` with arguments `args`. `args` must be an array of strings.
The returned array contains the standard output at index 0 and standard error at
index 1 of the executed command. An error object is returned if the command failed to execute.

## exec(cmd: string[, args: array]): error|nil

`exec` like `system` will will execute a system command. However `exec` will link
up the standard input, output, and error of the interpreter to the command. This
allows a script to give control to a user for interactive commands. An error object is
returned if the command failed to execute, nil otherwise.

## argv(): array

`argv` returns an array of command line arguments used when executing the application.

## env(): map

`env` returns a hashmap of the environment variables.
