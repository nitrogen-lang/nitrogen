import "test"

test.run("Compound equality", func(assert) {
    const a = 5
    const b = 6
    const c = 5

    assert.isFalse(a >= b)
    assert.isFalse(b <= a)
    assert.isTrue(a <= c)
    assert.isTrue(a >= c)
})

test.run("Compound assignment", func(assert) {
    let a = 5

    a += 2
    assert.isEq(a, 7)

    a -= 3
    assert.isEq(a, 4)

    a *= 2
    assert.isEq(a, 8)

    a /= 4
    assert.isEq(a, 2)
})
