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
        exit(1)
    }

    if (b <= a) {
        println("b <= a should be false")
        exit(1)
    }

    if (!(a <= c)) {
        println("a <= c should be true")
        exit(1)
    }

    if (!(a >= c)) {
        println("a >= c should be true")
        exit(1)
    }
}

func testCompoundAssignments() {
    let a = 5

    a += 2
    if (a != 7) {
        println("a != 7, got ", a)
        exit(1)
    }

    a -= 3
    if (a != 4) {
        println("a != 4, got ", a)
        exit(1)
    }

    a *= 2
    if (a != 8) {
        println("a != 8, got ", a)
        exit(1)
    }

    a /= 4
    if (a != 2) {
        println("a != 2, got ", a)
        exit(1)
    }
}

testCompoundComparisons()
testCompoundAssignments()
