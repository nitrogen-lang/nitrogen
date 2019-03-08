import "std/test"

test.run("Failed lookup in class init", fn(assert) {
    class MyClass {
        const init = fn() {
            println(things)
        }
    }

    assert.shouldThrow(fn() {
        const instance = new MyClass()
    })

    return nil // See bug regarding class init exceptions not propogating
})
