import "test"
import "json"
import "collections" as col

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
]

test.run("JSON encode", func(assert) {
    col.foreach(tests, func(i, el) {
        assert.isEq(json.encode(el.in), el.out)
    })
})

test.run("JSON encode bad value", func(assert) {
    assert.shouldThrow(func() {
        json.encode(func() {pass})
    })
})
