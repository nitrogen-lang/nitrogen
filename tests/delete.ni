let hello = "world"

if !isDefined("hello") {
    println("Test Failed: hello is not defined")
}

delete hello

if isDefined("hello") {
    println("Test Failed: hello is still defined")
}


const place = "Earth"

if !isDefined("place") {
    println("Test Failed: place is not defined")
}

try {
    delete place
    println("Test Failed: deleting constant didn't throw")
} catch { pass }
