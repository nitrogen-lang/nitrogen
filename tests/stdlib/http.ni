import "std/http"
import "std/encoding/json"
import "std/test"
import "std/collections"

use collections.contains

test.run("HTTP GET request", fn(assert, check) {
    const resp = http.get("https://jsonplaceholder.typicode.com/posts/1")

    check(assert.isTrue(contains(resp, "body")))
    check(assert.isTrue(contains(resp, "headers")))
    check(assert.isTrue(isString(resp.body)))
    check(assert.isNeq(resp.body, ""))
    check(assert.isEq(resp.headers["Content-Type"], "application/json; charset=utf-8"))
})

test.run("HTTP GET JSON request", fn(assert, check) {
    const resp = http.getJSON("https://jsonplaceholder.typicode.com/posts/1")

    check(assert.isTrue(isMap(resp)))
    check(assert.isTrue(contains(resp, "id")))
    check(assert.isTrue(contains(resp, "title")))
    check(assert.isTrue(contains(resp, "body")))
    check(assert.isTrue(contains(resp, "userId")))
})

test.run("HTTP HEAD request", fn(assert, check) {
    const resp = http.head("https://jsonplaceholder.typicode.com/posts/1")

    check(assert.isTrue(contains(resp, "body")))
    check(assert.isTrue(contains(resp, "headers")))
    check(assert.isEq(resp.body, ""))
    check(assert.isEq(resp.headers["Content-Type"], "application/json; charset=utf-8"))
})

test.run("HTTP DELETE request", fn(assert, check) {
    const resp = http.del("https://jsonplaceholder.typicode.com/posts/1")

    check(assert.isTrue(contains(resp, "body")))
    check(assert.isTrue(contains(resp, "headers")))
    check(assert.isEq(resp.body, "{}"))
    check(assert.isEq(resp.headers["Content-Type"], "application/json; charset=utf-8"))
})

test.run("HTTP POST request", fn(assert, check) {
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

    check(assert.isTrue(contains(resp, "body")))
    check(assert.isTrue(contains(resp, "headers")))
    check(assert.isTrue(isString(resp.body)))

    const respData = json.decode(resp.body)
    check(assert.isTrue(isMap(respData)))
    check(assert.isTrue(contains(respData, "id")))
    check(assert.isTrue(contains(respData, "title")))
    check(assert.isTrue(contains(respData, "body")))
    check(assert.isTrue(contains(respData, "userId")))
})

test.run("HTTP POST request automatic encoding", fn(assert, check) {
    const data = {
        "title": 'foo',
        "body": 'bar',
        "userId": 1,
    }

    const resp = http.post("https://jsonplaceholder.typicode.com/posts", data)

    check(assert.isTrue(contains(resp, "body")))
    check(assert.isTrue(contains(resp, "headers")))
    check(assert.isTrue(isString(resp.body)))

    const respData = json.decode(resp.body)
    check(assert.isTrue(isMap(respData)))
    check(assert.isTrue(contains(respData, "id")))
    check(assert.isTrue(contains(respData, "title")))
    check(assert.isTrue(contains(respData, "body")))
    check(assert.isTrue(contains(respData, "userId")))
})

test.run("HTTP PUT request automatic encoding", fn(assert, check) {
    const data = {
        "title": 'foo2',
        "body": 'bar2',
        "userId": 5,
    }

    const resp = http.put("https://jsonplaceholder.typicode.com/posts/1", data)

    check(assert.isTrue(contains(resp, "body")))
    check(assert.isTrue(contains(resp, "headers")))
    check(assert.isTrue(isString(resp.body)))

    const respData = json.decode(resp.body)
    check(assert.isTrue(isMap(respData)))
    check(assert.isTrue(contains(respData, "id")))
    check(assert.isTrue(contains(respData, "title")))
    check(assert.isTrue(contains(respData, "body")))
    check(assert.isTrue(contains(respData, "userId")))
})

test.run("HTTP PATCH request automatic encoding", fn(assert, check) {
    const data = {
        "title": 'foo3',
    }

    const resp = http.patch("https://jsonplaceholder.typicode.com/posts/1", data)

    check(assert.isTrue(contains(resp, "body")))
    check(assert.isTrue(contains(resp, "headers")))
    check(assert.isTrue(isString(resp.body)))

    const respData = json.decode(resp.body)
    check(assert.isTrue(isMap(respData)))
    check(assert.isTrue(contains(respData, "id")))
    check(assert.isTrue(contains(respData, "title")))
    check(assert.isTrue(contains(respData, "body")))
    check(assert.isTrue(contains(respData, "userId")))
})
