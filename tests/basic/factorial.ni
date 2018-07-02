import "test"

func fac(in) {
    if in == 0: return 1
    return in * fac(in - 1)
}

test.run("Factorial 20", func(assert) {
    assert.isEq(fac(20), 2432902008176640000)
})
