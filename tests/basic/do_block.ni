import "std/test"

test.run("Do block", fn(assert, check) {
    const t = do {
        1 + 2
    }

    check(assert.isEq(t, 3))
})

test.run("Do block scoping", fn(assert, check) {
    const t = do {
        const c = 6
        1 + 2
    }

    check(assert.isEq(t, 3))
    check(assert.isFalse(isDefined("c")))
})

test.run("Do block access outside scope", fn(assert, check) {
    const c = 6
    const t = do {
        1 + c
    }

    check(assert.isEq(t, 7))
})
