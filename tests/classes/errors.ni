import "std/test"

// Issue #37
test.run("Failed lookup in class init", fn(assert) {
    class MyClass {
        fn init() {
            println(thing)
        }
    }

    fn main() {
        recover {
            const things = new MyClass()
        }

        42 // Ensure the instance object doesn't linger on the stack
    }

    assert.isEq(main(), 42)
})
