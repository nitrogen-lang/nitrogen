import "std/test"

test.run("Attempt to redefine constant", fn(assert, check) {
    const thing = 42
    check(assert.shouldRecover(fn() { thing = 43 }))
})
