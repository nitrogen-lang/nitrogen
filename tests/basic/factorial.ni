import "std/test"

const fac = fn(in) {
    if in == 0: return 1
    return in * fac(in - 1)
}

test.run("Factorial 20", fn(assert) {
    assert.isEq(fac(20), 2432902008176640000)
})
