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
