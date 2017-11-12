/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates for loops including the use of continue and break
 */

func main() {
    // This loop will go 10 times (0-9) but break on 5 and skip 2
    for i = 0; i < 10; i + 1 {
        // Loops run in an enclosed scope. Variable defined here won't escape to
        // the outer scope. However, like other inner scopes, variables declared
        // before this loop can be modified within the loop

        if i == 2 { continue }

        println(i)

        if i == 5 { break }
    }

    // This shows that the loop counter doesn't leave the loop scope
    println(isDefined("i"))
}

main()
