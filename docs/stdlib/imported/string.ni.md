# string

string exposes a String class that can be used to manipulate strings.

To use: `import 'string'`

## class String(s: string)

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

##### Example

```
import "string"
const str = make string.String("Hello")
const dedupped = str.dedup("l")

dedupped == "Helo"
```

#### format(args...: T): string

Format inserts values into the string. `{}` is used to mark a replacement. Replacements are done
in order by `args`.

##### Example

```
import "string"
const str = make string.String("My name is {} and I'm {} years old")
const formatted = str.format("John", 25)

formatted == "My name is John and I'm 25 years old"
```
