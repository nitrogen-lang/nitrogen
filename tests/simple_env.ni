/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file tests to ensure lexical scoping is handled correctly.
 */
 
let string = "Hello, world!";

// Ensure this function changes the outer scope variable
func change_string(next) {
     string = next;
}

// Ensure this one doesn't
func not_change_string(next) {
     // Overshadow the global string variable
     let string = "";
     string = next;
}

if (string != "Hello, world!") {
   println("Outer scope initial value incorrect");
   return;
}

change_string("Hello, mars!");

if (string != "Hello, mars!") {
   println("Outer scope variable wasn't changed");
   return;
}

not_change_string("Hello, earth!");

if (string != "Hello, mars!") {
   println("Outer scope variable was changed");
}
