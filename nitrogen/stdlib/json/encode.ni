func encode(buf, obj) {
    if isString(obj): return buf + '"' + obj + '"'

    if isInt(obj) or isFloat(obj) or isBool(obj): return buf + toString(obj)

    if isNull(obj): return buf + "null"

    if isArray(obj): return encodeArray(buf, obj)

    if isMap(obj): return encodeMap(buf, obj)

    throw "Unsupported JSON object type: " + varType(obj)
}

func encodeArray(buf, arr) {
    const ln = len(arr)

    buf += '['

    for i = 0; i < ln; i += 1 {
        buf = encode(buf, arr[i])
        if i < ln-1: buf += ','
    }

    buf + ']'
}

func encodeMap(buf, obj) {
    const keys = hashKeys(obj)
    const ln = len(keys)

    buf += '{'

    for i = 0; i < ln; i += 1 {
        const key = keys[i]
        buf = encode(buf, key)
        buf += ':'
        buf = encode(buf, obj[key])

        if i < ln-1: buf += ','
    }

    buf + '}'
}

return {
    "encode": func(obj) { encode("", obj) },
}
