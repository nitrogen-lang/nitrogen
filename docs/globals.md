# Nitrogen Environment Globals

The interpreter exposes a few global variables that can be used during execution.
Nitrogen reserves all variables starting with a single underscore for interpreter
provided values. Use of variables starting with a single underscore in user code
is discouraged.

## _ENV

`_ENV` is a hashmap of string keys to string values. It contains the environment
variables present in the execution environment of the interpreter. Changing these
values doesn't affect execution or system calls.

## _ARGV

`_ARGV` is a string array that contains all arguments given to the script upon
execution. `len(_ARGV)` will return the number of arguments given. `_ARGV[0]` is
the path name of the main script as it was called.

## _FILE

`_FILE` is the absolute path to the currently executing script.
