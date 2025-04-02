import "std/test"
import "std/encoding/json"
import "std/collections" as col

let tests = [
    {
        "in": "Hello, World",
        "out": '"Hello, World"',
    },
    {
        "in": 42,
        "out": '42',
    },
    {
        "in": 42.5,
        "out": '42.5',
    },
    {
        "in": true,
        "out": 'true',
    },
    {
        "in": ["Hello", 123, nil],
        "out": '["Hello",123,null]',
    },
    {
        "in": {"key1": "val1"},
        "out": '{"key1":"val1"}',
    },
    {
        "in": ["Hello", [true, false]],
        "out": '["Hello",[true,false]]',
    },
    {
        "in": ["Hello", {"key1": "key2"}],
        "out": '["Hello",{"key1":"key2"}]',
    },
    {
        "in": ["Hello", {"key1": "key2"}, 42, [true, false]],
        "out": '["Hello",{"key1":"key2"},42,[true,false]]',
    },
    {
        "in": {"key1": ["Hello", 42, true, nil]},
        "out": '{"key1":["Hello",42,true,null]}',
    },
]

test.run("JSON encode", fn(assert, check) {
    col.foreach(tests, fn(i, el) {
        check(assert.isEq(json.encode(el["in"]), el["out"]))
    })
})

test.run("JSON encode bad value", fn(assert, check) {
    check(assert.shouldRecover(fn() {
        json.encode(fn() {pass})
    }))
})

test.run("JSON decode", fn(assert, check) {
    col.foreach(tests, fn(i, el) {
        const decoded = json.decode(el["out"])

        if isArray(el["in"]){
            check(assert.isTrue(col.arrayMatch(decoded, el["in"])), el["out"])
        } elif isMap(el["in"]) {
            check(assert.isTrue(col.mapMatch(decoded, el["in"])), el["out"])
        } else {
            check(assert.isEq(decoded, el["in"]), el["out"])
        }
    })

    const whitespaceTest = '{
    "title": "foo",
    "userId": 1,
    "body": "bar",
    "id": 101
}'

    const wsExpected = {
        "title": "foo",
        "userId": 1,
        "body": "bar",
        "id": 101,
    }

    const decoded = json.decode(whitespaceTest)
    check(assert.isTrue(col.mapMatch(decoded, wsExpected)), "map match")
})
