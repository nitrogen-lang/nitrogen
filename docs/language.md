# Nitrogen Language

## Source Files

Nitrogen source files are read as valid UTF-8 characters.

## Semicolons

Nitrogen's formal grammar uses semicolons to denote the end of a statement or expression. However, Nitrogen code doesn't need to have
explicit semicolons as the lexer will automatically insert a semicolon where it's needed. Namely, after an identifier, literal,
nil, a closing parenthesis, curly or square bracket, and after the keyword "return". This requires the programmer to keep in mind
how they're formatting code. For example, the following if statement is invalid:

```
if (var == 3) // The lexer attempts to insert a semicolon here, which will fail.
{
    print("It's a 3")
}
else
{
    print("No, it's not 3")
}
```

But the following is ok:

```
if (var == 3) { print("It's a 3") } else { print("No, it's not 3") }
```

The parser will catch any such errors and will warn the programmer.

## Comments

Nitrogen supports three styles of comments:

```
// This comment lasts to the end of the line
# As does this one, this is mainly for scripts and #! headers
/*
 * Multiline comments are also supported.
 */
```

## Literals

Nitrogen supports integers, floats, booleans, strings, null, arrays, and hash maps (dictionaries, associative arrays).

### Null

Nitrogen has a null value: `nil`. Nil can be returned from functions, returned from a bad array or map indexing, or
from other things.

### Numbers

Note: Currently, ints and floats CAN NOT be compared to each other. To compare, use the toInt() or toFloat() functions to convert
values between types.

#### Integers

Integers are implemented using Go's int64 type which means all ints are signed and 64 bits long. Integers can be declared using
decimal, octal, or hexadecimal notation. Here's a few examples:

```
45    // Decimal
0664  // Octal, leading 0
\xA4  // Hexadecimal, prefixed with \x
```

Ints support the standard arithmatic operations: addition, subtraction, multiplication, division, and modulo.
And of course ints can be compared to each other using <, >, ==, and !=.

#### Floats

Floating point numbers are implemented using Go's float64 type which means they are the same as a double in C or Java. Floats can only be represented
in dotted decimal notation. Exponential notation is coming soon. Like ints, floats support the standard arithmatic operations: addition, subtraction,
multiplication, and division. Floats may be compared to each other.

### Booleans

Booleans are simple `true` and `false`.

### Strings

Strings are made up of arbitrary bytes that may or may not represent a UTF-8 encoded string. There are two types of strings in Nitrogen,
interpreted strings and raw strings.

