# Nitrogen Documentation

Here you will find documentation for the language syntax, as well as the standard library:

- [Language Syntax](language)
- [Standard Library](std)
- [Globals](globals.md)
- [SCGI Server](scgi-server.md)
- [Elemental VM](vm.md)

## Function Notation

Throughout the documentation, you will find several function definitions. The following
syntax is used to denote the number and type of any function arguments as well as
return types. All functions have an implicit return. If a function doesn't list a
return type, it's assumed to be nil.

- Arguments are denoted by their name following by a colon and their type:
  - `someFunc(arg1: string)` - This functions takes a single arg which must be a string
- Variable arguments are denoted by an ellipse after the argument name:
  - `someFunc(multiple...: int)` - This function takes 1 or more arguments of type int
- Function return types are denoted by a colon and type after the function definition:
  - `someFunc(in: string): string` - This function takes a single string argument and returns a string
- Multiple types are denoted using pipes or the generic type "T":
  - `print(in...: T)` - Print takes 1 or more arguments of any type
  - `calc(op: string, nums...: string|int): int` - This function takes one argument of type string plus one or more arguments of type string OR int. The function returns an int.
  - `calc(op: string, nums...: string|int): int|error|nil` - This function takes one argument of type string plus one or more arguments of type string OR int. The function returns an int, error, or nil.
- Multiple arguments with the same type next to each other don't need to specify type:
  - `hashMerge(map1, map2: map): map` - hashMerge takes two maps. Since they're the same type, the type only needs to be on the last argument.
- Optional arguments are denoted with square braces:
  - `hashMerge(map1, map2: map[, overwrite: bool]): map` - hashMerge takes two maps and an optional boolean argument.
