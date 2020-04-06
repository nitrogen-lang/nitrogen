const exports = {}

const filter = fn(arr, func)/*: arr*/ {
    let newArr = [];

    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        if func(arr[i], i): newArr = push(newArr, arr[i])
    }

    newArr
}
exports.filter = filter

const map = fn(arr, func)/*: arr*/ {
    let newArr = [];

    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        newArr = push(newArr, func(arr[i], i))
    }

    newArr
}
exports.map = map

const reduce = fn(collection, func)/*: Object*/ {
    let accumulator = nil

    if len(arguments) > 0: accumulator = arguments[0]

    if isArray(collection): return reduceArray(collection, func, accumulator)
    if isMap(collection): return reduceMap(collection, func, accumulator)
    throw "reduce(): collection must be a map or array"
}
exports.reduce = reduce

const reduceArray = fn(arr, func, accumulator)/*: Object*/ {
    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        accumulator = func(accumulator, arr[i], i)
    }

    accumulator
}

const reduceMap = fn(map, func, accumulator)/*: Object*/ {
    const keys = hashKeys(map);
    const keyLn = len(keys)
    for i = 0; i < keyLn; i+=1 {
        const key = keys[i]
        accumulator = func(accumulator, map[key], key)
    }

    accumulator
}

const arrayMatch = fn(arr1, arr2)/*: bool*/ {
    if !isArray(arr1) or !isArray(arr2): throw "arrayMatch expected arrays as arguments"
    if len(arr1) != len(arr2): return false

    const ln = len(arr1)
    for i = 0; i < ln; i+=1 {
        if !valuesEqual(arr1[i], arr2[i]): return false
    }

    true
}
exports.arrayMatch = arrayMatch

const mapMatch = fn(map1, map2) {
    if !isMap(map1) or !isMap(map2): return false
    if len(map1) != len(map2): return false

    const keys = hashKeys(map1);
    const keyLn = len(keys)
    for i = 0; i < keyLn; i+=1 {
        const key = keys[i]
        if !hasKey(map2, key): return false
        if !valuesEqual(map1[key], map2[key]): return false
    }

    true
}
exports.mapMatch = mapMatch

const valuesEqual = fn(v1, v2) {
    if varType(v1) != varType(v2): return false

    if isArray(v1) {
        return arrayMatch(v1, v2)
    } elif isMap(v1) {
        return mapMatch(v1, v2)
    }

    return (v1 == v2)
}

const foreach = fn(collection, func) {
    if !isMap(collection) and !isArray(collection) and !isString(collection) {
        throw "foreach(): collection must be a map, array, or string"
    }

    for key, val in collection {
        func(key, val)
    }
}
exports.foreach = foreach

const contains = fn(arr, needle) {
    if isArray(arr): return arrayContains(arr, needle)
    if isMap(arr): return arrayContains(hashKeys(arr), needle)
    throw "contains expected an array or map but received " + varType(arr)
}
exports.contains = contains

const arrayContains = fn(arr, needle) {
    const arrLen = len(arr)
    const needleType = varType(needle)

    for i = 0; i < arrLen; i += 1 {
        const v = arr[i]
        if varType(v) != needleType: continue
        if (arr[i] == needle): return true
    }

    false
}

const join = fn(separator, arr) {
    const arrLen = len(arr)
    let str = ""

    for i = 0; i < arrLen; i += 1 {
        str += toString(arr[i])
        if i < arrLen - 1 {
            str += separator
        }
    }

    str
}
exports.join = join

return exports
