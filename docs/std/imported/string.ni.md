# string

string exposes a String class that can be used to manipulate strings.

To use: `import 'std/string'`

## contains(s, substr: string): bool

Returns if string s contains sub-string substr.

## count(s, substr: string): int

Returns the number of non-overlapping instances of substr in s. Count will throw
if substr is empty.

## dedup(s, char: string): string

`dedup` will reduce any consecutive substring of `char` to a single occurrence of `char`. `char` must be a single character.

Example:

```
let strings = module('strings.so')
strings.dedup("name:    John", " ") == "name: John"
```

This example replaces consecutive strings of spaces with a single space.

#### format(s: string, args...: T): string

Format inserts values into the string s. `{}` is used to mark a replacement. Replacements are done
in order by `args`.

##### Example

```
import "string"
const formatted = string.format("My name is {} and I'm {} years old", "John", 25)

formatted == "My name is John and I'm 25 years old"
```

## hasPrefix(s, prefix: string): bool

Returns if the string begins with `prefix`.

## hasSuffix(s, suffix: string): bool

Returns if the string ends with `suffix`.

## replace(s, old, new: string, n: int): string

Returns a copy of s with the first n non-overlapping instances of old replaced
by new. If old is an empty string, replace will throw an exception. If n < 0,
all instances of old are replaced.

## split(s, sep: string): array

Shorthand for `splitN(s, sep, -1)`.

## splitN(s, sep: string, n: int): array

`splitN` will split `s` on `sep` and return at most `n` array elements. If n is < 0, all substrings will be returned.
If n > 0, at most n substrings will be returned. If n == 0, an empty array is returned.
The returned array may have less than n elements.

## trimSpace(s: string): string

`trimSpace` removes any whitespace characters from the beginning and end of the string.

## class String(s: string)

The initializer takes a string literal.

### Fields

- `str: string` - The string.

### Methods

#### contains(substr: string): bool

Returns if string s contains sub-string substr.

#### count(substr: string): int

Returns the number of non-overlapping instances of substr in s. Count will throw
if substr is empty.

#### dedup(char: string): string

`dedup` will reduce any consecutive substring of `char` to a single occurrence of `char`. `char` must be a single character.

##### Example

```
import "string"
const str = new string.String("Hello")
const dedupped = str.dedup("l")

dedupped == "Helo"
```

#### format(args...: T): string

Format inserts values into the string. `{}` is used to mark a replacement. Replacements are done
in order by `args`.

##### Example

```
import "string"
const str = new string.String("My name is {} and I'm {} years old")
const formatted = str.format("John", 25)

formatted == "My name is John and I'm 25 years old"
```

#### hasPrefix(prefix: string): bool

Returns if the string begins with `prefix`.

#### hasSuffix(suffix: string): bool

Returns if the string ends with `suffix`.

#### replace(old, new: string, n: int): string

Returns a copy of s with the first n non-overlapping instances of old replaced
by new. If old is an empty string, replace will throw an exception. If n < 0,
all instances of old are replaced.

#### split(sep: string): array

Shorthand for `splitN(sep, -1)`.

#### splitN(sep: string, n: int): array

`splitN` will split `str` on `sep` and return at most `n` array elements. If n is < 0, all substrings will be returned.
If n > 0, at most n substrings will be returned. The returned array may have less than n elements.

#### trimSpace(): string

`trimSpace` removes any whitespace characters from the beginning and end of the string.
