# Strings

## strSplitN(s, sep: string, n: int): array

`strSplitN` will split `s` on `sep` and return at most `n` array elements. If n is < 0, all substrings will be returned.
If n > 0, at most n substrings will be returned. The returned array may have less than n elements.

## strTrim(s: string): string

`strTrim` removes any whitespace characters from the beginning and end of the string.

## strDedup(s, char: string): string

`strDedup` will reduce any consecutive substring of `char` to a single occurance of `char`. `char` must be a single character.

Example:

```
strDedup("name:    John", " ") == "name: John"
```

This example replaces consecutive strings of spaces with a single space.
