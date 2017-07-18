# Collections

## len(in: array|string|null): int

Returns the length of an array (number of elements), string (number of bytes), or null (always 0).

## first(in: array): T

Returns the first element of array in.

## last(in: array): T

Returns the last element of array in.

## rest(in: array): array

Returns a new array with elements from in starting at index 1 to the end.

## push(arr: array, val: T): array

Returns a new array with all elements of arr plus the element val added to the end.

## hashMerge(map1, map2: map[, overwrite: bool]): map

Returns a new map with the key-value pairs of map1 combined with those of map2. Map1 acts as the base
map. If the overwrite flag is true, or not provided, keys in map2 with the same name as those in map1
will overwrite the value in map1 with that in map2. If overwrite is false, any duplicate key is
simply ignored. Note, neither input map is modified during the operation.

## hashKeys(in: map): array

Creates and returns an array with the keys of the given map. ***NOTE***: Programmers should NOT rely
on the order of hash map keys. They are not guaranteed to be in a specific order.
