const encode2 = fn(buf, obj) {
    if isString(obj): return buf + '"' + obj + '"'

    if isInt(obj) or isFloat(obj) or isBool(obj): return buf + toString(obj)

    if isNull(obj): return buf + "null"

    if isArray(obj): return encodeArray(buf, obj)

    if isMap(obj): return encodeMap(buf, obj)

    return error("Unsupported JSON object type: " + varType(obj))
}

const encodeArray = fn(buf, arr) {
    const ln = len(arr)

    buf += '['

    for i = 0; i < ln; i += 1 {
        buf = encode2(buf, arr[i])
        if i < ln-1: buf += ','
    }

    buf + ']'
}

const encodeMap = fn(buf, obj) {
    const keys = hashKeys(obj)
    const ln = len(keys)

    buf += '{'

    for i = 0; i < ln; i += 1 {
        const key = keys[i]
        buf = encode2(buf, key)
        buf += ':'
        buf = encode2(buf, obj[key])

        if i < ln-1: buf += ','
    }

    buf + '}'
}

export fn encode(obj) { encode2("", obj) },
