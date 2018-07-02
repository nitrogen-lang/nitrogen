import "string"

func String(s) {
    return make string.String(s)
}

func sameTypes(a, b) {
    return varType(a) == varType(b)
}

const exports = {}

func isTrue(x) {
    if !isBool(x): throw "assertion must be a boolean to isTrue"
    if x: return
    throw String("Assertion Failed: Expected {} to be true.").format(x)
}
exports.isTrue = isTrue

func isFalse(x) {
    if !isBool(x): throw "assertion must be a boolean to isFalse"
    if !x: return
    throw String("Assertion Failed: Expected {} to be true.").format(x)
}
exports.isFalse = isFalse

func isEq(a, b) {
    if a == b: return
    throw String("Assertion Failed: Expected {} and {} to be equal.").format(a, b)
}
exports.isEq = isEq

func isNeq(a, b) {
    if !sameTypes(a, b): return
    if a != b: return
    throw String("Assertion Failed: Expected {} and {} to not be equal.").format(a, b)
}
exports.isNeq = isNeq

func shouldThrow(fn) {
    if !isFunc(fn): throw "assertion must be a func to shouldThrow"
    try {
        fn()
    } catch { return }

    throw "Assertion Failed: Expected test to throw."
}
exports.shouldThrow = shouldThrow

func shouldNotThrow(fn) {
    if !isFunc(fn): throw "assertion must be a func to shouldNotThrow"
    try {
        fn()
    } catch e {
        throw String("Assertion Failed: Expected test not to throw. {}").format(e)
    }
}
exports.shouldNotThrow = shouldNotThrow

return exports
