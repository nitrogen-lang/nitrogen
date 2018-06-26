# json.ni

Functions for encoding and decoding JSON.

To use: `import 'json.ni'`

## encode(obj: T): string

`encode` takes an value and converts it into JSON. Values that can't be serialized will cause
`encode` to throw an exception. Class instances are currently not supported but will be as soon
as possible.
