import "std/test"

test.run("isTrue", fn(assert, check) {
    check(assert.isTrue(true))
})

test.run("isFalse", fn(assert, check) {
    check(assert.isFalse(false))
})

test.run("isEq", fn(assert, check) {
    check(assert.isEq(1, 1))
    check(assert.isEq("test", "test"))
    check(assert.isEq(nil, nil))
})

test.run("isNeq", fn(assert, check) {
    check(assert.isNeq(1, 2))
    check(assert.isNeq("test", "test2"))
    check(assert.isNeq(nil, 1))
    check(assert.isNeq("string", nil))
})
