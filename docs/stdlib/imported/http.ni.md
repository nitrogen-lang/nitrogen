# http.ni

Functions for making HTTP requests.

To use: `import 'stdlib/http'`

## get(url: string[, options: HTTPOptions]): Response

`get` makes an HTTP GET request to the given URL.

## head(url: string[, options: HTTPOptions]): Response

`head` makes an HTTP HEAD request to the given URL.

## del(url: string[, options: HTTPOptions]): Response

`del` makes an HTTP DELETE request to the given URL.

## post(url: string[, data: T, options: HTTPOptions]): Response

`post` makes an HTTP POST request to the given URL. If data is not a string,
it will be JSON encoded and the header `Content-Type` will be set to "json/application".

## put(url: string[, data: T, options: HTTPOptions]): Response

`put` makes an HTTP PUT request to the given URL. If data is not a string,
it will be JSON encoded and the header `Content-Type` will be set to "json/application".

## patch(url: string[, data: T, options: HTTPOptions]): Response

`patch` makes an HTTP PATCH request to the given URL. If data is not a string,
it will be JSON encoded and the header `Content-Type` will be set to "json/application".

## getJSON(url: string[, options: HTTPOptions]): T

Calls `get` with the given URL and returns the output of `json.decode` on the
returned body.

## req(method: string, url: string[, data: string, options: HTTPOptions]): Response

`req` is a low-level command to the native HTTP implementation. `req` can be used
to send requests that aren't possible with the other convenience functions such
as other methods like DELETE or PUT.

## canonicalHeaderKey(s: string): string

`canonicalHeaderKey` returns the canonical format of the header key s. The
canonicalization converts the first letter and any letter following a hyphen to
upper case; the rest are converted to lowercase. For example, the canonical key
for "accept-encoding" is "Accept-Encoding". If s contains a space or invalid header
field bytes, it is returned without modifications.*


\* Description for `canonicalHeaderKey` taken from https://golang.org/pkg/net/http/#CanonicalHeaderKey.

## Response: map

Response is a map with the following structure:

```
{
    "body": string
    "headers": map
}
```

### Fields

#### body

The body of the returned request. No processing is done once received.

#### headers

`headers` is a map of string keys to string values containing the HTTP headers
sent or received. When received, the map keys will be in canonical header format.
If multiple headers with the same are received, the values will be concatenated
and separated by ", ".

## HTTPOptions: map

HTTPOptions is a map with the following structure:

```
{
    "headers": map
}
```

### Fields

#### headers

`headers` is a map of string keys to string values containing the HTTP headers
sent or received. When received, the map keys will be in canonical header format.
It is not currently possible to send multiple headers with the same name.
