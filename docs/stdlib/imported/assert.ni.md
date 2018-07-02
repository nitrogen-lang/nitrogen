# assert.ni

Utilities for asserting things during testing.

To use: `import 'assert'`

## isTrue(x: bool|func): nil

Check if `x` is true or throw an exception. If `x` is a function, it will be executed and its
return value checked. `isTrue` will throw if x, or the return value of x, is not a boolean type.

## isFalse(x: bool|func): nil

Check if `x` is false or throw an exception. If `x` is a function, it will be executed and its
return value checked. `isFalse` will throw if x, or the return value of x, is not a boolean type.

## isEq(a, b: T): nil

Check if `a` and `b` are equal and if not throw an exception. `isEq` will also throw if `a` and
`b` are not the same type.

## isNeq(a, b: T): nil

Check if `a` and `b` are not equal and if not throw an exception. `isNeq` will also throw if `a`
and `b` are not the same type.

## shouldThrow(fn: func): nil

Run `fn` and check if it throws an error and if it doesn't, throw an exception.

## shouldNotThrow(fn: func): nil

Run `fn` and check if it throws an error and if it does, throw an exception.
