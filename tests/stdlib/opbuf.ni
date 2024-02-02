import "std/test"
import "std/opbuf"

test.run("Stop output buffer", fn(assert, check) {
    opbuf.start()
    print("Hello")
    check(assert.isEq(opbuf.stop(), nil))

    check(assert.shouldRecover(fn() {
        opbuf.stop()
    }))
})

test.run("Basic output buffering", fn(assert, check) {
    opbuf.start()
    print("Hello")
    check(assert.isEq(opbuf.stopAndGet(), "Hello"))

    check(assert.shouldRecover(fn() {
        opbuf.stopAndGet()
    }))
})

test.run("No double output buffering", fn(assert, check) {
    opbuf.start()
    check(assert.shouldRecover(fn() {
        opbuf.start()
    }))
    opbuf.stop()
})

test.run("No double output buffering", fn(assert, check) {
    opbuf.start()
    check(assert.isTrue(opbuf.isStarted()))
    opbuf.stop()
})

test.run("Clear output buffer", fn(assert, check) {
    opbuf.start()
    print("Hello")
    opbuf.clear()
    print("Nitrogen")
    check(assert.isEq(opbuf.stopAndGet(), "Nitrogen"))
})

test.run("Get output buffer", fn(assert, check) {
    opbuf.start()
    print("Hello")
    check(assert.isEq(opbuf.get(), "Hello"))

    check(assert.shouldRecover(fn() {
        opbuf.start()
    }))
})
