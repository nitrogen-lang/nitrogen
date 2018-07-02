# test.ni

A simple testing framework.

To use: `import 'test'`

To print information on each test, set the environment variable `VERBOSE_TEST`
to anything when executing the test script.

## run(desc: string, fn: func): nil

`run` represents a single test. `fn` is executed in a try block and the test will
fail if an uncaught exception is bubbled up. `fn` is given a single argument which
is an assertion module. Currently, only the [assert](assert.ni.md) module is the
standard library is supported.

Example:

```
import "test"

test.run("Attempt to redefine constant", func(assert) {
    const thing = 42

    assert.shouldThrow(func() {
        thing = 43
    })
})
```
