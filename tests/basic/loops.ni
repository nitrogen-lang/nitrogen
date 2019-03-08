import "std/test"
import "std/collections" as col

test.run("Simple loop", fn(assert) {
    let outer = 0

    for (i = 0; i < 10; i += 1) {
        outer = outer + 1
    }

    assert.isEq(outer, 10)
})

test.run("Loop with continue", fn(assert) {
    let outer = 0

    for (i = 0; i < 10; i += 1) {
        if (i % 2 > 0): continue
        outer = outer + 1
    }

    assert.isEq(outer, 5)
})

test.run("Loop with break", fn(assert) {
    let outer = 0

    for (i = 0; i < 12; i += 1) {
        if (i > 9): break
        outer = outer + 1
    }

    assert.isEq(outer, 10)
})

test.run("While loop", fn(assert) {
    const testWhile = fn() {
        let finished = false
        let i = 0

        while !finished {
            i += 1
            if i == 5: finished = true
        }

        i
    }

    assert.isEq(testWhile(), 5)
})

test.run("Infinite loop", fn(assert) {
    let finished = false
    let i = 0

    loop {
        if finished: break
        i += 1
        if i == 5: finished = true
    }

    assert.isEq(i, 5)
})

test.run("Array iterator", fn(assert) {
    let sum = 0

    const nums = [2, 5, 10, 12, 5, 7]

    for num in nums {
        sum += num
    }

    assert.isEq(sum, 41)
})

test.run("Array iterator with index", fn(assert) {
    let last_i = 0

    const nums = [2, 5, 10, 12, 5, 7]

    for i, num in nums {
        last_i = i
    }

    assert.isEq(last_i, 5)
})

test.run("Hashmap iterator", fn(assert) {
    let vals = []

    const map = {
        "key1": "val1",
        "key2": "val2",
        "key3": "val3",
    }

    for item in map {
        vals = push(vals, item)
    }

    vals = sort(vals)

    assert.isTrue(col.arrayMatch(vals, ["val1", "val2", "val3"]))
})

test.run("Hashmap iterator with keys", fn(assert) {
    let keys = []

    const map = {
        "key1": "val1",
        "key2": "val2",
        "key3": "val3",
    }

    for key, item in map {
        keys = push(keys, key)
    }

    keys = sort(keys)

    assert.isTrue(col.arrayMatch(keys, ["key1", "key2", "key3"]))
})
