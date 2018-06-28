/*
 * Copyright (c) 2018, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This script tests script imports
 */

// Try a required import
try {
    import './includes/_not_exist.ni'
    println("Test Failed: Import didn't throw")
} catch {pass}

// Test a real import
import './includes/math.ni'

if isError(math.add) {
    println("Test Failed: Include failed")
    println(math.add)
    return
}

if !isFunc(math.add) {
    println("Test Failed: Add is not a function")
    return
}

if math.add(2, 4) != 6 {
    println("Test Failed: Add func failed")
}
