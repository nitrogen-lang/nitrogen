const exports = {}

func filter(arr, fn)/*: arr*/ {
    let newArr = [];

    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        if fn(arr[i], i): newArr = push(newArr, arr[i])
    }

    newArr
}
exports.filter = filter

func map(arr, fn)/*: arr*/ {
    let newArr = [];

    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        newArr = push(newArr, fn(arr[i], i))
    }

    newArr
}
exports.map = map

func reduce(arr, fn)/*: Object*/ {
    let accumulator = nil

    if len(arguments) > 0: accumulator = arguments[0]

    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        accumulator = fn(accumulator, arr[i], i)
    }

    accumulator
}
exports.reduce = reduce

func arrayMatch(arr1, arr2)/*: bool*/ {
    if !isArray(arr1) or !isArray(arr2): throw "arrayMatch expected arrays as arguments"
    if len(arr1) != len(arr2): return false

    const ln = len(arr1)
    for i = 0; i < ln; i+=1 {
        if !valuesEqual(arr1[i], arr2[i]): return false
    }

    true
}
exports.arrayMatch = arrayMatch

func mapMatch(map1, map2) {
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

func valuesEqual(v1, v2) {
    if varType(v1) != varType(v2): return false

    if isArray(v1) {
        return arrayMatch(v1, v2)
    } elif isMap(v1) {
        return mapMatch(v1, v2)
    }

    return (v1 == v2)
}

func foreach(collection, fn) {
    if isMap(collection): return foreachMap(collection, fn)
    if isArray(collection) or isString(collection): return foreachArray(collection, fn)
    throw "foreach(): collection must be a map, array, or string"
}
exports.foreach = foreach

func foreachMap(map, fn) {
    const keys = hashKeys(map);
    const keyLn = len(keys)
    for i = 0; i < keyLn; i+=1 {
        const key = keys[i]
        fn(key, map[key])
    }
}

func foreachArray(arr, fn) {
    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        fn(i, arr[i])
    }
}

func contains(arr, needle) {
    if isArray(arr): return arrayContains(arr, needle)
    if isMap(arr): return arrayContains(hashKeys(arr), needle)
    throw "contains expected an array or map but received " + varType(arr)
}
exports.contains = contains

func arrayContains(arr, needle) {
    const arrLen = len(arr)
    const needleType = varType(needle)

    for i = 0; i < arrLen; i += 1 {
        const v = arr[i]
        if varType(v) != needleType: continue
        if (arr[i] == needle): return true
    }

    false
}

return exports
