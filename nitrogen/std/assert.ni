import "std/string"

use string.format

const exports = {}

const isTrue = func(x) {
    if isFunc(x): x = x()
    if !isBool(x): throw "assertion must be a boolean to isTrue"
    if x: return
    throw format("Assertion Failed: Expected `{}` to be true.", x)
}
exports.isTrue = isTrue

const isFalse = func(x) {
    if isFunc(x): x = x()
    if !isBool(x): throw "assertion must be a boolean to isFalse"
    if !x: return
    throw format("Assertion Failed: Expected `{}` to be true.", x)
}
exports.isFalse = isFalse

const isEq = func(a, b) {
    if a == b: return
    throw format("Assertion Failed: Expected `{}` and `{}` to be equal.", a, b)
}
exports.isEq = isEq

const isNeq = func(a, b) {
    if a != b: return
    throw format("Assertion Failed: Expected `{}` and `{}` to not be equal.", a, b)
}
exports.isNeq = isNeq

const shouldThrow = func(fn) {
    if !isFunc(fn): throw "assertion must be a func to shouldThrow"
    try { fn() } catch { return }
    throw "Assertion Failed: Expected test to throw."
}
exports.shouldThrow = shouldThrow

const shouldNotThrow = func(fn) {
    if !isFunc(fn): throw "assertion must be a func to shouldNotThrow"
    try {
        fn()
    } catch e {
        throw format("Assertion Failed: Expected test not to throw. {}", e)
    }
}
exports.shouldNotThrow = shouldNotThrow

return exports
