/*
 * Copyright (c) 2018, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This script tests script imports
 */

// Try a required import
try {
    import('./includes/_not_exist.ni')
    println("Test Failed: Import didn't throw")
} catch {pass}

// Try and optional import
try {
    const err = import('./includes/_not_exist.ni', false)
    if !isError(err) {
        println("Test Failed: Err is not an error")
    }
} catch e {
    println("Test Failed: Import threw")
}

// Test a real import
const add = import('./includes/math.ni')

if isError(add) {
    println("Test Failed: Include failed")
    println(add)
    return
}

if !isFunc(add) {
    println("Test Failed: Add is not a function")
    return
}

if add(2, 4) != 6 {
    println("Test Failed: Add func failed")
}
