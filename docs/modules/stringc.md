# Stringc

stringc exposes a string object that can be used to manipulate strings.

## class string(s: string)

The initializer takes a string literal.

### Fields

- `str: string` - The string.

### Methods

#### splitN(sep: string, n: int): array

`splitN` will split `str` on `sep` and return at most `n` array elements. If n is < 0, all substrings will be returned.
If n > 0, at most n substrings will be returned. The returned array may have less than n elements.

#### trimSpace(): string

`trimSpace` removes any whitespace characters from the beginning and end of the string.

#### dedup(char: string): string

`dedup` will reduce any consecutive substring of `char` to a single occurrence of `char`. `char` must be a single character.

## Example

```
let stringc = module('strings.so')
let str = make stringc.string("Hello")
println(str.str)
```
