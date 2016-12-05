def hello = func(place) {
    return "Hello, " + place
}

def helloWorld = hello("World!")

println(helloWorld)

if (helloWorld == "Hello, World!") {
    println("Yep, that's right")
} else {
    println("That's not right...")
}

def places = ["America", "Africa", "Europe"]

println(hello(places[1]))

def placeMap = {
    "America": ["USA", "Canada", "Mexico"],
    "Europe": ["Germany", "France", "Spain"]
}

println(placeMap["Europe"])

println(push(placeMap["Europe"], "Denmark"))
