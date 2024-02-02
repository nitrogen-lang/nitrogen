import "std/test"

test.run("delete variable", fn(assert, check) {
    let hello = "world"
    check(assert.isTrue(isDefined("hello")))

    delete hello
    check(assert.isFalse(isDefined("hello")))
})

test.run("delete constant", fn(assert, check) {
    check(assert.shouldRecover(fn() {
        const place = "Earth"
        delete place
    }))
})
