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
- `>=`: Greater than or equal to ` `<=`: Less than or equal to

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

Blocks that would otherwise contain a single statement do not require
surrounding braces. Instead, add a colon after the condition and write the
statement. This form cannot be used with else or elif branches.

```
if a == b: return c
```

## Match Expressions

The match syntax is an alternative to the if/else syntax that can make code
cleaner and better express intent.

```
match expression {
    "case1" => expression,
    "case2" => {
        do more
    },
    _ => "default case",
}
```

Example:

```
const resp = http.get("")

const message = match resp.status_code {
    200 => "OK",
    400 => "Bad client request",
    500 => "Bad server request",
    _ => "Uknown status",
}

println(message)
```

Case branches must be literal values (string, int, float, bool, nil). The match
syntax does not support numeric conditionals or range checks.

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
loop {
    println("Infinity")
}

// Collection iterator
const items = ["item1", "item2"]
for val in items {
    println(val)
}

// Range iterator
// Same as C style for loop above
for i in range(10) {
    println(i)
}

// While loop
let finished = false
while !finished {
    // Do work eventually setting finished to true
}
```

The initlizer, condition, and incrementor may be enclosed in parentheses, but
they are optional.

A for loop has three parts in the header. An initializer which is ran before the
loop starts, a condition which is evaluated before each iteration, and an
iterator which is ran after the body but before the next condition check.

The initializer must be an assign statement. The condition needs to return a
boolean value. The iterator may be any expression but should increment the loop
counter somehow otherwise it will go in an infinite loop.

Only one variable can be assigned in the initializer.

A while loop is used when only a condition needs to be checked on each iteration
and an iterator isn't required. Think of it as "if some condition, continue to
loop".

### Loop control

The statements `continue` and `break` can be used to control a loop. `continue`
will stop executing the body and begin the next iteration. `break` will stop the
loop completely and continue execution after the loop body.

### Looping over collections and iterators

The `for..in` loop allows looping over a collection of objects. Builtin arrays,
strings, and hashmaps implement this. Custom classes can implement this as well.
When iterating over a collection, both the index/key and value are available.

```
const states = ["Alaska", "Maine", "Florida", "Hawaii"]

for state in states {
    println(state)
}

// Prints:
//
// Alaska
// Maine
// Florida
// Hawaii

for i, state in states {
    println(i, " -> ", state)
}

// Prints:
//
// 0 -> Alaska
// 1 -> Maine
// 2 -> Florida
// 3 -> Hawaii
```

Modifying the length of the collection during iteration is undefined behavior.

Loops with only one value on the left of `in` while be bound to the value of the
current loop. If two values are given separated by a comma, both the index and
value are bound.

#### Custom Iterators

The way iteration works is simple:

- Before starting the loop, the `_iter` method is called on the object being
  iterated.
- The `_iter` method returns an iterator object responsible for returning
  consecutive values. The method can return the instance itself if it implements
  an iterator.
- On each loop, the `_next` method is called on the iterator object. This method
  must return an array where index 0 is the index/key of the current value and
  index 1 is the value.
- To indicate the iterator is finished, the `_next` method must return `nil`.

Example:

```
// Implements the iterator interface to iterate over the collection
class MyCollectionIter {
    let collection
    let i = 0 // Current position in iteration

    const init = fn(collection) {
        this.collection = collection
    }

    const _next = fn() {
        if this.i == len(this.collection.items) {
            return nil // Iterator is finished
        }

        const i = this.i
        const item = this.collection.items[i]
        this.i += 1
        return [i, item] // Return current index and value
    }
}


class MyCollection {
    let items

    const init = fn(items) {
        this.items = items
    }

    const _iter = fn() {
        new MyCollectionIter(this) // Create an iterator over this object
    }
}

const things = new MyCollection(["computer", "mouse", "keyboard"])

for thing in things {
    println(thing)
}

// Prints:
//
// computer
// mouse
// keyboard
```
