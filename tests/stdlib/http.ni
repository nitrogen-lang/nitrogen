import "stdlib/http"
import "stdlib/encoding/json"
import "stdlib/test"
import "stdlib/collections"

use collections.contains

test.run("HTTP GET request", func(assert) {
    const resp = http.get("https://jsonplaceholder.typicode.com/posts/1")

    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "headers"))
    assert.isTrue(isString(resp.body))
    assert.isNeq(resp.body, "")
    assert.isEq(resp.headers["Content-Type"], "application/json; charset=utf-8")
})

test.run("HTTP GET JSON request", func(assert) {
    const resp = http.getJSON("https://jsonplaceholder.typicode.com/posts/1")

    assert.isTrue(isMap(resp))
    assert.isTrue(contains(resp, "id"))
    assert.isTrue(contains(resp, "title"))
    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "userId"))
})

test.run("HTTP HEAD request", func(assert) {
    const resp = http.head("https://jsonplaceholder.typicode.com/posts/1")

    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "headers"))
    assert.isEq(resp.body, "")
    assert.isEq(resp.headers["Content-Type"], "application/json; charset=utf-8")
})

test.run("HTTP DELETE request", func(assert) {
    const resp = http.del("https://jsonplaceholder.typicode.com/posts/1")

    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "headers"))
    assert.isEq(resp.body, "{}")
    assert.isEq(resp.headers["Content-Type"], "application/json; charset=utf-8")
})

test.run("HTTP POST request", func(assert) {
    const data = json.encode({
        "title": 'foo',
        "body": 'bar',
        "userId": 1,
    })

    const resp = http.post("https://jsonplaceholder.typicode.com/posts", data, {
        "headers": {
            "Content-Type": "application/json",
        },
    })

    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "headers"))
    assert.isTrue(isString(resp.body))

    const respData = json.decode(resp.body)
    assert.isTrue(isMap(respData))
    assert.isTrue(contains(respData, "id"))
    assert.isTrue(contains(respData, "title"))
    assert.isTrue(contains(respData, "body"))
    assert.isTrue(contains(respData, "userId"))
})

test.run("HTTP POST request automatic encoding", func(assert) {
    const data = {
        "title": 'foo',
        "body": 'bar',
        "userId": 1,
    }

    const resp = http.post("https://jsonplaceholder.typicode.com/posts", data)

    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "headers"))
    assert.isTrue(isString(resp.body))

    const respData = json.decode(resp.body)
    assert.isTrue(isMap(respData))
    assert.isTrue(contains(respData, "id"))
    assert.isTrue(contains(respData, "title"))
    assert.isTrue(contains(respData, "body"))
    assert.isTrue(contains(respData, "userId"))
})

test.run("HTTP PUT request automatic encoding", func(assert) {
    const data = {
        "title": 'foo2',
        "body": 'bar2',
        "userId": 5,
    }

    const resp = http.put("https://jsonplaceholder.typicode.com/posts/1", data)

    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "headers"))
    assert.isTrue(isString(resp.body))

    const respData = json.decode(resp.body)
    assert.isTrue(isMap(respData))
    assert.isTrue(contains(respData, "id"))
    assert.isTrue(contains(respData, "title"))
    assert.isTrue(contains(respData, "body"))
    assert.isTrue(contains(respData, "userId"))
})

test.run("HTTP PATCH request automatic encoding", func(assert) {
    const data = {
        "title": 'foo3',
    }

    const resp = http.patch("https://jsonplaceholder.typicode.com/posts/1", data)

    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "headers"))
    assert.isTrue(isString(resp.body))

    const respData = json.decode(resp.body)
    assert.isTrue(isMap(respData))
    assert.isTrue(contains(respData, "id"))
    assert.isTrue(contains(respData, "title"))
    assert.isTrue(contains(respData, "body"))
    assert.isTrue(contains(respData, "userId"))
})
