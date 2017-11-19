# Filepath

The filepath module returns a Module object. All documented functions are part of this returned object.

## cwd(): string

`cwd` returns the current working directory. This function may return an empty string if the working directory can't
be determined. This is not common but possible.

## dir(path: string): string

`dir` returns the directory portion of a file path. Effectively, everything before the last directory separator.
This function is the complement to `basename`.

## basename(path: string): string

`basename` returns the file portion of a path. Effectively, everything after the last directory separator.
This function is the complement to `dir`.

## ext(path: string): string

`ext` returns the extension of a file. Effectively, everything after and including the last period. This function
will return an empty string if the file doesn't have an extension.

## abs(path: string): string

`abs` returns path as an absolute filepath starting at the system root directory. This function may return an empty string
if the current working directory cannot be determined. This is not common but possible on some systems.

## join(paths...: string): string

`join` will concatenate all path parts and separate them with the system-specific directory separator.
