# collections.ni

This module contains several utility and convenience functions for managing and
working with collections such as arrays and maps.

To use: `import('collections.ni')`

## map(arr: array, fn: func(element, index): T): array

`map` applies the function `fn` on each element of `arr` and returns a new array
with the returned elements.

## filter(arr: array, fn: func(element, index): bool): array

`filter` applies the function `fn` on each element of `arr` and returns a new array
containing the elements of `arr` where `fn` returned true.

## reduce(arr: array, fn: func(accumulator, element, index): T[, initialValue: T]): T

`reduce` applies a function against an accumulator and each element in the array
`arr` (from left to right) to reduce it to a single value.

## arrayMatch(arr1, arr2: array): bool

`arrayMatch` returns if arr1 and arr2 have the same length and all elements match
in order. If the arrays have the same elements but in different orders, `arrayMatch`
will return false.
