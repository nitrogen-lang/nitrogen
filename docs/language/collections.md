# Collections

## Arrays

An array is a numerically indexed collection of items. Items don't have to be the same type.
Arrays look like `[1, "string", true]` and can be indexed using square bracket notation
`var[2]`. Indexing starts at 0. The [collections](../std/imported/collections.ni.md)
package contains several functions to manipulate and manage arrays.

## Hash Maps

Also known as dictionaries or associative arrays, these are data structures that use
key-value pairs. Keys can be strings, ints, or floats. Attempting to use any other data
type will result in an evaluation error. Maps can be created using the syntax
`{"key": "value", "key2": "value2"}`. Map definitions can span multiple lines but
be careful of automatic semicolon insertion, every key-value pair must have a comma after it:

```
myMap = {
    "key": "value",
    "key2": "value2", // Note the trailing comma, without it parsing will fail
}
```

Map values can be retrieved in two ways. The first is standard array index form `myMap["key2"]`.
The other using dot notation `myMap.key2`. Dot notation can be used when the key is a valid
identifier. If the key is not a valid identifier, the normal index notation must be used
with a string. The dot notation can also be used for assignment `myMap.key2 = "another value2"`.
The dot notation is left associative meaning any map index will be resolved before calling a
function. Example: `myMap.key2()` is syntactically the same as `(myMap.key2)()`.
