import "std/test"

test.run("Simple function call", fn(assert, check) {
    const somefn = fn(arg1) {
        arg1
    }

    check(assert.isEq(somefn('Hello'), 'Hello'))
})

test.run("Simple function call no args", fn(assert, check) {
    const somefn = fn() { 'called' }

    check(assert.isEq(somefn(), 'called'))
})

test.run("function call optional args", fn(assert, check) {
    const somefn = fn(arg1) {
        toString(arguments)
    }

    check(assert.isEq(somefn('Hello', 'World'), '["World"]'))
})

test.run("function call no required args", fn(assert, check) {
    const somefn = fn(arg1) {
        pass
    }

    check(assert.shouldRecover(fn() { somefn() }))
})

test.run("function call with sugar syntax", fn(assert, check) {
    fn somefn(arg1) {
        arg1
    }

    check(assert.isEq(somefn('Hello'), 'Hello'))
})
