# Assignments

Nitrogen supports reassignable and constant variables. All variables must be declared
before they can be assigned. A declaration and assignment can happen at the same time.
Here's some examples:

```
// Standard variable
let var = "Hello, Earth"

// Constant
const var2 = "Hello, Mars"

// Variables must be declared before assignment
// The following works because var was defined above
var = "Goodbye, Earth"

// The following fails since it hasn't be declared
anotherVar = "This will fail"

// Constants can't be reassigned
var2 = "This also causes an error"
```

Compound operations and assignments are supported using the compound operators
`+=`, `-=`, `*=`, `/=`, and `%=`. Each operator will perform the given operation
then assign it to the identifier on the left side:

```
let a = 5

a += 2 // a == 7
a -= 3 // a == 4
a *= 2 // a == 8
a /= 4 // a == 2
a %= 4 // a == 1
```

## Constants

Constants refer to constant references not immutable data. Meaning, a variable
cannot be assigned a different value, but the value of the variable can be changed.

For example, the following code is valid. The variable `obj` is assigned only once,
but the object assigned to obj is still mutable.

```
const obj = {
    a: 1,
    b: 2,
}

obj.a = 3

// obj = {} <- This is invalid because obj cannot be reassigned.
```

## Identifiers

Identifiers are the name of a variable, constant, or function. Identifies must start
with a valid UTF-8 character from a Letter category but can be followed by any number
of letters, decimal digits, or underscores.

## Delete

A non-constant variable can be deleted using the `delete` statement.

```
let someVar = 42
delete someVar
// someVar doesn't exist anymore
```
