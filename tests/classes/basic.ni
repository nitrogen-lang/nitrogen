import "std/test"

class printer {
    let x
    const t = "Thing"

    const init = func(x) {
        this.x = x
    }

    const doStuff = func(msg) {
        return 'ID: ' + toString(this.x) + ' Msg: ' + msg
    }
}

test.run("Basic classes", func(assert) {
    const myPrinter = new printer(1)
    const myPrinter2 = new printer(2)

    assert.isTrue(instanceOf(myPrinter, printer))
    assert.isTrue(instanceOf(myPrinter, 'printer'))

    assert.isTrue(instanceOf(myPrinter2, printer))
    assert.isTrue(instanceOf(myPrinter2, 'printer'))

    assert.isEq('printer', classOf(myPrinter))
    assert.isEq('printer', classOf(myPrinter2))

    const mp1 = myPrinter.doStuff('Hello, world!')
    assert.isEq(mp1, 'ID: 1 Msg: Hello, world!')

    const mp2 = myPrinter2.doStuff('Hello, world!')
    assert.isEq(mp2, 'ID: 2 Msg: Hello, world!')
})
