// Comment types:
// Single line comment
# Single line comment but using a pound/hash
/*
 * Multi-line comments are also supported
 */

// Primitive types:
print("Integer: ", 4);
println("");
print("Float: ", 3.14);
println("");
print("Boolean: ", true);
println("");
print("String: ", "This is a string");
println("");
print("Function: ", func(x, y) { x + y; x * y; });
println("");

// Here we define the variable "hello" and assign it a function
let hello = func(place) {
    return "Hello, " + place;
}

// Syntax sugar for the above
func hello2(place) {
     return "Hello2, " + place;
}

// Here we assign the output of calling hello()
let helloWorld = hello("World!");

// And constants
always helloWorld2 = hello("Mars!");

// The following code would fail since a constant can't be reassigned
// helloWorld2 = "Something else";

// Also, constants must be an int, float, string, bool or null. Arrays, hashmaps or
// other types aren't allowed as constants.

/*
 * Nitrogen features many builtin functions that are implemented directly in the
 * interpreter. These functions generally deal with any type of I/O or
 * manipulation operations.
 */

// "println" will print all arguments with a newline after each
// The similarly named function "print" will print all arguments without a newline.
println(helloWorld);
println(hello2("Earth!"));
println(helloWorld2);

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
 * Currently, only one statement can be evaluated as a condition.
 */

if (helloWorld == "Hello, World!") {
    println("Yep, that's right");
} else {
    println("That's not right...");
}

// Arrays may be arbitrarily long and contain any variable type
let places = ["America", "Africa", "Europe"];

// Like any proper language, Nitrogen is zero-based
println(hello(places[1]));

// Hash maps are also supported. Keys can be either strings or ints.
// Values can be of any type.
let placeMap = {
    "America": ["USA", "Canada", "Mexico"],
    "Europe": ["Germany", "France", "Spain"]
};

println(placeMap["Europe"]);

// Variables in Nitrogen are immutable which means the following push
// does not change the value of placeMap["Europe"].
println(push(placeMap["Europe"], "Denmark"));

// But we can reassign keys and array indices
placeMap["Europe"] = push(placeMap["Europe"], "Norway");
println(placeMap["Europe"]);

// We can add new values to a map
placeMap["Africa"] = ["Egypt", "South Africa", "Madagascar"];
println(placeMap["Africa"]);

placeMap["Africa"][1] = "Ethiopia";
println(placeMap["Africa"]);

println(placeMap);

// And nil (null) for all your null needs
let thisIsNull = nil;

println(thisIsNull);
