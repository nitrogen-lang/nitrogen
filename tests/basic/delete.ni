let hello = "world"

if !isDefined("hello") {
    println("Test Failed: hello is not defined")
    exit(1)
}

delete hello

if isDefined("hello") {
    println("Test Failed: hello is still defined")
    exit(1)
}


const place = "Earth"

if !isDefined("place") {
    println("Test Failed: place is not defined")
    exit(1)
}

try {
    delete place
    println("Test Failed: deleting constant didn't throw")
    exit(1)
} catch { pass }
