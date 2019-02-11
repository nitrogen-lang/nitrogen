# collections.ni

This module contains several utility and convenience functions for managing and
working with collections such as arrays and maps.

To use: `import 'std/collections'`

## map(arr: array, func: fn(element, index): T): array

`map` applies the function `fn` on each element of `arr` and returns a new array
with the returned elements.

## filter(arr: array, func: fn(element, index): bool): array

`filter` applies the function `fn` on each element of `arr` and returns a new array
containing the elements of `arr` where `fn` returned true.

## reduce(col: array|map, func: fn(accumulator, element, index): T[, initialValue: T]): T

`reduce` applies a function against an accumulator and each element in the array/map
`col` (from left to right) to reduce it to a single value.

## foreach(col: array|map, func: fn(key, val))

`foreach` will iterate over the supplied collection calling `fn` on each element.
The function `fn` is given the index or map key and the element value. Returned
values are ignored. To actually modify the element, use the `map()` function
instead.

## arrayMatch(arr1, arr2: array): bool

`arrayMatch` returns if `arr1` and `arr2` have the same length and all elements match
in order. If the arrays have the same elements but in different orders, `arrayMatch`
will return false.

## mapMatch(map1, map2: map): bool

`mapMatch` returns if `map1` and `map2` have the same length and all elements match.
`mapMatch` will recursively check nested maps and arrays.

## contains(haystack: array|map, needle: T): bool

`contains` searches `haystack` for `needle` and returns true if the needle is in the
array, false otherwise. If `haystack` is a map, then `contains` returns if the map
has a key `needle`.
