import "std/test"

test.run("Attempt to redefine constant", fn(assert) {
    const thing = 42

    assert.shouldThrow(fn() {
        thing = 43
    })
})

test.run("Variable in try goes out", fn(assert) {
    try {
        let me_out = "please"
    } catch { pass }

    assert.isTrue(isDefined("me_out"))
})
