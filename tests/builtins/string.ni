import 'string'
import "test"

func String(s) {
    return make string.String(s)
}

test.run("String module contains", func(assert) {
    const str = String("one,two,three")
    assert.isTrue(str.contains("two"))
    assert.isFalse(str.contains("Two"))
})

test.run("String module count", func(assert) {
    const str = String("one,two,three")
    assert.isEq(str.count(","), 2)
    assert.shouldThrow(func() {
        str.count("")
    })
})

test.run("String module dedup", func(assert) {
    const str = String("name:    John")
    const expected = "name: John"
    assert.isEq(str.dedup(" "), expected)
})

test.run("String module format", func(assert) {
    const str = make string.String("My name is {} and I'm {} years old")
    const expected = "My name is John and I'm 25 years old"
    assert.isEq(str.format("John", 25), expected)
})

test.run("String module prefix/suffix", func(assert) {
    const str = String("abcefguvwxyz")

    assert.isTrue(str.hasPrefix("abc"))
    assert.isTrue(str.hasSuffix("xyz"))
    assert.isFalse(str.hasPrefix("abcf"))
    assert.isFalse(str.hasSuffix("vxyz"))
})

test.run("String module replace", func(assert) {
    const str = String("oink oink oink")

    let expected = "oinky oinky oink"
    assert.isEq(str.replace("k", "ky", 2), expected)

    expected = "moo moo moo"
    assert.isEq(str.replace("oink", "moo", -1), expected)

    assert.shouldThrow(func() {
        str.replace("", "moo", -1)
    })
})

test.run("String module split", func(assert) {
    const str = String("one,two,three")

    let splits = str.split(",")
    assert.isEq(len(splits), 3)

    splits = str.splitN(",", 2)
    assert.isEq(len(splits), 2)
})

test.run("String module trim space", func(assert) {
    const str = String("       test        ")
    const expected = "test"
    assert.isEq(str.trimSpace(), expected)
})
