# Types

## toInt(in: int|float): int

Convert a number to an int. Some information will be lost when converting from a float to an integer.

## toFloat(in: int|float): float

Convert a number to a float.

## isFloat(in: T): bool
## isInt(in: T): bool
## isBool(in: T): bool
## isString(in: T): bool
## isNull(in: T): bool
## isFunc(in: T): bool
## isArray(in: T): bool
## isMap(in: T): bool

Return if a variable is a specific type.

## parseInt(in: string): int|nil

Attempts to parse the given string as an integer. If parsing fails, nil is returned.

## parseFloat(in: string): float|nil

Same as parseInt() but with floats.

## varType(in: T): string

Returns the type of the variable as a string.

## isDefined(ident: string): bool

Returns if the given identifier is defined.
