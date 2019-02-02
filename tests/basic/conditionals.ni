import "std/test"

test.run("Single statement if expressions", func(assert) {
    let counter = 0

    if true {
        counter += 1
    }

    if (true) {
        counter += 1
    }

    if true: counter += 1
    if (true): counter += 1

    assert.isEq(counter, 4)
})

test.run("Various if expressions", func(assert) {
    assert.shouldThrow(func() {
        if ["a"] == ["b"]: return
        println("Hello")
    })

    assert.isTrue(func() {
        if 42 == 42: return true
        return false
    }())

    assert.isFalse(func() {
        if "42" == 42: return true
        return false
    })
})

test.run("If statement with else", func(assert) {
    const theTest = func(a) {
        if a {
            return "Hello"
        } else {
            return "Good bye"
        }
    }

    assert.isEq(theTest(true), "Hello")
    assert.isEq(theTest(false), "Good bye")
})

test.run("If statement with else and compound conditional", func(assert) {
    const theTest = func(a, b) {
        if a or b {
            return "Hello"
        } else {
            return "Good bye"
        }
    }

    assert.isEq(theTest(true, false), "Hello")
    assert.isEq(theTest(false, true), "Hello")
    assert.isEq(theTest(true, true), "Hello")
    assert.isEq(theTest(false, false), "Good bye")
})

test.run("If statement with else assigned to variable", func(assert) {
    const theTest = func(a) {
        const msg = if a {
            "Hello"
        } else {
            "Good bye"
        }

        return msg
    }

    assert.isEq(theTest(true), "Hello")
    assert.isEq(theTest(false), "Good bye")
})

test.run("If statement within an if true branch", func(assert) {
    const theTest = func(a, b) {
        if a {
            if b {
                return "Hello1"
            } else {
                return "Hello2"
            }
        } else {
            return "Good bye"
        }
    }

    assert.isEq(theTest(true, true), "Hello1")
    assert.isEq(theTest(true, false), "Hello2")
    assert.isEq(theTest(false, true), "Good bye")
    assert.isEq(theTest(false, false), "Good bye")
})

test.run("If statement within an if false branch", func(assert) {
    const theTest = func(a, b) {
        if a {
            return "Hello"
        } else {
            if b {
                return "Good bye1"
            } else {
                return "Good bye2"
            }
        }
    }

    assert.isEq(theTest(true, true), "Hello")
    assert.isEq(theTest(true, false), "Hello")
    assert.isEq(theTest(false, true), "Good bye1")
    assert.isEq(theTest(false, false), "Good bye2")
})

test.run("If statement with elif block", func(assert) {
    const theTest = func(a) {
        if a == 42 {
            return 82
        } elif a == 43 {
            return 83
        } elif a == 44 {
            return 84
        } else {
            return 90
        }
    }

    assert.isEq(theTest(42), 82)
    assert.isEq(theTest(43), 83)
    assert.isEq(theTest(44), 84)
    assert.isEq(theTest(45), 90)
})
