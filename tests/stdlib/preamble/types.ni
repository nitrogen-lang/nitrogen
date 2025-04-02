import "std/test"

test.run("toInt()", fn(assert, check) {
    const tests = [
        [23.5, 23],
        [1, 1],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(toInt(input), expected))
    }
})

test.run("toInt() error handling", fn(assert, check) {
    const tests = [
        [fn() { toInt("four") }, "Argument to `toInt` must be FLOAT, INT, or BYTESTRING, got STRING"],
        [fn() { toInt() }, "Incorrect number of arguments. Got 0, expected 1"],
        [fn() { toInt([]) }, "Argument to `toInt` must be FLOAT, INT, or BYTESTRING, got ARRAY"],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.shouldRecover(input, expected))
    }
})

test.run("toFloat()", fn(assert, check) {
    const tests = [
        [23.5, 23.5],
        [1, 1.0],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(toFloat(input), expected))
    }
})

test.run("toFloat() error handling", fn(assert, check) {
    const tests = [
        [fn() { toFloat("four") }, "Argument to `toFloat` must be FLOAT or INT, got STRING"],
        [fn() { toFloat() }, "Incorrect number of arguments. Got 0, expected 1"],
        [fn() { toFloat([]) }, "Argument to `toFloat` must be FLOAT or INT, got ARRAY"],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.shouldRecover(input, expected))
    }
})

test.run("builtin isX()", fn(assert, check) {
    const tests = [
        [fn() { isFloat(3.14159) }, true],
        [fn() { isFloat(3) }, false],
        [fn() { isInt(3) }, true],
        [fn() { isInt(3.14159) }, false],
        [fn() { isBool(true) }, true],
        [fn() { isBool(false) }, true],
        [fn() { isNull(nil) }, true],
        [fn() { isNull("nil") }, false],
        [fn() { isFunc(fn() { 10; }) }, true],
        [fn() { isFunc(10) }, false],
        [fn() { isString("Hello") }, true],
        [fn() { isString(10) }, false],
        [fn() { isArray([10, "true", false]) }, true],
        [fn() { isArray("array") }, false],
        [fn() { isMap({"key": "value"}) }, true],
        [fn() { isMap("array") }, false],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(input(), expected))
    }
})

test.run("toString", fn(assert, check) {
    const tests = [
        [1, "1"],
        [1.1, "1.1"],
        [true, "true"],
        [false, "false"],
        [nil, "nil"],
        ["hello", "hello"],
        [[1, 2], "[1, 2]"],
        [{"key": "value"}, "{key: \"value\"}"],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(toString(input), expected))
    }
})

test.run("parseInt", fn(assert, check) {
    const tests = [
        ["1", 1],
        ["1.1", nil],
        ["-1", -1],
        ["-1.1", nil],
        ["hello", nil],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(parseInt(input), expected), input)
    }
})

test.run("parseFloat", fn(assert, check) {
    const tests = [
        ["1", 1.0],
        ["1.1", 1.1],
        ["-1", -1.0],
        ["-1.1", -1.1],
        ["hello", nil],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(parseFloat(input), expected), input)
    }
})
