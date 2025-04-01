import "std/string"
import "std/assert"
import "std/os"

const verbose = isString(os.env()['VERBOSE_TEST'])

const exports = {
    "fatal": true,
    "assertLib": assert,
}

fn run(desc, func) {
    let assertLib = exports.assertLib
    let cleanup = nil

    if len(arguments) > 0: cleanup = arguments[0]

    if verbose: println("Test: ", desc)

    let assertionError = recover {
        func(assertLib, check(desc))
    }

    recover {
        if !isNil(cleanup): cleanup()
    }

    if isError(assertionError) {
        printerrln(string.format("Test '{}' failed: {}", desc, assertionError))
        if exports.fatal: exit(1)
    }
}
exports.run = run

fn check(desc) {
    return fn(val) {
        if isError(val) {
            printerrln(string.format("Test '{}' failed: {}", desc, val))
            if exports.fatal: exit(1)
        }
    }
}

return exports
