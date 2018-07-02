import "test"

test.run("Single statement if expressions", func(assert) {
    let counter = 0

    if true {
        counter += 1
    }

    if (true) {
        counter += 1
    }

    if true: counter += 1
    if (true): counter += 1

    assert.isEq(counter, 4)
})

test.run("Various if expressions", func(assert) {
    assert.shouldThrow(func() {
        if ["a"] == ["b"]: return
        println("Hello")
    })

    assert.isTrue(func() {
        if 42 == 42: return true
        return false
    }())

    assert.isFalse(func() {
        if "42" == 42: return true
        return false
    })
})
