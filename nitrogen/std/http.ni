import "std/encoding/json"

const do = func native (method, url)
const canonicalHeaderKey = func native (header)

const exports = {
    "req": do,
    "canonicalHeaderKey": canonicalHeaderKey,
}

const getJSON = func(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    const resp = get(url, options)
    return json.decode(resp.body)
}
exports.getJSON = getJSON

const get = func(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return do("GET", url, "", options)
}
exports.get = get

const head = func(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return do("HEAD", url, "", options)
}
exports.head = head

const del = func(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return do("DELETE", url, "", options)
}
exports.del = del

const post = func(url) {
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

const put = func(url) {
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

const patch = func(url) {
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
