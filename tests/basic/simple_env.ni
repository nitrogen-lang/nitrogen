import "std/test"

test.run("Basic environment", fn(assert) {
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

    assert.isEq(string, "Hello, world!")

    change_string("Hello, mars!")

    assert.isEq(string, "Hello, mars!")

    not_change_string("Hello, earth!")

    assert.isEq(string, "Hello, mars!")
})
