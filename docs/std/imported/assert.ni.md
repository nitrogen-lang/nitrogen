# assert.ni

Utilities for asserting things during testing. If a test passes, nil is returned.
Otherwise, a string is returned as an error message.

To use: `import 'std/assert'`

## isTrue(x: bool|fn): nil|string

Check if `x` is true. If `x` is a function, it will be executed and its return
value checked. `isTrue` will fail if x, or the return value of x, is not a
boolean type.

## isFalse(x: bool|fn): nil|string

Check if `x` is false. If `x` is a function, it will be executed and its return
value checked. `isFalse` will fail if x, or the return value of x, is not a
boolean type.

## isEq(a, b: T): nil|string

Check if `a` and `b` are equal. `isEq` will also fail if `a` and `b` are not the
same type.

## isNeq(a, b: T): nil|string

Check if `a` and `b` are not equal. `isNeq` will also fail if `a` and `b` are
not the same type.

## shouldRecover(func: fn): nil|string

Run `fn` and check if it recovers from an error.

## shouldNotRecover(func: fn): nil|string

Run `fn` and check if it does not recover from an error.
