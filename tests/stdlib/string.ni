import "std/string"
import "std/test"

test.run("String module contains", fn(assert, check) {
    const testStr = "one,two,three"
    check(check(assert.isTrue(string.contains(testStr, "two"))))
    check(check(assert.isFalse(string.contains(testStr, "Two"))))
})

test.run("String module count", fn(assert, check) {
    const testStr = "one,two,three"
    check(assert.isEq(string.count(testStr, ","), 2))
    check(assert.shouldRecover(fn() {
        string.count(testStr, "")
    }))
})

test.run("String module dedup", fn(assert, check) {
    const testStr = "name:    John"
    const expected = "name: John"
    check(assert.isEq(string.dedup(testStr, " "), expected))
})

test.run("String module format", fn(assert, check) {
    const formatted = string.format("My name is {} and I'm {} years old", "John", 25)
    const expected = "My name is John and I'm 25 years old"

    check(assert.isEq(formatted, expected))
})

test.run("String module hasPrefix", fn(assert, check) {
    const testStr = "abcefguvwxyz"

    check(assert.isTrue(string.hasPrefix(testStr, "abc")))
    check(assert.isFalse(string.hasPrefix(testStr, "abcf")))
})

test.run("String module hasSuffix", fn(assert, check) {
    const testStr = "abcefguvwxyz"

    check(assert.isTrue(string.hasSuffix(testStr, "xyz")))
    check(assert.isFalse(string.hasSuffix(testStr, "vxyz")))
})

test.run("String module replace", fn(assert, check) {
    const testStr = "oink oink oink"

    let expected = "oinky oinky oink"
    check(assert.isEq(string.replace(testStr, "k", "ky", 2), expected))

    expected = "moo moo moo"
    check(assert.isEq(string.replace(testStr, "oink", "moo", -1), expected))

    check(assert.shouldRecover(fn() {
        string.replace(testStr, "", "moo", -1)
    }))
})

test.run("String module split", fn(assert, check) {
    const testStr = "one,two,three"

    const splits = string.split(testStr, ",")
    check(assert.isEq(len(splits), 3))
})

test.run("String module splitN", fn(assert, check) {
    const testStr = "one,two,three"

    const splits = string.splitN(testStr, ",", 2)
    check(assert.isEq(len(splits), 2))
})

test.run("String module trim space", fn(assert, check) {
    const testStr = "       test        "
    const expected = "test"
    check(assert.isEq(string.trimSpace(testStr), expected))
})

const String = fn(s) {
    return new string.String(s)
}

test.run("String module class contains", fn(assert, check) {
    const str = String("one,two,three")
    check(assert.isTrue(str.contains("two")))
    check(assert.isFalse(str.contains("Two")))
})

test.run("String module class count", fn(assert, check) {
    const str = String("one,two,three")
    check(assert.isEq(str.count(","), 2))
    check(assert.shouldRecover(fn() {
        str.count("")
    }))
})

test.run("String module class dedup", fn(assert, check) {
    const str = String("name:    John")
    const expected = "name: John"
    check(assert.isEq(str.dedup(" "), expected))
})

test.run("String module class format", fn(assert, check) {
    const str = new string.String("My name is {} and I'm {} years old")
    const expected = "My name is John and I'm 25 years old"
    check(assert.isEq(str.format("John", 25), expected))
})

test.run("String module class prefix/suffix", fn(assert, check) {
    const str = String("abcefguvwxyz")

    check(assert.isTrue(str.hasPrefix("abc")))
    check(assert.isTrue(str.hasSuffix("xyz")))
    check(assert.isFalse(str.hasPrefix("abcf")))
    check(assert.isFalse(str.hasSuffix("vxyz")))
})

test.run("String module class replace", fn(assert, check) {
    const str = String("oink oink oink")

    let expected = "oinky oinky oink"
    check(assert.isEq(str.replace("k", "ky", 2), expected))

    expected = "moo moo moo"
    check(assert.isEq(str.replace("oink", "moo", -1), expected))

    check(assert.shouldRecover(fn() {
        str.replace("", "moo", -1)
    }))
})

test.run("String module class split", fn(assert, check) {
    const str = String("one,two,three")

    let splits = str.split(",")
    check(assert.isEq(len(splits), 3))

    splits = str.splitN(",", 2)
    check(assert.isEq(len(splits), 2))
})

test.run("String module class trim space", fn(assert, check) {
    const str = String("       test        ")
    const expected = "test"
    check(assert.isEq(str.trimSpace(), expected))
})
