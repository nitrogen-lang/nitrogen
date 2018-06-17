const exports = {}

exports.filter = func(arr, fn)/*: arr*/ {
    let newArr = [];

    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        if fn(arr[i], i) {
            newArr = push(newArr, arr[i])
        }
    }
    return newArr
}

exports.map = func(arr, fn)/*: arr*/ {
    let newArr = [];

    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        newArr = push(newArr, fn(arr[i], i))
    }
    return newArr
}

exports.reduce = func(arr, fn)/*: Object*/ {
    let accumulator = nil

    if len(arguments) > 0 { accumulator = arguments[0] }

    const ln = len(arr)
    for i = 0; i < ln; i+=1 {
        accumulator = fn(accumulator, arr[i], i)
    }
    return accumulator
}

exports.arrayMatch = func(arr1, arr2)/*: bool*/ {
    if isMap(arr1) or isMap(arr2) { return false }
    if len(arr1) != len(arr2) { return false }

    const ln = len(arr1)
    for i = 0; i < ln; i+=1 {
        if arr1[i] != arr2[i] {
            return false
        }
    }
    return true
}

exports.foreach = func(collection, fn) {
    if isMap(collection) {
        return foreachMap(collection, fn)
    }
    if isArray(collection) {
        return foreachArray(collection, fn)
    }
    throw "foreach(): collection must be a map or array"
}

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

exports.mapMatch = func(map1, map2) {
    throw "This function is broken, need to fix if statment compiler"

    if !isMap(map1) or !isMap(map2) { return false }
    if len(map1) != len(map2) { return false }

    const keys = hashKeys(map1);
    const keyLn = len(keys)
    for i = 0; i < keyLn; i+=1 {
        const key = keys[i]
        if !hasKey(map2, key) {
            return false
        }
        const v1 = map1[key]
        const v2 = map2[key]
        if varType(v1) != varType(v2) {
            return false
        }

        if isArray(v1) {
            if !exports.arrayMatch(v1, v2) {
                return false
            }
            continue
        }

        if isMap(v1) {
            if !exports.mapMatch(v1, v2) {
                return false
            }
            continue
        }

        if map1[key] != map2[key] {
            return false
        }
    }
    return true
}

return exports
