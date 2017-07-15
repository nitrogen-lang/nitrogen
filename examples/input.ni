/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates getting input from stdin.
 */

func getName() {
    readline("What's your name? ")
}

func getAge() {
    readline("How old are you? ")
}

println(getName())
println(getAge())
