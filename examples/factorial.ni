/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates a factorial implementation in Nitrogen.
 * This example calculates factorials from 0! -> 20!. 20! being
 * the largest that will fit in an int64 data type.
 */

// Calculate n!
func fac(in) {
     // Don't mess with negative numbers, too complex
     if (in < 0) {
        println("in must be non-negative");
        return;
     }

     // 0! is defined as 1
     if (in == 0) { return 1; }

     // n! = n * (n-1)!
     return in * fac(in - 1);
}

// This is what happens when you don't have loops yet.
// Also this is a great example of a case for tail call
// optimization.
func main(i) {
     // Print extra space for alignment
     if (i < 10) { print(" "); }

     // Print numbers and their results
     print(i);
     print("!");
     print(" = ");
     println(fac(i));

     // 20! is the largest that fits in an int64
     if (i == 20) { return; }

     // Since no loops, use recursion
     main(i+1);
}

// Call entry method
main(0)
