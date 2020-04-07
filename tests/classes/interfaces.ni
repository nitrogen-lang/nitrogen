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

test.run("Class implements interface", fn(assert) {
    assert.isTrue(StdOutPrinter implements Printer)
})

test.run("Instance implements interface", fn(assert) {
    const p = new StdOutPrinter()
    assert.isTrue(p implements Printer)
})

test.run("Another interface implements interface", fn(assert) {
    assert.isTrue(AdvancedPrinter implements Printer)
})
