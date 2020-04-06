# file.ni

The files module exposes functions to open, close, and manipulate files and directories.

To use: `import 'std/file'`

## fileResourceID: string

The internal ID of a file resource. This can be used with `resourceID()` to check
if a resource is a file.

## open(path, mode: string): resource

`open` will attempt to open the file at `path` using mode `mode`. The returned
value is a file resource that is used by other file methods in this module.

| mode | description                                                                                                                        |
| ---- | ---------------------------------------------------------------------------------------------------------------------------------- |
| r    | Open for reading only                                                                                                              |
| r+   | Open for reading and writing                                                                                                       |
| w    | Open for writing only; truncates the file to zero length; if the file doesn't exist, attempts to create it                         |
| w+   | Open for reading and writing; truncates the file to zero length; if the file doesn't exist, attempts to create it                  |
| a    | Opening for writing only; places file pointer at the end of file (append); if the file doesn't exist, attempts to create it        |
| a+   | Opening for reading and writing; places file pointer at the end of file (append); if the file doesn't exist, attempts to create it |

## close(file: resource)

Closes an open file. If the file is already closed, nothing happens.

## write(file: resource, data: string): int

Writes `data` to `file`. `file` must have been open using a mode that allows
writing otherwise a runtime exception will occur. The function returns the
number of bytes written.

## readAll(filepath: string): string

Reads the entire file at `filepath` and returns its contents as a string.

## remove(filepath: string)

Deletes the file at `filepath`. If the file doesn't exist, nothing happens.

## exists(filepath: string): bool

Returns if the file at `filepath` exists.

## rename(oldPath, newPath: string): error

Attempts to rename a file from `oldPath` to `newPath`. If no error occurs, nil
is returned.

## dirlist(path: string): array

Returns an array which is the directory listing of path. If path is not a
directory, an exception is thrown.

## isdir(path: string): bool

Returns if path is a directory.
