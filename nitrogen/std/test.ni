import "std/string"
import "std/assert"

const verbose = isString(_ENV['VERBOSE_TEST'])

const exports = {
    "fatal": true,
    "assertLib": assert,
}

func run(desc, fn) {
    if verbose: println("Test: ", desc)

    try {
        fn(exports.assertLib)
    } catch e {
        printerrln(string.format("Test '{}' failed: {}", desc, e))
        if exports.fatal: exit(1)
    }
}
exports.run = run

return exports
