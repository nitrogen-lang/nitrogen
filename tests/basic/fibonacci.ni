import "std/test"

const fib = fn(x) {
    if x == 0 or x == 1: return x
    return fib(x-1) + fib(x-2)
}

test.run("Fibonacci 10", fn(assert, check) {
    check(assert.isEq(fib(10), 55))
})
