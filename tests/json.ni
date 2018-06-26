import "json.ni"
import "collections.ni" as col

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

col.foreach(tests, func(i, el) {
    const encoded = json.encode(el.in)
    if encoded != el.out {
        println("JSON Test Failed: Expected ", el.out, ", got ", encoded)
    }
})

try {
    json.encode(func() {pass})
    println("JSON Test Failed: Expected encode to throw when given func")
} catch { pass }
