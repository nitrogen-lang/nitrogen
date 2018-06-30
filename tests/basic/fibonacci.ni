/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates recursion using the Fibonacci sequence.
 */

func fib(x) {
    if x == 0 or x == 1 {
        return x
    }

    return fib(x-1) + fib(x-2)
}

func main() {
    let fibTest = fib(10)
    if (fibTest != 55) {
        println("Fibonacci is broken!\nGot: ", fibTest, ", Expected: 55")
        exit(1)
    }
}

main()
