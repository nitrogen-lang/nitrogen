import "std/test"

test.run("Attempt to redefine constant", fn(assert) {
    const thing = 42

    assert.shouldRecover(fn() {
        thing = 43
    })
})
