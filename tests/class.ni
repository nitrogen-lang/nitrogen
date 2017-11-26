always conf = 'This thing'

let printer = class {
    let x
    always t = "Thing"

    func init(x) {
        this.x = x
    }

    func doStuff(msg) {
        return 'ID: ' + toString(x) + ' Msg: ' + msg
    }
}

let myPrinter = make printer(1);
let myPrinter2 = make printer(2);

if !is_a(myPrinter, printer) {
    println("myPrinter isn't a printer (class)")
}
if !is_a(myPrinter, 'printer') {
    println("myPrinter isn't a printer (string)")
}

if !is_a(myPrinter2, printer) {
    println("myPrinter2 isn't a printer (class)")
}
if !is_a(myPrinter2, 'printer') {
    println("myPrinter2 isn't a printer (string)")
}

if classOf(myPrinter) != 'printer' {
    println("myPrinter isn't classOf printer")
}
if classOf(myPrinter2) != 'printer' {
    println("myPrinter2 isn't classOf printer")
}

let mp1 = myPrinter.doStuff('Hello, world!')
if mp1 != 'ID: 1 Msg: Hello, world!' {
    println('myPrinter wrong message: ', mp1)
}

let mp2 = myPrinter2.doStuff('Hello, world!')
if mp2 != 'ID: 2 Msg: Hello, world!' {
    println('myPrinter2 wrong message: ', mp2)
}
