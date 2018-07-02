import "test"

test.run("Simple loop", func(assert) {
    let outer = 0

    for (i = 0; i < 10; i += 1) {
        outer = outer + 1
    }

    assert.isEq(outer, 10)
})

test.run("Loop with continue", func(assert) {
    let outer = 0

    for (i = 0; i < 10; i += 1) {
        if (i % 2 > 0) { continue }
        outer = outer + 1
    }

    assert.isEq(outer, 5)
})

test.run("Loop with break", func(assert) {
    let outer = 0

    for (i = 0; i < 12; i += 1) {
        if (i > 9) { break }
        outer = outer + 1
    }

    assert.isEq(outer, 10)
})
