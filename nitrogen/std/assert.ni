import "std/string"

use string.format

const exports = {}

const isTrue = fn(x) {
    if isFunc(x): x = x()
    if !isBool(x): throw "assertion must be a boolean to isTrue"
    if x: return
    throw format("Assertion Failed: Expected `{}` to be true.", x)
}
exports.isTrue = isTrue

const isFalse = fn(x) {
    if isFunc(x): x = x()
    if !isBool(x): throw "assertion must be a boolean to isFalse"
    if !x: return
    throw format("Assertion Failed: Expected `{}` to be false.", x)
}
exports.isFalse = isFalse

const isEq = fn(a, b) {
    if a == b: return
    throw format("Assertion Failed: Expected `{}` and `{}` to be equal.", a, b)
}
exports.isEq = isEq

const isNeq = fn(a, b) {
    if a != b: return
    throw format("Assertion Failed: Expected `{}` and `{}` to not be equal.", a, b)
}
exports.isNeq = isNeq

const shouldThrow = fn(func) {
    if !isFunc(func): throw "assertion must be a func to shouldThrow"
    try { func() } catch { return }
    throw "Assertion Failed: Expected test to throw."
}
exports.shouldThrow = shouldThrow

const shouldNotThrow = fn(func) {
    if !isFunc(func): throw "assertion must be a func to shouldNotThrow"
    try {
        func()
    } catch e {
        throw format("Assertion Failed: Expected test not to throw. {}", e)
    }
}
exports.shouldNotThrow = shouldNotThrow

return exports
