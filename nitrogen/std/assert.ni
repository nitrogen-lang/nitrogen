import "std/string"

use string.format

const exports = {}

fn isTrue(x) {
    if isFunc(x): x = x()
    if !isBool(x): return error("assertion must be a boolean to isTrue")
    if x: return
    return error(format("Assertion Failed: Expected `{}` to be true.", x))
}
exports.isTrue = isTrue

fn isFalse(x) {
    if isFunc(x): x = x()
    if !isBool(x): return error("assertion must be a boolean to isFalse")
    if !x: return
    return error(format("Assertion Failed: Expected `{}` to be false.", x))
}
exports.isFalse = isFalse

fn isEq(a, b) {
    if a == b: return
    return error(format("Assertion Failed: Expected `{}` and `{}` to be equal.", a, b))
}
exports.isEq = isEq

fn isNeq(a, b) {
    if a != b: return
    return error(format("Assertion Failed: Expected `{}` and `{}` to not be equal.", a, b))
}
exports.isNeq = isNeq

fn shouldRecover(func) {
    if !isFunc(func): return error("assertion must be a func to shouldRecover")
    const r = recover { func() }
    if isNil(r): return error("Assertion Failed: Expected test to recove.")
}
exports.shouldRecover = shouldRecover

fn shouldNotRecover(func) {
    if !isFunc(func): return error("assertion must be a func to shouldNotRecover")
    const r = recover { func() }
    if !isNil(r): return error(format("Assertion Failed: Expected test not to recove. {}", e))
}
exports.shouldNotRecover = shouldNotRecover

return exports
