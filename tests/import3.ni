/*
 * Copyright (c) 2018, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This script tests script imports
 */

// Test directory import
import './includes/math2' as math

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
