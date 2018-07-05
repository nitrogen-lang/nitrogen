import "stdlib/test"
import "stdlib/opbuf"

test.run("Stop output buffer", func(assert) {
    opbuf.start()
    print("Hello")
    assert.isEq(opbuf.stop(), nil)

    assert.shouldThrow(func() {
        opbuf.stop()
    })
})

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
    opbuf.stop()
})

test.run("No double output buffering", func(assert) {
    opbuf.start()
    assert.isTrue(opbuf.isStarted())
    opbuf.stop()
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
