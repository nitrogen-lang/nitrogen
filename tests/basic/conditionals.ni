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
