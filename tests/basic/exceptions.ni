import "test"

test.run("Exceptions", func(assert) {
    let fastVar = 42

    func myException() {
        myException2()
    }

    func myException2() {
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
