/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates a factorial implementation in Nitrogen.
 * 20! being is the largest int that will fit in an int64 data type.
 */

// Calculate n!
func fac(in) {
    // Don't mess with negative numbers, too complex
    if (in < 0) {
        println("in must be non-negative")
        return
    }

    // 0! is defined as 1
    if (in == 0) { return 1 }

    // n! = n * (n-1)!
    return in * fac(in - 1)
}

func main() {
    let facTest = fac(20)

    if (facTest != 2432902008176640000) {
        println("Factorial is broken!")
        println("Got: ", facTest, ", Expected: 2432902008176640000")
        exit(1)
    }
}

// Call entry method
main()
