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

if m1 != "Nope" {
    println('m1 is not the correct value')
    println('Expected "Nope", got "', m1, '"')
}
