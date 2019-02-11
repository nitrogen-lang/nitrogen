import "std/string"
import "std/assert"
import "std/os"

const verbose = isString(os.env()['VERBOSE_TEST'])

const exports = {
    "fatal": true,
    "assertLib": assert,
}

const run = fn(desc, func) {
    if verbose: println("Test: ", desc)

    try {
        func(exports.assertLib)
    } catch e {
        printerrln(string.format("Test '{}' failed: {}", desc, e))
        if exports.fatal: exit(1)
    }
}
exports.run = run

return exports
