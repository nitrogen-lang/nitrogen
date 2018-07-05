# json.ni

Functions for encoding and decoding JSON.

To use: `import 'stdlib/json'`

## encode(obj: T): string

`encode` takes an value and converts it into JSON. Values that can't be serialized will cause
`encode` to throw an exception. Class instances are currently not supported but will be as soon
as possible.

## decode(json: string): T

`decode` takes a string and returns a Nitrogen value object that represents the parsed JSON
string. Decode may return any valid JSON type including string, int, float, map, array,
boolean, or nil If the JSON is invalid, `decode` will throw and exception.
