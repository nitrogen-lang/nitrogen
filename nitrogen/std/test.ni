import "std/string"
import "std/assert"
import "std/os"

const verbose = isString(os.env['VERBOSE_TEST'])

const exports = {
    "fatal": true,
    "assertLib": assert,
}

const run = fn(desc, func) {
    let cleanup = nil

    if len(arguments) > 0: cleanup = arguments[0]

    if verbose: println("Test: ", desc)

    try {
        func(exports.assertLib)
        if !isNil(cleanup): cleanup()
    } catch e {
        if !isNil(cleanup): cleanup()

        printerrln(string.format("Test '{}' failed: {}", desc, e))
        if exports.fatal: exit(1)
    }
}
exports.run = run

return exports
