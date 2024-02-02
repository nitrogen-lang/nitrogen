import "std/test"

const arr = ["one", "two"]

test.run("array prepend", fn(assert, check) {
    const prependArr = prepend(arr, "zero")
    check(assert.isEq(prependArr[0], "zero"))
})

test.run("array push", fn(assert, check) {
    const pushArr = push(arr, "three")
    check(assert.isEq(len(pushArr), 3))
    check(assert.isEq(pushArr[2], "three"))
})

test.run("array pop", fn(assert, check) {
    const popArr = pop(arr)
    check(assert.isEq(len(popArr), 1))
    check(assert.isEq(popArr[0], "one"))
})

const arr2 = arr + ["three", "four"]

test.run("array splice with offset and length", fn(assert, check) {
    const spliceArr = splice(arr2, 1, 2)

    check(assert.isEq(len(spliceArr), 2))
    check(assert.isEq(spliceArr[0], "one"))
    check(assert.isEq(spliceArr[1], "four"))
})

test.run("array splice with offset", fn(assert, check) {
    const spliceArr = splice(arr2, 2)

    check(assert.isEq(len(spliceArr), 2))
    check(assert.isEq(spliceArr[0], "one"))
    check(assert.isEq(spliceArr[1], "two"))
})

test.run("array splice 0 offset, no length", fn(assert, check) {
    const spliceArr = splice(arr2, 0)

    check(assert.isEq(len(spliceArr), 0))
})

test.run("array splice with negative offset and length", fn(assert, check) {
    check(assert.shouldRecover(fn() { splice(arr2, -1) }))
    check(assert.shouldRecover(fn() { splice(arr2, 1, -1) }))
    check(assert.shouldRecover(fn() { splice(arr2, -1, -1) }))
})

test.run("array splice with 0 length", fn(assert, check) {
    const spliceArr = splice(arr2, 1, 0)
    check(assert.isEq(len(spliceArr), 4))
})

test.run("array slice with 0 offset", fn(assert, check) {
    const sliceArr = slice(arr2, 0)
    check(assert.isEq(len(sliceArr), 4))
})

test.run("array slice with offset and length", fn(assert, check) {
    const sliceArr = slice(arr2, 1, 2)
    check(assert.isEq(len(sliceArr), 2))
    check(assert.isEq(sliceArr[0], "two"))
    check(assert.isEq(sliceArr[1], "three"))
})
