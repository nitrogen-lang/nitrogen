# Nitrogen Documentation

Here you will find documentation for the language syntax, as well as the standard library:

- [Language](language.md)
- [Elemental VM](vm.md)
- [Standard Library](stdlib)
- [Including Script Files](stdlib/including-scripts.md)
- [Optional Modules](modules)
- [Module Documentation](modules.md)
- [SCGI Server](scgi-server.md)

## Function Notation

Throughout the documentation, you will find several function definitions. The following syntax is used to denote the number and type
of any function arguments as well as return types. All functions have an implicit return. If a function doesn't list a return type,
it's assumed to be nil.

- Arguments are denoted by their name following by a colon and their type:
  - `func someFunc(arg1: string)` - This functions takes a single arg which must be a string
- Variable arguments are denoted by an ellipse after the argument name:
  - `func someFunc(multiple...: int)` - This function takes 1 or more arguments of type int
- Function return types are denoted by a colon and type after the function definition:
  - `func someFunc(in: string): string` - This function takes a single string argument and returns a string
- Multiple types are denoted using pipes or the generic type "T":
  - `func print(in...: T)` - Print takes 1 or more arguments of any type
  - `func calc(op: string, nums...: string|int): int` - This function takes one argument of type string plus one or more arguments of type string OR int. The function returns an int.
  - `func calc(op: string, nums...: string|int): int|error|nil` - This function takes one argument of type string plus one or more arguments of type string OR int. The function returns an int, error, or nil.
- Multiple arguments with the same type next to each other don't need to specify type:
  - `func hashMerge(map1, map2: map): map` - hashMerge takes two maps. Since they're the same type, the type only needs to be on the last argument.
- Optional arguments are denoted with square braces:
  - `func hashMerge(map1, map2: map[, overwrite: bool]): map` - hashMerge takes two maps and an optional boolean argument.
