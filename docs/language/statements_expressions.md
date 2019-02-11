# Statements/Expressions

Nitrogen is an expression oriented language. Meaning, almost everything is an expression
and will return its last value. To clarify, an expression is anything that returns a value.
This includes literals which return themselves, as well as functions calls, index lookups, etc.
Statements don't return anything instead they just do something. These include assignments,
loops, and imports.

## Functions

Functions will return either a explicit value when using the `return` statement or
it will return the last expression in the function body.

```
const someFunc = fn() {
    doSomething()
    doSomethingElse()
    "hello"
}

const result = someFunc() // result == "hello"
```

The function above calls two other functions then returns the string "hello". Notice
the `return` keyword wasn't used. Since the string was the last expression in the body
it was returned implicitly. If the string was removed, then the return value of
`doSomethingElse()` would've been returned instead since it would be the last expression.

## If "Statements"

This idea of block expressions extends to the if "statement" as well.

```
const someValue = if thing_is_true {
    "Hello"
} else {
    "World"
}
```

`someValue` will either be "Hello" or "World" depending on if `thing_is_true` is actually true.

## Try/Catch

Try/catch blocks work the same way. They can "return" their last expression. If an exception is
caught, then the last expression of the catch block is returned otherwise the last expression
of the try block is returned.
