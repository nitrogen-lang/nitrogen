import "std/test"
import "std/collections" as col

const testArr = ["asia", "north america", "south america", "africa", "europe", "australia", "antarctica"]

test.run("collections array match", fn(assert) {
    const testArr2 = ["asia", "north america", "south america", "africa", "europe", "australia", "antarctica"]
    const testArr3 = ["asia", "north america", "south america", "africa", "europe", "antarctica", "australia"]
    assert.isTrue(col.arrayMatch(testArr, testArr2))
    assert.isFalse(col.arrayMatch(testArr, testArr3))
})

test.run("collections filter", fn(assert) {
    const expected = ["asia", "africa", "australia", "antarctica"]
    const filterTest = col.filter(testArr, fn(i){ i[0] == "a" })
    assert.isTrue(col.arrayMatch(filterTest, expected))
})

test.run("collections map", fn(assert) {
    const expected = ["C-asia", "C-north america", "C-south america", "C-africa", "C-europe", "C-australia", "C-antarctica"]
    const mapTest = col.map(testArr, fn(i){ "C-"+i })
    assert.isTrue(col.arrayMatch(mapTest, expected))
})

test.run("collections reduce array", fn(assert) {
    const expected = 61
    const reduceTest = col.reduce(testArr, fn(r, e){ r + len(e) }, 0)
    assert.isEq(reduceTest, expected)
})

test.run("collections reduce map", fn(assert) {
    const theMap = {
        "key1": 1,
        "key2": 50,
        "key3": 10,
    }
    const expected = 61
    const reduceTest = col.reduce(theMap, fn(acc, val){ acc + val }, 0)
    assert.isEq(reduceTest, expected)
})

test.run("collections foreach with map", fn(assert) {
    const expected = {
        "key1": "value1",
        "key2": "value2",
        "key3": "value3",
    }

    const foreachTest = {}
    col.foreach(expected, fn(k, v) {
        foreachTest[k] = v
    })
    assert.isTrue(col.mapMatch(foreachTest, expected))
})

test.run("collections foreach with array", fn(assert) {
    let foreachArrTest = []
    col.foreach(testArr, fn(i, e) {
        foreachArrTest = push(foreachArrTest, e)
    })
    assert.isTrue(col.arrayMatch(foreachArrTest, testArr))
})

test.run("collections foreach with string", fn(assert) {
    const testStr = "Hello"
    let counter = 0

    col.foreach(testStr, fn(i, c) {
        counter += 1
    })

    assert.isEq(counter, 5)
})

test.run("collections foreach not collection", fn(assert) {
    assert.shouldThrow(fn() {
        col.foreach(42, fn(){})
    })
})

test.run("collections map match", fn(assert) {
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
