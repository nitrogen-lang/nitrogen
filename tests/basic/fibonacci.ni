import "test"

func fib(x) {
    if x == 0 or x == 1: return x
    return fib(x-1) + fib(x-2)
}

test.run("Fibonacci 10", func(assert) {
    assert.isEq(fib(10), 55)
})
