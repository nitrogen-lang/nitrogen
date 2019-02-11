/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates recursion using the Fibonacci sequence.
 */

const fib = fn(x) {
    if x == 0 or x == 1 {
        return x
    }

    return fib(x-1) + fib(x-2)
}

const main = fn() {
    for i = 0; i < 31; i += 1 {
        println("Fib of ", i, " is ", fib(i))
    }
}

main()
