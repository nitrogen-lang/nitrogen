# Strings

The strings module returns a Module object. All documented functions are part of this returned object.

## splitN(s, sep: string, n: int): array

`splitN` will split `s` on `sep` and return at most `n` array elements. If n is < 0, all substrings will be returned.
If n > 0, at most n substrings will be returned. The returned array may have less than n elements.

## trimSpace(s: string): string

`trimSpace` removes any whitespace characters from the beginning and end of the string.

## dedup(s, char: string): string

`dedup` will reduce any consecutive substring of `char` to a single occurance of `char`. `char` must be a single character.

Example:

```
let strings = module('strings.so')
strings->dedup("name:    John", " ") == "name: John"
```

This example replaces consecutive strings of spaces with a single space.
