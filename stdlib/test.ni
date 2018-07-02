import "string"
import "assert"

const verbose = isString(_ENV['VERBOSE_TEST'])

func String(s) {
    return new string.String(s)
}

const exports = {
    "fatal": true
}

func run(desc, fn) {
    if verbose: println("Test: ", desc)

    try {
        fn(assert)
    } catch e {
        println(String("Test '{}' failed: {}").format(desc, e))
        if exports.fatal: exit(1)
    }
}
exports.run = run

return exports
