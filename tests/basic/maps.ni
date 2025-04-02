import "std/test"
import "std/collections" as col

test.run("Maps", fn(assert, check) {
    const hash = {
        "key1": "value1",
        "key2": "value2",
        "key3": "value3",
        "key4": "value4",
    }

    const keys = hashKeys(hash);
    const expectedKeys = ["key1", "key2", "key3", "key4"]

    check(assert.isEq(len(keys), len(expectedKeys)))

    col.foreach(expectedKeys, fn(i, v) {
        check(assert.isTrue(col.contains(keys, v)))
    })
})

test.run("Map ident keys", fn(assert, check) {
    const item1 = "Hello"

    const hash = {
        item1,
    }

    check(assert.isEq(hash.item1, item1))
    check(assert.isEq(hash["item1"], item1))
})
