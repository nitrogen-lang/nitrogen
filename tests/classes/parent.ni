import "test"

let grandparentInitRan = false

class grandParent {
    func init() {
        grandparentInitRan = true
    }

    func doStuff(msg) {
        return 'Grandparent: '+msg
    }
}

class parentPrinter ^ grandParent {
    let z

    func init() {
        parent()
        this.z = "parent thing"
    }

    func doStuff(msg) {
        return 'Parent: ' + this.z + ' Msg: ' + msg
    }

    func parentOnly() {
        return "I'm the parent"
    }

    func doStuff2(msg) {
        return parent.doStuff(msg)
    }
}

class printer ^ parentPrinter {
    let x
    const t = "Thing"

    func init(x) {
        parent()
        this.x = x
    }

    // Overloaded function
    func doStuff(msg) {
        return 'ID: ' + toString(this.x) + ' Msg: ' + msg
    }

    func doStuff2(msg) {
        return parent.doStuff(msg)
    }

    func doStuff3(msg) {
        return parent.doStuff2(msg)
    }

    func setX(x) {
        this.x = x
    }
}

test.run("Class inheritance", func(assert) {
    let myPrinter = make printer(1)
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
