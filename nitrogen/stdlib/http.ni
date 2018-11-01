import "stdlib/encoding/json"

func native do(method, url)
func native canonicalHeaderKey(header)

const exports = {
    "req": do,
    "canonicalHeaderKey": canonicalHeaderKey,
}

func getJSON(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    const resp = get(url, options)
    return json.decode(resp.body)
}
exports.getJSON = getJSON

func get(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return do("GET", url, "", options)
}
exports.get = get

func head(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return do("HEAD", url, "", options)
}
exports.head = head

func del(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return do("DELETE", url, "", options)
}
exports.del = del

func post(url) {
    let data = if len(arguments) >= 1 { arguments[0] } else { nil }
    let options = if len(arguments) >= 2 { arguments[1] } else { nil }

    if !isNull(data) and !isString(data) {
        data = json.encode(data)

        if isNull(options): options = {}
        options["headers"] = { "Content-Type": "application/json" }
    }

    return do("POST", url, data, options)
}
exports.post = post

func put(url) {
    let data = if len(arguments) >= 1 { arguments[0] } else { nil }
    let options = if len(arguments) >= 2 { arguments[1] } else { nil }

    if !isNull(data) and !isString(data) {
        data = json.encode(data)

        if isNull(options): options = {}
        options["headers"] = { "Content-Type": "application/json" }
    }

    return do("PUT", url, data, options)
}
exports.put = put

func patch(url) {
    let data = if len(arguments) >= 1 { arguments[0] } else { nil }
    let options = if len(arguments) >= 2 { arguments[1] } else { nil }

    if !isNull(data) and !isString(data) {
        data = json.encode(data)

        if isNull(options): options = {}
        options["headers"] = { "Content-Type": "application/json" }
    }

    return do("PATCH", url, data, options)
}
exports.patch = patch

return exports
