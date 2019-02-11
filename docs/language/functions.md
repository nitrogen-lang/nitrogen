# Functions

Functions allow a programmer to break apart a program into separate chunks that focus
on specific tasks. Functions are first class citizens in Nitrogen. They can be passed
around just like any other literal.

Functions are defined using the `fn` keyword:

```
let myFunc = fn(thing) {
    println(thing)
}

// Functions are called like so

myFunc("Some variable")
```

A function definition includes a parameter name list and a statement block.

## Parameters

Functions can be called with parameters. A funciton must be called with at least
the same number of parameters as its declaration. Functions can be called with
more parameters, but they won't assigned to individual identifiers.

```
const noParam = fn() {
    println('This function has no required arguments')
}

const withParam = fn(arg1, arg2) {
    println('This function takes two required arguments')
}
```

Arguments beyond the required ones, are inserted into an array and assigned to
the variable `arguments`.

```
const someFunc = fn() {
    println(arguments) // Prints any arguments passed in
}

someFunc('Hello', 'there') // Will print ['Hello', 'there']
```

Calling a function without the required number of arguments will throw an exception.

```
const someFunc = fn(arg1) {
    println(arg1)
}

someFunc() // Will throw
```

## Variable Scope

All code blocks have their own local scope. Any variable declared inside a function body
will not be visible outside that function. Any variable declared in the environment
in which the function is declared, will be available to that function.
