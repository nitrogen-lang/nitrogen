let fastVar = "I'm fast"

func myException() {
    myException2()
}

func myException2() {
    throw "Nope"
}

let m1 = try {
    myException()
} catch e {
    errorVal(e)
}

always expected = "Nope"
if m1 != expected {
    println('m1 is not the correct value')
    println('Expected "', expected, '", got "', m1, '"')
}

if !isDefined(fastVar) {
    println('Try/catch block test failed, fastVar not defined locally')
}
