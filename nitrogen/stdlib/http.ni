import "stdlib/encoding/json"

func native do(method, url)
func native canonicalHeaderKey(header)

const exports = {
    "req": do,
    "canonicalHeaderKey": canonicalHeaderKey,
}

func getJSON(url) {
    const resp = get(url)
    return json.decode(resp.body)
}
exports.getJSON = getJSON

func get(url) {
    let options
    if len(arguments) >= 1: options = arguments[0]
    return do("GET", url, "", options)
}
exports.get = get

func post(url) {
    let data
    let options = {}

    if len(arguments) >= 1: data = arguments[0]
    if len(arguments) >= 2: options = arguments[1]

    if !isNull(data) and !isString(data) {
        data = json.encode(data)
        options["Content-Type"] = "application/json"
    }

    return do("POST", url, data, options)
}
exports.post = post

return exports
