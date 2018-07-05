import "stdlib/test"

test.run("Attempt to redefine constant", func(assert) {
    const thing = 42

    assert.shouldThrow(func() {
        thing = 43
    })
})

test.run("Variable in try goes out", func(assert) {
    try {
        let me_out = "please"
    } catch { pass }

    assert.isTrue(isDefined("me_out"))
})
