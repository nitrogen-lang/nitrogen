import "std/assert"

// This script doesn't use the "test" module since assert is used by test and that's what
// we're testing here. We're testing the test harness.

const shouldFail = fn(name, func) {
    if isError(recover { func() }): return

    println("Assertion Failed: Expected " + name + " to fail.")
}

const shouldPass = fn(name, func) {
    if !isError(recover { func() }): return
    println("Assertion Failed: Expected " + name + " to pass.")
}

// isTrue
shouldPass("001", fn() {
    assert.isTrue(true)
})

shouldFail("002", fn() {
    assert.isTrue(false)
})

// isFalse
shouldPass("003", fn() {
    assert.isFalse(false)
})

shouldFail("004", fn() {
    assert.isFalse(true)
})

// isEq
shouldPass("005", fn() {
    assert.isEq("hello", "hello")
})

shouldFail("006", fn() {
    assert.isEq("hello", "Hello")
})

shouldFail("007", fn() {
    assert.isEq("hello", 42)
})

// isNeq
shouldPass("008", fn() {
    assert.isNeq("hello", "Hello")
})

shouldFail("009", fn() {
    assert.isNeq("hello", "hello")
})

shouldFail("010", fn() {
    assert.isNeq("hello", 42)
})
