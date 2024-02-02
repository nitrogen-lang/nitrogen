import "std/test"

test.run("Compound equality", fn(assert, check) {
    const a = 5
    const b = 6
    const c = 5

    check(assert.isFalse(a >= b))
    check(assert.isFalse(b <= a))
    check(assert.isTrue(a <= c))
    check(assert.isTrue(a >= c))
})

test.run("Compound assignment", fn(assert, check) {
    let a = 5

    a += 2
    check(assert.isEq(a, 7))

    a -= 3
    check(assert.isEq(a, 4))

    a *= 2
    check(assert.isEq(a, 8))

    a /= 4
    check(assert.isEq(a, 2))
})
