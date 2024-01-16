# Collections

## len(in: array|map|string|null): int

Returns the length of an array or map (number of elements), string (number of
bytes), or null (always 0).

## first(in: array): T

Returns the first element of array in.

## last(in: array): T

Returns the last element of array in.

## rest(in: array): array

Returns a new array with elements from in starting at index 1 to the end.

## pop(in: array): array

Returns a new array with the last element of in removed.

## push(arr: array, val: T): array

Returns a new array with all elements of arr plus the element val added to the
end.

## prepend(arr: array, val: T): array

Returns a new array with all elements of arr plus the element val added to the
front.

## splice(arr: array, offset: int[, length: int]): array|error

Returns an array with length elements of arr beginning at offset removed. Length
defaults to the size of the array. splice will return an error if either offset
or length are negative. Using 0 as an offset with no length (thus the default)
will return an empty array.

## slice(arr: array, offset: int[, length: int]): array|error

Returns an array with length elements of arr beginning at offset. Length
defaults to the size of the array. slice will return an error if either offset
or length are negative. Using 0 as an offset with no length (thus the default)
will return a clone of the array.

## sort(arr: array): array

Returns a sorted version of the input array. Array elements must be strings.

## hashMerge(map1, map2: map[, overwrite: bool]): map

Returns a new map with the key-value pairs of map1 combined with those of map2.
Map1 acts as the base map. If the overwrite flag is true, or not provided, keys
in map2 with the same name as those in map1 will overwrite the value in map1
with that in map2. If overwrite is false, any duplicate key is simply ignored.
Note, neither input map is modified during the operation.

## hashKeys(in: map): array

Creates and returns an array with the keys of the given map. **_NOTE_**:
Programmers should NOT rely on the order of hash map keys. They are not
guaranteed to be in a specific order.

## range([start: int, ]end: int[, step: int]): instance

`range` returns an instance implementing an iterator over the integer range.
`range` takes between 1

- 3 arguments. 1 arg is the end with start = 0 and step = 1. 2 args sets start
  and end with step = 1. 3 args sets start, end, and step. Range iterates over
  the range [start, end) meaning `end` is not in the set. `range(10)` returns
  integers 0 - 9. `range(3, 10)` returns integers 3 - 9. `range(0, 10, 2)`
  returns integers 0, 2, 4, 6, 8.
