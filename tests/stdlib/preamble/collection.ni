import "std/test"

test.run("len()", fn(assert, check) {
    const tests = [
        ["", 0],
        ["four", 4],
        ["hello world", 11],
        [[1, 2, 3], 3],
        [[], 0],
        [nil, 0],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(len(input), expected))
    }
})

test.run("len() error handling", fn(assert, check) {
    const tests = [
        [fn() { len(1) }, "len(): Unsupported type INTEGER"],
        [fn() { len("one", "two") }, "Incorrect number of arguments. Got 2, expected 1"],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.shouldRecover(input, expected))
    }
})

test.run("first()", fn(assert, check) {
    const tests = [
        ["", nil],
        ["four", "f"],
        ["hello world", "h"],
        [[1, 2, 3], 1],
        [[], nil],
        [nil, nil],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(first(input), expected))
    }
})

test.run("first() error handling", fn(assert, check) {
    const tests = [
        [fn() { first(1) }, "Argument to `first` must be ARRAY|STRING|BYTESTRING, got INTEGER"],
        [fn() { first("one", "two") }, "Incorrect number of arguments. Got 2, expected 1"],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.shouldRecover(input, expected))
    }
})

test.run("last()", fn(assert, check) {
    const tests = [
        ["", nil],
        ["four", "r"],
        ["hello world", "d"],
        [[1, 2, 3], 3],
        [[], nil],
        [nil, nil],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(last(input), expected))
    }
})

test.run("last() error handling", fn(assert, check) {
    const tests = [
        [fn() { last(1) }, "Argument to `last` must be ARRAY|STRING|BYTESTRING, got INTEGER"],
        [fn() { last("one", "two") }, "Incorrect number of arguments. Got 2, expected 1"],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.shouldRecover(input, expected))
    }
})

test.run("rest()", fn(assert, check) {
    const tests = [
        [[1], []],
        [[1, 2, 3], [2, 3]],
        [[], nil],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.isEq(rest(input), expected))
    }
})

test.run("rest() error handling", fn(assert, check) {
    const tests = [
        [fn() { rest(1) }, "Argument to `rest` must be ARRAY, got INTEGER"],
        [fn() { rest("one", "two") }, "Incorrect number of arguments. Got 2, expected 1"],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.shouldRecover(input, expected))
    }
})

test.run("push()", fn(assert, check) {
    const tests = [
        [[1], 2, [1, 2]],
        [[1, 2, 3], 4, [1, 2, 3, 4]],
        [[], 1, [1]],
    ]

    for test in tests {
        const input = test[0]
        const add = test[1]
        const expected = test[2]
        check(assert.isEq(push(input, add), expected))
    }
})

test.run("push() error handling", fn(assert, check) {
    const tests = [
        [fn() { push("four", "five") }, "Argument to `push` must be ARRAY, got STRING"],
        [fn() { push() }, "Incorrect number of arguments. Got 0, expected 2"],
        [fn() { push([1]) }, "Incorrect number of arguments. Got 1, expected 2"],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.shouldRecover(input, expected))
    }
})

test.run("hashmerge()", fn(assert, check) {
    const test1 = hashMerge({"key2": "value"}, {"key2": "value2"})
    check(assert.isEq(test1.key2, "value2"))


    const test2 = hashMerge({"key2": "value"}, {"key2": "value2"}, false)
    check(assert.isEq(test2.key2, "value"))

    const test3 = hashMerge({"key": "value"}, {"key2": "value2"})
    check(assert.isEq(sort(hashKeys(test3)), ["key", "key2"]))
    check(assert.isEq(test3.key, "value"))
    check(assert.isEq(test3.key2, "value2"))
})

test.run("hashmerge() error handling", fn(assert, check) {
    const tests = [
        [fn() { hashMerge() }, "hashMerge requires at least 2 arguments. Got 0"],
        [fn() { hashMerge({"key": "value"}, 10) }, "First two arguments must be maps"],
        [fn() { hashMerge(10, {"key": "value"}) }, "First two arguments must be maps"],
    ]

    for test in tests {
        const input = test[0]
        const expected = test[1]
        check(assert.shouldRecover(input, expected))
    }
})
