import "test"
import "stdlib/opbuf"

test.run("Basic output buffering", func(assert) {
    opbuf.start()
    print("Hello")
    assert.isEq(opbuf.stopAndGet(), "Hello")

    assert.shouldThrow(func() {
        opbuf.stopAndGet()
    })
})

test.run("No double output buffering", func(assert) {
    opbuf.start()
    assert.shouldThrow(func() {
        opbuf.start()
    })
})

test.run("Clear output buffer", func(assert) {
    opbuf.start()
    print("Hello")
    opbuf.clear()
    print("Nitrogen")
    assert.isEq(opbuf.stopAndGet(), "Nitrogen")
})

test.run("Get output buffer", func(assert) {
    opbuf.start()
    print("Hello")
    assert.isEq(opbuf.get(), "Hello")

    assert.shouldThrow(func() {
        opbuf.start()
    })
})

test.run("Get output buffer", func(assert) {
    opbuf.start()
    print("Hello")
    assert.isEq(opbuf.stop(), nil)

    assert.shouldThrow(func() {
        opbuf.stop()
    })
})
