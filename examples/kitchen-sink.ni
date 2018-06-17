// Comment types:
// Single line comment
# Single line comment but using a pound/hash
/*
 * Multi-line comments are also supported
 */

// Primitive types:
println("Integer: ", 4)
println("Float: ", 3.14)
println("Boolean: ", true)
println("String: ", "This is a string", 'This is so a string with single quotes')
println("Function: ", func(x, y) { x + y; x * y; })

// Here we define the variable "hello" and assign it a function
let hello = func(place) {
    return "Hello, " + place
}

// Syntax sugar for the above
func hello2(place) {
    return "Hello2, " + place
}

// Here we assign the output of calling hello()
let helloWorld = hello("World!")

// And constants
const helloWorld2 = hello("Mars!")

// The following code would fail since a constant can't be reassigned
// helloWorld2 = "Something else"

// Also, constants must be an int, float, string, bool or null. Arrays, hashmaps or
// other types aren't allowed as constants.

/*
 * Nitrogen features many builtin functions that are implemented directly in the
 * interpreter. These functions generally deal with any type of I/O or
 * manipulation operations.
 */

// "println" will print all arguments with a newline after each
// The similarly named function "print" will print all arguments without a newline.
println(helloWorld)
println(hello2("Earth!"))
println(helloWorld2)

/*
 * Nitrogen supports simple if statements
 *
 * Like many scripting languages, the following are considered true:
 *  TRUE
 *  Integer or floats not equal to 0
 *  Non-empty strings
 *
 * All other values are considered false.
 *
 * Compound logical expressions can be evalualted using "and" and "or"
 *
 * If statements will short circuit when possible, so below will evaluate properly
 * even though somethingElse isn't defined anywhere.
 */

if helloWorld == "Hello, World!" and true or somethingElse {
    println("Yep, that's right")
} else {
    println("That's not right...")
}

// Arrays may be arbitrarily long and contain any variable type
let places = ["America", "Africa", "Europe"]

// Like any proper language, Nitrogen is zero-based
println(hello(places[1])) // Africa

// Hash maps are also supported. Keys can be either strings or ints.
// Values can be of any type.
let placeMap = {
    "America": ["USA", "Canada", "Mexico"],
    "Europe": ["Germany", "France", "Spain"], // Comma required due to automatic semicolon insertion
}

println(placeMap["Europe"])

// Standard library manipulation functions do not alter the actual array
// The push below does not alter the array, but rather returns a new array
// with the elements pf placeMap["Europe"] plus the new element "Denmark".
println(push(placeMap["Europe"], "Denmark"))

// Map keys and array indices can be reassigned
placeMap["Europe"] = push(placeMap["Europe"], "Norway")
println(placeMap["Europe"])

// We can add new values to a map
placeMap["Africa"] = ["Egypt", "South Africa", "Madagascar"]
println(placeMap["Africa"])

placeMap["Africa"][1] = "Ethiopia"
println(placeMap["Africa"])

println(placeMap)

// And nil (null) for all your null needs
let thisIsNull = nil

println(thisIsNull)

// Functions can take more arguments than declared, this can be used for optional args
func extra(a) {
    // The local variable "args" is an array that contains all parameters after those
    // that were declared. So here, "args[0]" will be the SECOND parameter given since
    // the first paramter is bound to "a".
    println(args)
}

extra(1, 2)

// Simple loops are also possible. An infinite loop can be achieved by omitting the loop header
for i = 0; i < 5; i + 1 {
    println(i)
}
