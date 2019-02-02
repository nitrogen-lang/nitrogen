import "std/test"

test.run("Exceptions", func(assert) {
    let fastVar = 42

    const myException = func() {
        myException2()
    }

    const myException2 = func() {
        throw "Nope"
    }

    let m1 = try {
        myException()
    } catch e {
        errorVal(e)
    }

    assert.isEq(m1, "Nope")
    assert.isTrue(isDefined("fastVar"))
    assert.isFalse(isDefined("e"))
})
