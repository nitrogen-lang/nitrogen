import "std/test"

interface Printer {
    print(data)
}

interface AdvancedPrinter {
    print(data)
    yell(data, other)
}

class StdOutPrinter {
    fn print(data) {
        println(data)
    }
}

test.run("Class implements interface", fn(assert, check) {
    check(assert.isTrue(StdOutPrinter implements Printer))
})

test.run("Instance implements interface", fn(assert, check) {
    const p = new StdOutPrinter()
    check(assert.isTrue(p implements Printer))
})

test.run("Another interface implements interface", fn(assert, check) {
    check(assert.isTrue(AdvancedPrinter implements Printer))
})
