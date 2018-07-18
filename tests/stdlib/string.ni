if !modulesSupported() {
    return
}

import 'stdlib/string' as str
import "stdlib/test"

test.run("String module contains", func(assert) {
    const testStr = "one,two,three"
    assert.isTrue(str.contains(testStr, "two"))
    assert.isFalse(str.contains(testStr, "Two"))
})

test.run("String module count", func(assert) {
    const testStr = "one,two,three"
    assert.isEq(str.count(testStr, ","), 2)
    assert.shouldThrow(func() {
        str.count(testStr, "")
    })
})

test.run("String module dedup", func(assert) {
    const testStr = "name:    John"
    const expected = "name: John"
    assert.isEq(str.dedup(testStr, " "), expected)
})

test.run("String module format", func(assert) {
    const formatted = str.format("My name is {} and I'm {} years old", "John", 25)
    const expected = "My name is John and I'm 25 years old"

    assert.isEq(formatted, expected)
})

test.run("String module hasPrefix", func(assert) {
    const testStr = "abcefguvwxyz"

    assert.isTrue(str.hasPrefix(testStr, "abc"))
    assert.isFalse(str.hasPrefix(testStr, "abcf"))
})

test.run("String module hasSuffix", func(assert) {
    const testStr = "abcefguvwxyz"

    assert.isTrue(str.hasSuffix(testStr, "xyz"))
    assert.isFalse(str.hasSuffix(testStr, "vxyz"))
})

test.run("String module replace", func(assert) {
    const testStr = "oink oink oink"

    let expected = "oinky oinky oink"
    assert.isEq(str.replace(testStr, "k", "ky", 2), expected)

    expected = "moo moo moo"
    assert.isEq(str.replace(testStr, "oink", "moo", -1), expected)

    assert.shouldThrow(func() {
        str.replace(testStr, "", "moo", -1)
    })
})

test.run("String module split", func(assert) {
    const testStr = "one,two,three"

    const splits = str.split(testStr, ",")
    assert.isEq(len(splits), 3)
})

test.run("String module splitN", func(assert) {
    const testStr = "one,two,three"

    const splits = str.splitN(testStr, ",", 2)
    assert.isEq(len(splits), 2)
})

test.run("String module trim space", func(assert) {
    const testStr = "       test        "
    const expected = "test"
    assert.isEq(str.trimSpace(testStr), expected)
})
