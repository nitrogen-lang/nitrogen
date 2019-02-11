/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates getting input from stdin.
 */

const getName = fn() {
    readline("What's your name? ")
}

const getAge = fn() {
    readline("How old are you? ")
}

println(getName())
println(getAge())
