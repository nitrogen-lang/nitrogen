import "std/assert"

// This script doesn't use the "test" module since assert is used by test and that's what
// we're testing here. We're testing the test harness.

const shouldThrow = fn(name, func) {
    try { func() } catch { return }

    throw "Assertion Failed: Expected " + name + " to throw."
}

const shouldNotThrow = fn(name, func) {
    try {
        func()
    } catch {
        throw "Assertion Failed: Expected " + name + " not to throw."
    }
}

// isTrue
shouldNotThrow("001", fn() {
    assert.isTrue(true)
})

shouldThrow("002", fn() {
    assert.isTrue(false)
})

// isFalse
shouldNotThrow("003", fn() {
    assert.isFalse(false)
})

shouldThrow("004", fn() {
    assert.isFalse(true)
})

// isEq
shouldNotThrow("005", fn() {
    assert.isEq("hello", "hello")
})

shouldThrow("006", fn() {
    assert.isEq("hello", "Hello")
})

shouldThrow("007", fn() {
    assert.isEq("hello", 42)
})

// isNeq
shouldNotThrow("008", fn() {
    assert.isNeq("hello", "Hello")
})

shouldThrow("009", fn() {
    assert.isNeq("hello", "hello")
})

shouldThrow("010", fn() {
    assert.isNeq("hello", 42)
})

// shouldThrow
shouldNotThrow("011", fn() {
    assert.shouldThrow(fn() { throw "Hello" })
})

shouldThrow("012", fn() {
    assert.shouldThrow(fn() { 42 })
})

// shouldNotThrow
shouldNotThrow("013", fn() {
    assert.shouldNotThrow(fn() { 42 })
})

shouldThrow("014", fn() {
    assert.shouldNotThrow(fn() { throw "Hello" })
})
