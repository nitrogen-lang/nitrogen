import "std/test"

test.run("Do block", fn(assert) {
    const t = do {
        1 + 2
    }

    assert.isEq(t, 3)
})

test.run("Do block scoping", fn(assert) {
    const t = do {
        const c = 6
        1 + 2
    }

    assert.isEq(t, 3)
    assert.isFalse(isDefined("c"))
})

test.run("Do block access outside scope", fn(assert) {
    const c = 6
    const t = do {
        1 + c
    }

    assert.isEq(t, 7)
})
