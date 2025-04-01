import "std/test"

const fac = fn(num) {
    if num == 0: return 1
    return num * fac(num - 1)
}

test.run("Factorial 20", fn(assert, check) {
    check(assert.isEq(fac(20), 2432902008176640000))
})
