import "std/test"

test.run("Exceptions", fn(assert) {
    let fastVar = 42

    const myException = fn() {
        myException2()
    }

    const myException2 = fn() {
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
