import "stdlib/http"
import "stdlib/json"
import "stdlib/test"
import "stdlib/collections"

use collections.contains

test.run("HTTP GET request", func(assert) {
    const resp = http.get("https://jsonplaceholder.typicode.com/posts/1")

    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "headers"))
    assert.isTrue(isString(resp.body))
})

test.run("HTTP GET JSON request", func(assert) {
    const resp = http.getJSON("https://jsonplaceholder.typicode.com/posts/1")

    assert.isTrue(isMap(resp))
    assert.isTrue(contains(resp, "id"))
    assert.isTrue(contains(resp, "title"))
    assert.isTrue(contains(resp, "body"))
    assert.isTrue(contains(resp, "userId"))
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
