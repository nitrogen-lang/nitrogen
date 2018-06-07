/*
 * Copyright (c) 2018, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This script tests script imports
 */

// Try a required import
try {
    import('./includes/_not_exist.ni')
    println("import didn't throw")
} catch e {
    print("") # Get around a bug with empty blocks
}

// Try and optional import
try {
    always err = import('./includes/_not_exist.ni', false)
    if !isError(err) {
        println("err is not an error")
    }
} catch e {
    println("import threw")
}

// Test a real import
always add = import('./includes/math.ni')

if isError(add) {
    println("Include failed")
    println(add)
    return
}

if !isFunc(add) {
    println("add is not a function")
    return
}

if add(2, 4) != 6 {
    println("add func failed")
}
