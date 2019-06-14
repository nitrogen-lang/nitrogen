import "std/test"

test.run("Simple function call", fn(assert) {
    const somefn = fn(arg1) {
        arg1
    }

    assert.isEq(somefn('Hello'), 'Hello')
})

test.run("Simple function call no args", fn(assert) {
    const somefn = fn() { 'called' }

    assert.isEq(somefn(), 'called')
})

test.run("function call optional args", fn(assert) {
    const somefn = fn(arg1) {
        toString(arguments)
    }

    assert.isEq(somefn('Hello', 'World'), '["World"]')
})

test.run("function call no required args", fn(assert) {
    const somefn = fn(arg1) {
        pass
    }

    assert.shouldThrow(fn() {
        somefn()
    })
})

test.run("function call with sugar syntax", fn(assert) {
    fn somefn(arg1) {
        arg1
    }

    assert.isEq(somefn('Hello'), 'Hello')
})