Interpreted strings are surrounded by double quotes and cannot contain any new lines (it can't span lines), but it can contain escape sequences:

- \b - Backspace
- \n - Newline
- \r - Carriage return
- \t - Horizontal tab
- \v - Vertical tab
- \f - Form feed
- \\\\ - Backspace
- \\" - Double quote

If any other escape sequence is found, the backslash and following character are left untouched. For example the string `"He\llo World"` would
not change in its interpreted form since the escape sequence `\l` isn't valid. It's always good practice to explicitly escape a backslash
rather than relying on this behavior.

Raw strings are slightly different. They're surrounded by single quotes and may span multiple lines. The only valid escape sequence
is `\'`, escaping a single quote. Raw strings can be helpful for templates or large bodies of inline text.

## Collections

### Arrays

Arrays work similarly to other languages. They take the form: `[1, "string", true]`. Array elements can be of any type. Arrays can be
indexed using square bracket notation `var[2]`. Nitrogen, like any proper language, uses 0-based array indexing. Please consult the standard
library documentation for functions that can manipulate arrays.

### Hash Maps

Also known as dictionaries or associative arrays are data structures that use key-value pairs. Keys can be strings, ints, or floats. Attempting
to use any other data type will result in an evaluation error. Maps can be created using the syntax `{"key": "value", "key2": "value2"}`.
Map definitions can span multiple lines but be careful of automatic semicolon insertion, every key-value pair must have a comma after it:

```
myMap = {
    "key": "value",
    "key2": "value2", // Note the trailing comma, without it parsing will fail
}
```

## Assignments

Nitrogen supports both variables and constants. All variables must be declared before they can be assigned. A declaration and assignment can
happen at the same time. Here's some examples:

```
// Standard variable
let var = "Hello, Earth"

// Constant
always var2 = "Hello, Mars"

// Variables must be declared before assignment
// The following works because var was defined above
var = "Goodbye, Earth"

// The following fails since it hasn't be declared
anotherVar = "This will fail"

// Obviously, constants can't be changed
var2 = "This also causes an error"
```

Compound operations and assignments are supported using the compound operators +=, -=, *=, and /=. Each operator will perform the given operation
then assign it to the identifier on the left side:

```
let a = 5

a += 2 // a == 7
a -= 3 // a == 4
a *= 2 // a == 8
a /= 4 // a == 2
```

### Constants

Constants can be a string, int, float, or bool. They CANNOT be an array or map.

## Identifiers

Identifiers are the name of a variable, constant, or function. Identifies must start with a valid UTF-8 character from a Letter category
but can be followed by any number of letters, decimal digits, or underscores.

## Functions

Functions allow a programmer to break apart a program into separate chunks that focus on specific tasks. Functions are first class citizens
in Nitrogen. They can be passed around just like any other variable.

```
// Functions can be defined using two syntaxes

let myFunc = func(thing) {
    println(thing)
}

// or

func myFunc(thing) {
    println(thing)
}

// Attempting to combine the two will result in an error

let myFunc = func myFunc(thing) { println(thing) } // This is bad

// Functions are called like so

myFunc("Some variable")
```

## Comparisons/Control Flow

If expressions in Nitrogen are very similar to other languages:

```
if condition {
    ... do stuff
} else {
    ... do other stuff
}
```

The condition may be enclosed in parentheses, but they are completely optional.

Nitrogen supports standard comparison operators:

- `==`: Equal
- `!=`: Not equal
- `>`: Greater than
- `<`: Less than
- `>=`: Greater than or equal to
` `<=`: Less than or equal to

An expression can be prefixed with the bang operator to negate it:

```
!true == false
```

Compound comparisons are also possible with the keywords `and` and `or`:

```
if a == b or a == c {
    ... do stuff
}

if a == b and b == c {
    ... then a == c
}
```

Conditions can be groups to change the order or precidence:

```
if a == b or (a == c and a == d) {
    ... do more stuff
}
```

## For loops

Nitrogen supports a version of the traditional C for loop:

```
// Limited loop
for (i = 0; i < 10; i + 1) {
    println(i)
}

// Parentheses are optional
for i = 0; i < 10; i + 1 {
    println(i)
}

// Infinite loop
for {
    println("Infinity")
}
```

The initlizer, condition, and incrementor may be enclosed in parentheses, but they are completely optional.

A for loop has three parts in the header. An initializer which is ran before the loop starts, a condition which is evaluated before each
iteration, and an iterator which is ran after the body but before the next condition check.

The initializer must be an assign statement. The keyword `let` is not needed. The condition needs to return a boolean value. See Nitrogen's
boolean logic for what constitutes as true/false. The iterator should be an assign or other expression. If it's not an assign statment, the
value of the iterator will be assigned to the initialized identifier. For example, in the above for loop, since the iterator isn't assigning
the value of `i + 1` to a variable, it will automatically be assigned to `i`. The assignment to i has nothing to do with i in the iterator,
but because i is in the initalizer.

Only one variable can be assigned in the initializer.

An inifinate loop can be achieved my simply omitting the entire loop header.

### Loop control

The statements `continue` and `break` can be used to control a loop. `continue` will stop executing the body and begin the next iteration.
`break` will stop the loop completely and continue execution after the loop body.

### Looping over arrays/maps

Loops over arrays can be done by using the length of the array and then getting the value from the array by index.

```
let arr = ["one", "two", "three"]

for (i = 0; i < len(arr); i + 1) {
    println(arr[i])
}

// Outputs:
//  one
//  two
//  three
```

Hash maps can be iterated over by getting the map keys with `hashKeys()` and then iterating over the returned array like above.

```
let map = {
    "key1": "value1",
    "key2": "value2",
    "key3": "value3",
    "key4": "value4",
}

let keys = hashKeys(map)

for (i = 0; i < len(keys); i + 1) {
    ley key = keys[i]
    println(key, ": ", map[key])
}

// Output:
//  key1: value1
//  key2: value2
//  key3: value3
//  key4: value4
```
