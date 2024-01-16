import "std/test"

const arr = ["one", "two"]

test.run("array prepend", fn(assert) {
    const prependArr = prepend(arr, "zero")
    assert.isEq(prependArr[0], "zero")
})

test.run("array push", fn(assert) {
    const pushArr = push(arr, "three")
    assert.isEq(len(pushArr), 3)
    assert.isEq(pushArr[2], "three")
})

test.run("array pop", fn(assert) {
    const popArr = pop(arr)
    assert.isEq(len(popArr), 1)
    assert.isEq(popArr[0], "one")
})

const arr2 = arr + ["three", "four"]

test.run("array splice with offset and length", fn(assert) {
    const spliceArr = splice(arr2, 1, 2)

    assert.isEq(len(spliceArr), 2)
    assert.isEq(spliceArr[0], "one")
    assert.isEq(spliceArr[1], "four")
})

test.run("array splice with offset", fn(assert) {
    const spliceArr = splice(arr2, 2)

    assert.isEq(len(spliceArr), 2)
    assert.isEq(spliceArr[0], "one")
    assert.isEq(spliceArr[1], "two")
})

test.run("array splice 0 offset, no length", fn(assert) {
    const spliceArr = splice(arr2, 0)

    assert.isEq(len(spliceArr), 0)
})

test.run("array splice with negative offset and length", fn(assert) {
    assert.shouldRecover(fn() {
        splice(arr2, -1)
    })

    assert.shouldRecover(fn() {
        splice(arr2, 1, -1)
    })

    assert.shouldRecover(fn() {
        splice(arr2, -1, -1)
    })
})

test.run("array splice with 0 length", fn(assert) {
    const spliceArr = splice(arr2, 1, 0)
    assert.isEq(len(spliceArr), 4)
})

test.run("array slice with 0 offset", fn(assert) {
    const sliceArr = slice(arr2, 0)
    assert.isEq(len(sliceArr), 4)
})

test.run("array slice with offset and length", fn(assert) {
    const sliceArr = slice(arr2, 1, 2)
    assert.isEq(len(sliceArr), 2)
    assert.isEq(sliceArr[0], "two")
    assert.isEq(sliceArr[1], "three")
})
