import "std/encoding/json"

const do = fn native (method, url)
const canonicalHeaderKey = fn native (header)

const exports = {
    "req": do,
    "canonicalHeaderKey": canonicalHeaderKey,
}

const getJSON = fn(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    const resp = get(url, options)
    return json.decode(resp.body)
}
exports.getJSON = getJSON

const get = fn(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return do("GET", url, "", options)
}
exports.get = get

const head = fn(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return do("HEAD", url, "", options)
}
exports.head = head

const del = fn(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return do("DELETE", url, "", options)
}
exports.del = del

const post = fn(url) {
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

const put = fn(url) {
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

const patch = fn(url) {
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
