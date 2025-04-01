# Literals

Nitrogen literals include integers, floats, booleans, strings, nil, arrays, hash
maps (also known as dictionaries or associative arrays), and functions.

## Nil

Nitrogen's null value is `nil`. `nil` is used to denote the absence of a value.
It's returned from a out of range array or map index among other situations.

## Numbers

Nitrogen makes a distinction between an integer an a floating point number. The
two are not comparable to each other without explicit conversion.

### Integers

Integers are implemented as signed 64-bit numbers giving them a minimum value of
−9,223,372,036,854,775,808 and a maximum of 9,223,372,036,854,775,807. Integers
can be declared using decimal, octal, or hexadecimal notation. Here's a few
examples:

```
45     // Decimal
0o664  // Octal, leading 0
0xA4   // Hexadecimal, prefixed with 0x
0b1010 // Binary, prefixed with 0b
```

Integers support all the standard arithmetic and bitwise operators and are
comparable. Integers may contain underscores "\_" to separate digits:
`10_000_000`. Rules around using the digit separator are very loose. They may
appear anywhere in the number with not limitation on how many are used
consecutively. Obviously, one should be mindful and consistent when using them.

### Floats

Floating point numbers are implemented as 64-bit IEEE floating point numbers.
Floats can only be represented in dotted decimal notation. Exponential notation
is coming soon. Like ints, floats support the standard arithmetic operations and
are comparable. Floats may contain digit separators as well.

## Booleans

Booleans are `true` and `false`. Nitrogen has no such thing as "falsy" or
"truthy" values. Only the literals `true` and `false` evaluate as a boolean
value.

## Strings

Strings are made up of a collection of UTF-8 code points. There are two types of
strings in Nitrogen, interpreted strings and raw strings.

Interpreted strings are surrounded by double quotes and cannot contain any new
lines (it can't span lines), but it can contain escape sequences:

- \b - Backspace
- \e - Escape
- \f - Form feed
- \n - Newline
- \r - Carriage return
- \t - Horizontal tab
- \v - Vertical tab
- \\\\ - Backspace
- \\" - Double quote

If any other escape sequence is found, the backslash and following character are
left untouched. For example the string `"He\llo World"` would not change in its
interpreted form since the escape sequence `\l` isn't valid. It's always good
practice to explicitly escape a backslash rather than relying on this behavior.

Raw strings are slightly different. They're surrounded by single quotes and may
span multiple lines. The only valid escape sequence is `\'`, escaping a single
quote. Raw strings can be helpful for templates or large bodies of inline text.

Strings may be indexed like an array using square brackets `"Hello, world"[0] ==
"H"`. The value of an index expression is another string with the character at
the index of the original string.

## Functions

Functions are literals just like anything else but they have their own docs
[here](functions.md).
