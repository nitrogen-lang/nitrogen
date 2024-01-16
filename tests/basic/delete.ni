import "std/test"

test.run("delete variable", fn(assert) {
    let hello = "world"
    assert.isTrue(isDefined("hello"))

    delete hello
    assert.isFalse(isDefined("hello"))
})

test.run("delete constant", fn(assert) {
    assert.shouldRecover(fn() {
        const place = "Earth"
        delete place
    })
})
