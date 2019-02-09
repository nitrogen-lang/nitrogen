# Classes

Nitrogen has support for simple classes. Classes are like Python where methods
and properties don't have visibility limitations. It's up to the programmer to use
methods as "public" or "private". Classes allow to encapsulate functionality
and data. Instances of a class have different field values while sharing the same method definitions.

## Class Example

```
class name ^ parent {
    let field1
    const field2 = "constant"

    const init = func(x) { // Initializer function
        this.x = x // The instance can be referenced with "this"
    }

    // Class methods
    const doStuff = func(msg) {
        return 'ID: ' + toString(this.x) + ' Msg: ' + msg
    }

    const setX = func(x) {
        this.x = x
    }
}
```

Class `name` extends a class named `parent` and defines several properties and methods.

## Class Inheritance

Classes can inherit methods and fields from another class using the inheritance operator `^`.
Classes may have only one parent. Calling the parent init method can be done by calling
`parent()`. In methods, the variable `parent` is bound to the parent class if one is available.
If a class doesn't have a parent, `parent` is not defined. Parent methods can be retrieved like so: `parent.overridenMethod()`. If a method isn't redefined in a child class, the method can be
called directly without consulting the `parent` variable.

```
class parentPrinter {
    let z

    const init = func() {
        z = "parent thing"
    }

    const doStuff = func(msg) {
        return 'Parent: ' + z + ' Msg: ' + msg
    }

    const parentOnly = func() {
        return "I'm the parent"
    }
}

class printer ^ parentPrinter {
    let x
    const t = "Thing"

    const init = func(x) {
        parent()
        this.x = x
    }

    // Overridden function
    const doStuff = func(msg) {
        return 'ID: ' + toString(x) + ' Msg: ' + msg
    }

    const doStuff2 = func(msg) {
        return parent.doStuff(msg)
    }
}

let myPrinter = new printer(1)

println(myPrinter.doStuff('Hello')) // Redefined on child class
println(myPrinter.z) // Field from parent class
println(myPrinter.parentOnly()) // Method only on parent class
println(myPrinter.doStuff2('Hello')) // Method that calls the parent's doStuff() method
```

## Creating an Instance

An instance of a class can be created using the `new` keyword followed by the class
and any arguments to the init method.

```
class name ^ parent {
    let field1
    const field2 = "constant"

    const init = func(x) { // Initializer function
        this.x = x // The instance can be referenced with "this"s
    }

    // Class methods
    const doStuff = func(msg) {
        return 'ID: ' + toString(x) + ' Msg: ' + msg
    }

    const setX = func(x) {
        this.x = x
    }
}

let myObject = new name(10)

println(classOf(myObject)) // Prints "name"
```
