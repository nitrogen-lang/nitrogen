# Types

Type conversion and checking functions.

## toInt(in: int|float): int

Convert a number to an int. Some information will be lost when converting from a float to an integer.

## toFloat(in: int|float): float

Convert a number to a float.

## toString(in: T): string

Convert any value into its stringified form.

## isFloat(in: T): bool

## isInt(in: T): bool

## isBool(in: T): bool

## isString(in: T): bool

## isNull(in: T): bool (alias for isNil)

## isNil(in: T): bool

## isFunc(in: T): bool

## isArray(in: T): bool

## isMap(in: T): bool

## isResource(in: T): bool

## isClass(in: T): bool

## isInstance(in: T): bool

## isError(in: T): bool

## isException(in: T): bool

Return if a variable is a specific type.

## parseInt(in: string): int|nil

Attempts to parse the given string as an integer. If parsing fails, nil is returned.

## parseFloat(in: string): float|nil

Same as parseInt() but with floats.

## varType(in: T): string

Returns the type of the variable as a string.

## isDefined(ident: string): bool

Returns if the given identifier is defined.

## instanceOf(i: T, className: string|class): bool

Returns if object i is an instance of `className`. `className` can be either a string or an actual class object.
`instanceOf` will throw an exception if `className` is not a class or string.

## classOf(i: T): string

Returns the name of the class that i is an instance of. Returns empty string if i is not an object.

## resourceID(i: resource): string

Return the internal ID name of a resource object.
