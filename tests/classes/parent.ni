const conf = 'This thing'
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

let myPrinter = make printer(1)

if !grandparentInitRan {
    println("Test Failed: Grandparent init() did not run")
    exit(1)
}

let expected = 'ID: 1 Msg: Hello'
let test = myPrinter.doStuff('Hello')
if test != expected {
    println("Test Failed: instance function failed. Expected ", expected, ", got ", test)
    exit(1)
}

expected = 'parent thing'
test = myPrinter.z
if test != expected {
    println("Test Failed: inherited property not right. Expected ", expected, ", got ", test)
    exit(1)
}

expected = "I'm the parent"
test = myPrinter.parentOnly()
if test != expected {
    println("Test Failed: parent only method failed. Expected ", expected, ", got ", test)
    exit(1)
}

expected = 'Parent: parent thing Msg: Hello'
test = myPrinter.doStuff2('Hello')
if test != expected {
    println("Test Failed: calling parent method. Expected ", expected, ", got ", test)
    exit(1)
}

expected = 'Grandparent: Hello'
test = myPrinter.doStuff3('Hello')
if test != expected {
    println("Test Failed: calling grand-parent method. Expected ", expected, ", got ", test)
    exit(1)
}
