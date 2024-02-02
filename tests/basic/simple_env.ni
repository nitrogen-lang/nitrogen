import "std/test"

test.run("Basic environment", fn(assert, check) {
    let string = "Hello, world!"

    // Ensure this fntion changes the outer scope variable
    const change_string = fn(next) {
        string = next
    }

    // Ensure this one doesn't
    const not_change_string = fn(next) {
        // Overshadow the global string variable
        let string = ""
        string = next
    }

    check(assert.isEq(string, "Hello, world!"))

    change_string("Hello, mars!")

    check(assert.isEq(string, "Hello, mars!"))

    not_change_string("Hello, earth!")

    check(assert.isEq(string, "Hello, mars!"))
})
