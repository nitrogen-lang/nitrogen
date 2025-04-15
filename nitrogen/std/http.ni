import "std/encoding/json"

export const req = fn native (method, url)
export const canonicalHeaderKey = fn native (header)

export const getJSON = fn(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    const resp = get(url, options)
    return json.decode(resp.body)
}

export const get = fn(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return req("GET", url, "", options)
}

export const head = fn(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return req("HEAD", url, "", options)
}

export const del = fn(url) {
    let options = if len(arguments) >= 1 { arguments[0] } else { nil }
    return req("DELETE", url, "", options)
}

export const post = fn(url) {
    let data = if len(arguments) >= 1 { arguments[0] } else { nil }
    let options = if len(arguments) >= 2 { arguments[1] } else { nil }

    if !isNull(data) and !isString(data) {
        data = json.encode(data)

        if isNull(options): options = {}
        options["headers"] = { "Content-Type": "application/json" }
    }

    return req("POST", url, data, options)
}

export const put = fn(url) {
    let data = if len(arguments) >= 1 { arguments[0] } else { nil }
    let options = if len(arguments) >= 2 { arguments[1] } else { nil }

    if !isNull(data) and !isString(data) {
        data = json.encode(data)

        if isNull(options): options = {}
        options["headers"] = { "Content-Type": "application/json" }
    }

    return req("PUT", url, data, options)
}

export const patch = fn(url) {
    let data = if len(arguments) >= 1 { arguments[0] } else { nil }
    let options = if len(arguments) >= 2 { arguments[1] } else { nil }

    if !isNull(data) and !isString(data) {
        data = json.encode(data)

        if isNull(options): options = {}
        options["headers"] = { "Content-Type": "application/json" }
    }

    return req("PATCH", url, data, options)
}
