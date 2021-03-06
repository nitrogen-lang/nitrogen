# Try Catch/Exceptions

Nitrogen has a simple concept of exceptions and runtime errors. Exceptions can be created
at runtime when a script attempts to perform an unsafe action. Errors are returned when
something goes wrong but isn't necessarily a huge issue. Some exceptions are labeled as
panics. When such an exception is thrown, all execution stops. These are thrown when something
happens internally in the interpreter that causes execution to be in a state from which the
program can't recover.

A try/catch block allows catching non-panic exceptions and provides the programmer an
opportunity to handle the exception gracefully. The exception can be bound to an identifier
in the catch block so its message can be printed or checked.

Try and catch blocks are in the same scope as surrounding code. Try/catch is an expression
in Nitrogen and will return the last expression, just like a function, or nil.

Exceptions can be generated by user code using the `throw` keyword followed by some value.

## Try Catch Examples

A try/catch is an expression and will return whatever the try block evaluates to if an
exception isn't thrown:

```
let m1 = try {
    "hello" // Potientially exceptional code
} catch e {
    println('Something bad happened: ', e)
}

println(m1) // "hello"
```

Try blocks can be nested for fallback functionality:

```
try {
    import 'non-existant-module'
} catch e {
    println('Import failed: ', e)
    try {
        import 'std/os'
    } catch e {
        println('Import2 failed: ', e)
    }
}
```

Here, the first try block will fail, but the second one will succeed.

A catch block is required, however it can be empty. In that case it evaluates to nil.

```
try {
    import 'non-existant-module'
} catch {
    pass
}
```

## User generated exceptions

Using the `throw` keyword, a script can also generate an exception:

```
const myException = fn() {
    throw "Nope"
}

try {
    myException()
} catch e {
    println(e) // Will print "Nope"
}
```

Exceptions can also be rethrown:

```
const myException = fn() {
    try {
        myException2()
    } catch e {
        throw e
    }
}

const myException2 = fn() {
    throw "Nope"
}

try {
    myException()
} catch e {
    println('Outer ', e) // Will print "Outer Nope"
}
```
