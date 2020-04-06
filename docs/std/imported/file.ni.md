# file.ni

The files module exposes functions to open, close, and manipulate files and directories.

To use: `import 'std/file'`

## readFile(filepath: string): string

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

## class File(path: string, mode: string)

Represents a file object. Creating this class will attempt to open the file
at `path` with mode `mode`.

| mode | description                                                                                                                        |
| ---- | ---------------------------------------------------------------------------------------------------------------------------------- |
| r    | Open for reading only                                                                                                              |
| r+   | Open for reading and writing                                                                                                       |
| w    | Open for writing only; truncates the file to zero length; if the file doesn't exist, attempts to create it                         |
| w+   | Open for reading and writing; truncates the file to zero length; if the file doesn't exist, attempts to create it                  |
| a    | Opening for writing only; places file pointer at the end of file (append); if the file doesn't exist, attempts to create it        |
| a+   | Opening for reading and writing; places file pointer at the end of file (append); if the file doesn't exist, attempts to create it |

### Fields

### Methods

#### close()

Closes the open file. If the file is already closed, nothing happens.

#### write(data: string): int

Writes `data` to file. File must have been open using a mode that allows
writing otherwise a runtime exception will occur. The function returns the
number of bytes written.

#### readAll(): string

Reads the entire file contents and returns it as a string.

#### readLine(): string

Reads a single line from the file and returns it.

#### readChar(): string

Reads a single character from the file and returns it.

#### remove(): null

Closes the file and deletes it.

#### rename(newpath: string): null

Closes the file and renames it.
