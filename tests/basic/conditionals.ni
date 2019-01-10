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
    func test(a) {
        if a {
            return "Hello"
        } else {
            return "Good bye"
        }
    }

    assert.isEq(test(true), "Hello")
    assert.isEq(test(false), "Good bye")
})

test.run("If statement with else and compound conditional", func(assert) {
    func test(a, b) {
        if a or b {
            return "Hello"
        } else {
            return "Good bye"
        }
    }

    assert.isEq(test(true, false), "Hello")
    assert.isEq(test(false, true), "Hello")
    assert.isEq(test(true, true), "Hello")
    assert.isEq(test(false, false), "Good bye")
})

test.run("If statement with else assigned to variable", func(assert) {
    func test(a) {
        const msg = if a {
            "Hello"
        } else {
            "Good bye"
        }

        return msg
    }

    assert.isEq(test(true), "Hello")
    assert.isEq(test(false), "Good bye")
})

test.run("If statement within an if true branch", func(assert) {
    func test(a, b) {
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

    assert.isEq(test(true, true), "Hello1")
    assert.isEq(test(true, false), "Hello2")
    assert.isEq(test(false, true), "Good bye")
    assert.isEq(test(false, false), "Good bye")
})

test.run("If statement within an if false branch", func(assert) {
    func test(a, b) {
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

    assert.isEq(test(true, true), "Hello")
    assert.isEq(test(true, false), "Hello")
    assert.isEq(test(false, true), "Good bye1")
    assert.isEq(test(false, false), "Good bye2")
})

test.run("If statement with elif block", func(assert) {
    func test(a) {
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

    assert.isEq(test(42), 82)
    assert.isEq(test(43), 83)
    assert.isEq(test(44), 84)
    assert.isEq(test(45), 90)
})
