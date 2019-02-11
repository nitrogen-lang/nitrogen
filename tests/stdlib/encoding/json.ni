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

test.run("JSON encode", fn(assert) {
    col.foreach(tests, fn(i, el) {
        assert.isEq(json.encode(el.in), el.out)
    })
})

test.run("JSON encode bad value", fn(assert) {
    assert.shouldThrow(fn() {
        json.encode(fn() {pass})
    })
})

test.run("JSON decode", fn(assert) {
    col.foreach(tests, fn(i, el) {
        const decoded = json.decode(el.out)

        if isArray(el.in){
            assert.isTrue(col.arrayMatch(decoded, el.in))
        } elif isMap(el.in) {
            assert.isTrue(col.mapMatch(decoded, el.in))
        } else {
            assert.isEq(decoded, el.in)
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
    assert.isTrue(col.mapMatch(decoded, wsExpected))
})
