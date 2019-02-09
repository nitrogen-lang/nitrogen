import "std/test"

test.run("Simple function call", func(assert) {
    const someFunc = func(arg1) {
        arg1
    }

    assert.isEq(someFunc('Hello'), 'Hello')
})

test.run("Simple function call no args", func(assert) {
    const someFunc = func() { 'called' }

    assert.isEq(someFunc(), 'called')
})

test.run("Function call optional args", func(assert) {
    const someFunc = func(arg1) {
        toString(arguments)
    }

    assert.isEq(someFunc('Hello', 'World'), '["World"]')
})

test.run("Function call no required args", func(assert) {
    const someFunc = func(arg1) {
        pass
    }

    assert.shouldThrow(func() {
        someFunc()
    })
})
