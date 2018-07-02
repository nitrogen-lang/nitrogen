import "test"

test.run("Basic environment", func(assert) {
    let string = "Hello, world!"

    // Ensure this function changes the outer scope variable
    func change_string(next) {
        string = next
    }

    // Ensure this one doesn't
    func not_change_string(next) {
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
