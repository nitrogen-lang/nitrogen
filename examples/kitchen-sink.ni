// Comment types:
// Single line comment
# Single line comment but using a pound/hash
/*
 * Multi-line comments are also supported
 */

// Here we define a function
def hello = func(place) {
    return "Hello, " + place
}

// Here we assign the output of calling hello()
def helloWorld = hello("World!")

// println is a builtin
println(helloWorld)

// Nitrogen supports if statements
if (helloWorld == "Hello, World!") {
    println("Yep, that's right")
} else {
    println("That's not right...")
}

// And arrays
def places = ["America", "Africa", "Europe"]

println(hello(places[1]))

// And hash maps
def placeMap = {
    "America": ["USA", "Canada", "Mexico"],
    "Europe": ["Germany", "France", "Spain"]
}

println(placeMap["Europe"])

println(push(placeMap["Europe"], "Denmark"))

// And nil (null) for all your null needs
def thisIsNull = nil

println(thisIsNull)
