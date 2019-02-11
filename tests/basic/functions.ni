import "std/test"

test.run("Simple fntion call", fn(assert) {
    const somefn = fn(arg1) {
        arg1
    }

    assert.isEq(somefn('Hello'), 'Hello')
})

test.run("Simple fntion call no args", fn(assert) {
    const somefn = fn() { 'called' }

    assert.isEq(somefn(), 'called')
})

test.run("fntion call optional args", fn(assert) {
    const somefn = fn(arg1) {
        toString(arguments)
    }

    assert.isEq(somefn('Hello', 'World'), '["World"]')
})

test.run("fntion call no required args", fn(assert) {
    const somefn = fn(arg1) {
        pass
    }

    assert.shouldThrow(fn() {
        somefn()
    })
})
