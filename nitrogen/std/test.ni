import "std/string"
import "std/assert"
import "std/os"

const verbose = isString(os.env()['VERBOSE_TEST'])

export let fatal = true
export let assertLib = assert

export fn run(desc, func) {
    let cleanup = nil

    if len(arguments) > 0: cleanup = arguments[0]

    if verbose: println("Test: ", desc)

    let assertionError = recover {
        func(assertLib, check(desc))
    }

    recover {
        if !isNil(cleanup): cleanup()
    }

    if isError(assertionError) or isException(assertionError) {
        printerrln(string.format("Test '{}' failed: {}", desc, assertionError))
        if fatal: exit(1)
    }
}

fn check(desc) {
    return fn(val) {
        let check_desc = ""
        if len(arguments) > 0: check_desc = arguments[0]
        if verbose and check_desc != "": println("Check: ", check_desc)

        if isError(val) {
            printerrln(string.format("Test '{}' failed: {}", desc, val))
            if fatal: exit(1)
        }
    }
}
