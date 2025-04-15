import "std/string"

use string.format


export fn isTrue(x) {
    if isFunc(x): x = x()
    if !isBool(x): return error("assertion must be a boolean to isTrue")
    if x: return
    return error(format("Assertion Failed: Expected `{}` to be true.", x))
}

export fn isFalse(x) {
    if isFunc(x): x = x()
    if !isBool(x): return error("assertion must be a boolean to isFalse")
    if !x: return
    return error(format("Assertion Failed: Expected `{}` to be false.", x))
}

export fn isEq(a, b) {
    if a == b: return
    return error(
        format(
            "Assertion Failed: Expected {} `{}` and {} `{}` to be equal.",
            varType(a), a,
            varType(b), b,
        ),
    )
}

export fn isNeq(a, b) {
    if a != b: return
    return error(
        format(
            "Assertion Failed: Expected {} `{}` and {} `{}` to not be equal.",
            varType(a), a,
            varType(b), b,
        ),
    )
}

export fn shouldRecover(func) {
    let recoverMsg = nil
    if len(arguments) > 0: recoverMsg = arguments[0]

    if !isFunc(func): return error("assertion must be a func to shouldRecover")
    const r = recover { func() }
    if isNil(r): return error("Assertion Failed: Expected test to recover.")

    if recoverMsg != nil and recoverMsg != toString(r) {
        return error(format("Assertion Failed: Expected test to recover with message `{}` but got `{}`", recoverMsg, r))
    }
}

export fn shouldNotRecover(func) {
    if !isFunc(func): return error("assertion must be a func to shouldNotRecover")
    const r = recover { func() }
    if !isNil(r): return error(format("Assertion Failed: Expected test not to recover. {}", r))
}
