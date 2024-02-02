import "std/test"

class printer {
    let x
    const t = "Thing"

    const init = fn(x) {
        this.x = x
    }

    const doStuff = fn(msg) {
        return 'ID: ' + toString(this.x) + ' Msg: ' + msg
    }
}

test.run("Basic classes", fn(assert, check) {
    const myPrinter = new printer(1)
    const myPrinter2 = new printer(2)

    check(assert.isTrue(instanceOf(myPrinter, printer)))
    check(assert.isTrue(instanceOf(myPrinter, '__main.printer')))

    check(assert.isTrue(instanceOf(myPrinter2, printer)))
    check(assert.isTrue(instanceOf(myPrinter2, '__main.printer')))

    check(assert.isEq('__main.printer', classOf(myPrinter)))
    check(assert.isEq('__main.printer', classOf(myPrinter2)))

    const mp1 = myPrinter.doStuff('Hello, world!')
    check(assert.isEq(mp1, 'ID: 1 Msg: Hello, world!'))

    const mp2 = myPrinter2.doStuff('Hello, world!')
    check(assert.isEq(mp2, 'ID: 2 Msg: Hello, world!'))
})
