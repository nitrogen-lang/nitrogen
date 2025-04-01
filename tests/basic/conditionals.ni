import "std/test"

test.run("Single statement if expressions", fn(assert, check) {
    let counter = 0

    if true {
        counter += 1
    }

    if (true) {
        counter += 1
    }

    if true: counter += 1
    if (true): counter += 1

    check(assert.isEq(counter, 4))
})

test.run("Various if expressions", fn(assert, check) {
    check(assert.shouldRecover(fn() {
        if ["a"] == ["b"]: return
        println("Hello")
    }))

    check(assert.isTrue(fn() {
        if 42 == 42: return true
        return false
    }()))

    check(assert.isFalse(fn() {
        if "42" == 42: return true
        return false
    }))
})

test.run("If statement with else", fn(assert, check) {
    const theTest = fn(a) {
        if a {
            return "Hello"
        } else {
            return "Good bye"
        }
    }

    check(assert.isEq(theTest(true), "Hello"))
    check(assert.isEq(theTest(false), "Good bye"))
})

test.run("If statement with else and compound conditional", fn(assert, check) {
    const theTest = fn(a, b) {
        if a or b {
            return "Hello"
        } else {
            return "Good bye"
        }
    }

    check(assert.isEq(theTest(true, false), "Hello"))
    check(assert.isEq(theTest(false, true), "Hello"))
    check(assert.isEq(theTest(true, true), "Hello"))
    check(assert.isEq(theTest(false, false), "Good bye"))
})

test.run("If statement with else assigned to variable", fn(assert, check) {
    const theTest = fn(a) {
        const msg = if a {
            "Hello"
        } else {
            "Good bye"
        }

        return msg
    }

    check(assert.isEq(theTest(true), "Hello"))
    check(assert.isEq(theTest(false), "Good bye"))
})

test.run("If statement within an if true branch", fn(assert, check) {
    const theTest = fn(a, b) {
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

    check(assert.isEq(theTest(true, true), "Hello1"))
    check(assert.isEq(theTest(true, false), "Hello2"))
    check(assert.isEq(theTest(false, true), "Good bye"))
    check(assert.isEq(theTest(false, false), "Good bye"))
})

test.run("If statement within an if false branch", fn(assert, check) {
    const theTest = fn(a, b) {
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

    check(assert.isEq(theTest(true, true), "Hello"))
    check(assert.isEq(theTest(true, false), "Hello"))
    check(assert.isEq(theTest(false, true), "Good bye1"))
    check(assert.isEq(theTest(false, false), "Good bye2"))
})

test.run("If statement with elif block", fn(assert, check) {
    const theTest = fn(a) {
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

    check(assert.isEq(theTest(42), 82))
    check(assert.isEq(theTest(43), 83))
    check(assert.isEq(theTest(44), 84))
    check(assert.isEq(theTest(45), 90))
})
