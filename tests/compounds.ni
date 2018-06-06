/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This script tests the <= and >= comparison operators.
 */

func testCompoundComparisons() {
    let a = 5
    let b = 6
    let c = 5

    if (a >= b) {
        println("a >= b should be false")
    }

    if (b <= a) {
        println("b <= a should be false")
    }

    if (!(a <= c)) {
        println("a <= c should be true")
    }

    if (!(a >= c)) {
        println("a >= c should be true")
    }
}

func testCompoundAssignments() {
    let a = 5

    a += 2
    if (a != 7) {
        println("a != 7, got ", a)
    }

    a -= 3
    if (a != 4) {
        println("a != 4, got ", a)
    }

    a *= 2
    if (a != 8) {
        println("a != 8, got ", a)
    }

    a /= 4
    if (a != 2) {
        println("a != 2, got ", a)
    }
}

testCompoundComparisons()
testCompoundAssignments()
