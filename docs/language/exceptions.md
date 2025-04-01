# Runtime Exceptions and Recover

Nitrogen has a simple recovery mechanism to handle runtime errors. Exceptions can be created
at runtime when a script attempts to perform an unsafe or unallowed action. Errors are returned when
something goes wrong but isn't necessarily a huge issue. Some exceptions are labeled as
panics. When such an exception is thrown, all execution stops. These are thrown when something
happens internally in the interpreter that causes execution to be in a state from which the
program cannot recover.

A recover block can catch non-panic exceptions and provides the programmer an
opportunity to handle the exception gracefully. A recovery block has the same
semantics as a [Do Blocks](statements_expressions.md#do-blocks) except if a runtime
exception occurs, execution will continue from the outer scope.

## recover Examples

```
let p = recover {
    1 + "hello"
}
println(p)
println("hello")
```
