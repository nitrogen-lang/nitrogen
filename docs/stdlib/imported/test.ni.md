# test.ni

A simple testing framework.

To use: `import 'std/test'`

To print information on each test, set the environment variable `VERBOSE_TEST`
to anything when executing the test script.

## fatal: bool (Default: true)

If true, calls `exit(1)` if a test fails.

## assertLib: T

The assertion object to pass to each test. Defaults to the standard library
assertion package. The object needs to throw and exception to indicate a
particular assertion failed. The exception should be a clear explanation
of how the assert failed.

## run(desc: string, fn: func): nil

`run` represents a single test. `fn` is executed in a try block and the test will
fail if an uncaught exception is bubbled up. `fn` is given a single argument which
is an assertion module. The assertion module is the standard library assert package
by default but can be changed by setting the `assertLib` variable above.

Example:

```
import "std/test"

test.run("Attempt to redefine constant", func(assert) {
    const thing = 42

    assert.shouldThrow(func() {
        thing = 43
    })
})
```
