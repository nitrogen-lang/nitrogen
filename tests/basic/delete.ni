import "std/test"

test.run("delete variable", func(assert) {
    let hello = "world"
    assert.isTrue(isDefined("hello"))

    delete hello
    assert.isFalse(isDefined("hello"))
})

test.run("delete constant", func(assert) {
    assert.shouldThrow(func() {
        const place = "Earth"
        delete place
    })
})
