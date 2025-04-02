# json.ni

Functions for encoding and decoding JSON.

To use: `import 'std/json'`

## encode(obj: T): string|error

`encode` takes an value and converts it into JSON. Values that can't be serialized will cause
`encode` to return an error. Class instances are currently not supported.

## decode(json: string): T|error

`decode` takes a string and returns a Nitrogen value object that represents the parsed JSON
string. Decode may return any valid JSON type including string, int, float, map, array,
boolean, or nil. If the JSON is invalid, `decode` will return an error.
