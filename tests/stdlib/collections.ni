import "stdlib/test"
import "stdlib/collections" as col

const testArr = ["asia", "north america", "south america", "africa", "europe", "australia", "antarctica"]

test.run("collections array match", func(assert) {
    const testArr2 = ["asia", "north america", "south america", "africa", "europe", "australia", "antarctica"]
    const testArr3 = ["asia", "north america", "south america", "africa", "europe", "antarctica", "australia"]
    assert.isTrue(col.arrayMatch(testArr, testArr2))
    assert.isFalse(col.arrayMatch(testArr, testArr3))
})

test.run("collections filter", func(assert) {
    const expected = ["asia", "africa", "australia", "antarctica"]
    const filterTest = col.filter(testArr, func(i){ i[0] == "a" })
    assert.isTrue(col.arrayMatch(filterTest, expected))
})

test.run("collections map", func(assert) {
    const expected = ["C-asia", "C-north america", "C-south america", "C-africa", "C-europe", "C-australia", "C-antarctica"]
    const mapTest = col.map(testArr, func(i){ "C-"+i })
    assert.isTrue(col.arrayMatch(mapTest, expected))
})

test.run("collections reduce", func(assert) {
    const expected = 61
    const reduceTest = col.reduce(testArr, func(r, e){ r + len(e) }, 0)
    assert.isEq(reduceTest, expected)
})

test.run("collections foreach with map", func(assert) {
    const expected = {
        "key1": "value1",
        "key2": "value2",
        "key3": "value3",
    }

    const foreachTest = {}
    col.foreach(expected, func(k, v) {
        foreachTest[k] = v
    })
    assert.isTrue(col.mapMatch(foreachTest, expected))
})

test.run("collections foreach with array", func(assert) {
    let foreachArrTest = []
    col.foreach(testArr, func(i, e) {
        foreachArrTest = push(foreachArrTest, e)
    })
    assert.isTrue(col.arrayMatch(foreachArrTest, testArr))
})

test.run("collections foreach with string", func(assert) {
    const testStr = "Hello"
    let counter = 0

    col.foreach(testStr, func(i, c) {
        counter += 1
    })

    assert.isEq(counter, 5)
})

test.run("collections foreach not collection", func(assert) {
    assert.shouldThrow(func() {
        col.foreach(42, func(){})
    })
})

test.run("collections map match", func(assert) {
    let map1 = {
        "key1": "value1",
        "key2": "value2",
        "key3": "value3",
    }
    let map2 = {
        "key1": "value1",
        "key2": "value2",
        "key3": "value3",
    }
    assert.isTrue(col.mapMatch(map1, map2))

    map1 = {
        "key1": ["value1", "value1.1", "value1.2"],
        "key2": "value2",
        "key3": "value3",
    }
    map2 = {
        "key1": ["value1", "value1.1", "value1.2"],
        "key2": "value2",
        "key3": "value3",
    }
    assert.isTrue(col.mapMatch(map1, map2))

    map1 = {
        "key1": ["value8", "value1.1", "value1.2"],
        "key2": "value2",
        "key3": "value3",
    }
    map2 = {
        "key1": ["value1", "value1.1", "value1.2"],
        "key2": "value2",
        "key5": "value3",
    }
    assert.isFalse(col.mapMatch(map1, map2))
})
