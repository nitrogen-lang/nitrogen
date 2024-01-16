import "std/test"
import "std/opbuf"

test.run("Stop output buffer", fn(assert) {
    opbuf.start()
    print("Hello")
    assert.isEq(opbuf.stop(), nil)

    assert.shouldRecover(fn() {
        opbuf.stop()
    })
})

test.run("Basic output buffering", fn(assert) {
    opbuf.start()
    print("Hello")
    assert.isEq(opbuf.stopAndGet(), "Hello")

    assert.shouldRecover(fn() {
        opbuf.stopAndGet()
    })
})

test.run("No double output buffering", fn(assert) {
    opbuf.start()
    assert.shouldRecover(fn() {
        opbuf.start()
    })
    opbuf.stop()
})

test.run("No double output buffering", fn(assert) {
    opbuf.start()
    assert.isTrue(opbuf.isStarted())
    opbuf.stop()
})

test.run("Clear output buffer", fn(assert) {
    opbuf.start()
    print("Hello")
    opbuf.clear()
    print("Nitrogen")
    assert.isEq(opbuf.stopAndGet(), "Nitrogen")
})

test.run("Get output buffer", fn(assert) {
    opbuf.start()
    print("Hello")
    assert.isEq(opbuf.get(), "Hello")

    assert.shouldRecover(fn() {
        opbuf.start()
    })
})
