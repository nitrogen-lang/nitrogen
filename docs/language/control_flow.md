# Comparisons/Control Flow

## If Expressions
If expressions in Nitrogen are very similar to other languages:

```
if condition {
    ... do stuff
} elif other_condition {
    ... do something else
} else {
    ... do other stuff
}
```

The condition may be enclosed in parentheses, but they are optional.

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

Blocks that would otherwise contain a single statement do not require surrounding braces.
Instead, add a colon after the condition and write the statement. This form cannot be
used with else or elif branches.

```
if a == b: return c
```

## Loop Statements

Nitrogen supports for and while loops:

```
// Limited loop
for (i = 0; i < 10; i += 1) {
    println(i)
}

// Like if statements, parentheses are optional
for i = 0; i < 10; i += 1 {
    println(i)
}

// Infinite loop
for {
    println("Infinity")
}

let finished = false
while !finished {
    // Do work eventually setting finished to true
}
```

The initlizer, condition, and incrementor may be enclosed in parentheses, but they are optional.

A for loop has three parts in the header. An initializer which is ran before the loop starts,
a condition which is evaluated before each iteration, and an iterator which is ran after the
body but before the next condition check.

The initializer must be an assign statement. The condition needs to return a boolean value.
The iterator may be any expression but should increment the loop counter somehow otherwise it
will go in an infinite loop.

Only one variable can be assigned in the initializer.

An infinite loop can be achieved my simply omitting the entire loop header.

A while loop is used when only a condition needs to be checked on each iteration and
an iterator isn't required. Think of it as "if some condition, continue to loop".

### Loop control

The statements `continue` and `break` can be used to control a loop. `continue` will stop
executing the body and begin the next iteration. `break` will stop the loop completely and
continue execution after the loop body.

### Looping over arrays/maps

Loops over arrays can be done by using the length of the array and then getting the
value from the array by index.

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

Hash maps can be iterated over by getting the map keys with `hashKeys()` and then iterating
over the returned array like above.

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

### Looping with collections.foreach()

The [collections package](../std/imported/collections.ni.md) has a `foreach` method which
loops over an array, map, or string without having to use a loop. The functions take a
function that receives the index and value of each element.

```
let continents = ["Asia", "Africa", "North America", "South America", "Antarctica", "Europe", "Australia"]

import 'std/collections'

collections.foreach(continents, fn(i, v) {
    println(v) // Prints each continent
})
```
