import "std/test"

let grandparentInitRan = false

class grandParent {
    const init = fn() {
        grandparentInitRan = true
    }

    const doStuff = fn(msg) {
        return 'Grandparent: '+msg
    }
}

class parentPrinter ^ grandParent {
    let z

    const init = fn() {
        parent()
        this.z = "parent thing"
    }

    const doStuff = fn(msg) {
        return 'Parent: ' + this.z + ' Msg: ' + msg
    }

    const parentOnly = fn() {
        return "I'm the parent"
    }

    const doStuff2 = fn(msg) {
        return parent.doStuff(msg)
    }
}

class printer ^ parentPrinter {
    let x
    const t = "Thing"

    const init = fn(x) {
        parent()
        this.x = x
    }

    // Overloaded fntion
    const doStuff = fn(msg) {
        return 'ID: ' + toString(this.x) + ' Msg: ' + msg
    }

    const doStuff2 = fn(msg) {
        return parent.doStuff(msg)
    }

    const doStuff3 = fn(msg) {
        return parent.doStuff2(msg)
    }

    const setX = fn(x) {
        this.x = x
    }
}

test.run("Class inheritance", fn(assert, check) {
    let myPrinter = new printer(1)
    check(assert.isTrue(grandparentInitRan))

    let expected = 'ID: 1 Msg: Hello'
    check(assert.isEq(myPrinter.doStuff('Hello'), expected)) // Overloaded method

    expected = 'parent thing'
    check(assert.isEq(myPrinter.z, expected))

    expected = "I'm the parent"
    check(assert.isEq(myPrinter.parentOnly(), expected))

    expected = 'Parent: parent thing Msg: Hello'
    check(assert.isEq(myPrinter.doStuff2('Hello'), expected))

    expected = 'Grandparent: Hello'
    check(assert.isEq(myPrinter.doStuff3('Hello'), expected))
})
