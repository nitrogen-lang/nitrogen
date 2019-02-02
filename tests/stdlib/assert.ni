import "std/assert"

// This script doesn't use the "test" module since assert is used by test and that's what
// we're testing here. We're testing the test harness.

const shouldThrow = func(name, fn) {
    try { fn() } catch { return }

    throw "Assertion Failed: Expected " + name + " to throw."
}

const shouldNotThrow = func(name, fn) {
    try {
        fn()
    } catch {
        throw "Assertion Failed: Expected " + name + " not to throw."
    }
}

// isTrue
shouldNotThrow("001", func() {
    assert.isTrue(true)
})

shouldThrow("002", func() {
    assert.isTrue(false)
})

// isFalse
shouldNotThrow("003", func() {
    assert.isFalse(false)
})

shouldThrow("004", func() {
    assert.isFalse(true)
})

// isEq
shouldNotThrow("005", func() {
    assert.isEq("hello", "hello")
})

shouldThrow("006", func() {
    assert.isEq("hello", "Hello")
})

shouldThrow("007", func() {
    assert.isEq("hello", 42)
})

// isNeq
shouldNotThrow("008", func() {
    assert.isNeq("hello", "Hello")
})

shouldThrow("009", func() {
    assert.isNeq("hello", "hello")
})

shouldThrow("010", func() {
    assert.isNeq("hello", 42)
})

// shouldThrow
shouldNotThrow("011", func() {
    assert.shouldThrow(func() { throw "Hello" })
})

shouldThrow("012", func() {
    assert.shouldThrow(func() { 42 })
})

// shouldNotThrow
shouldNotThrow("013", func() {
    assert.shouldNotThrow(func() { 42 })
})

shouldThrow("014", func() {
    assert.shouldNotThrow(func() { throw "Hello" })
})
