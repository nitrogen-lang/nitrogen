import "std/test"

let grandparentInitRan = false

class grandParent {
    const init = func() {
        grandparentInitRan = true
    }

    const doStuff = func(msg) {
        return 'Grandparent: '+msg
    }
}

class parentPrinter ^ grandParent {
    let z

    const init = func() {
        parent()
        this.z = "parent thing"
    }

    const doStuff = func(msg) {
        return 'Parent: ' + this.z + ' Msg: ' + msg
    }

    const parentOnly = func() {
        return "I'm the parent"
    }

    const doStuff2 = func(msg) {
        return parent.doStuff(msg)
    }
}

class printer ^ parentPrinter {
    let x
    const t = "Thing"

    const init = func(x) {
        parent()
        this.x = x
    }

    // Overloaded function
    const doStuff = func(msg) {
        return 'ID: ' + toString(this.x) + ' Msg: ' + msg
    }

    const doStuff2 = func(msg) {
        return parent.doStuff(msg)
    }

    const doStuff3 = func(msg) {
        return parent.doStuff2(msg)
    }

    const setX = func(x) {
        this.x = x
    }
}

test.run("Class inheritance", func(assert) {
    let myPrinter = new printer(1)
    assert.isTrue(grandparentInitRan)

    let expected = 'ID: 1 Msg: Hello'
    assert.isEq(myPrinter.doStuff('Hello'), expected) // Overloaded method

    expected = 'parent thing'
    assert.isEq(myPrinter.z, expected)

    expected = "I'm the parent"
    assert.isEq(myPrinter.parentOnly(), expected)

    expected = 'Parent: parent thing Msg: Hello'
    assert.isEq(myPrinter.doStuff2('Hello'), expected)

    expected = 'Grandparent: Hello'
    assert.isEq(myPrinter.doStuff3('Hello'), expected)
})
