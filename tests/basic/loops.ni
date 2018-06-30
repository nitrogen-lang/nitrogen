/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates for loops including the use of continue and break
 */

func main() {
    // Test normal loop
    let outer = 0

    for (i = 0; i < 10; i += 1) {
        outer = outer + 1
    }

    if (outer != 10) {
        println("outer should be 10, got ", outer)
        exit(1)
    }

    // Test skipping every other iteration
    outer = 0

    for (i = 0; i < 10; i += 1) {
        if (i % 2 > 0) { continue }
        outer = outer + 1
    }

    if (outer != 5) {
        println("outer should be 5, got ", outer)
        exit(1)
    }

    // Test breaking
    outer = 0

    for (i = 0; i < 12; i += 1) {
        if (i > 9) { break }
        outer = outer + 1
    }

    if (outer != 10) {
        println("outer should be 10, got ", outer)
        exit(1)
    }
}

main()
